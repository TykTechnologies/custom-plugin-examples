package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	coprocess "github.com/TykTechnologies/tyk-protobuf/bindings/go"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

const (
	listenAddress = ":9111"
)

var (
	policiesToApply = []string{
		"admin",
	}
	mongoClient *mongo.Client
	collection  *mongo.Collection
)

type AccessToken struct {
	ID      string             `bson:"_id"`
	Ttl     int64              `bson:"ttl"`
	Created primitive.DateTime `bson:"created"`
}

func init() {
	// bootstrapping the user DB

	//Here you write the code for db connection
	// Set client options
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URL"))

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logrus.Fatal(err)
	}

	// Check the connection
	mongoClient = client
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		logrus.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	// Get a handle for your collection
	collection = client.Database(os.Getenv("DATABASE")).Collection(os.Getenv("COLLECTION"))
}

func main() {
	lis, err := net.Listen("tcp", listenAddress)
	fatalOnError(err, "failed to start tcp listener")

	logrus.Infof("starting grpc middleware on %s", listenAddress)
	s := grpc.NewServer()
	coprocess.RegisterDispatcherServer(s, &Dispatcher{})

	healthService := NewHealthChecker()
	grpc_health_v1.RegisterHealthServer(s, healthService)

	fatalOnError(s.Serve(lis), "unable to start grpc middleware")
}

type Dispatcher struct{}

func (d *Dispatcher) Dispatch(ctx context.Context, object *coprocess.Object) (*coprocess.Object, error) {
	switch object.HookName {
	case "MongoAuth":
		println("calling MongoAuth")
		return MongoAuth(object)
	}
	logrus.Warnf("unknown hook: %v", object.HookName)

	return object, nil
}

func (d *Dispatcher) DispatchEvent(ctx context.Context, event *coprocess.Event) (*coprocess.EventReply, error) {
	return &coprocess.EventReply{}, nil
}

func MongoAuth(object *coprocess.Object) (*coprocess.Object, error) {
	authKey := object.Request.Headers["Authorization"]
	// try to get session by API key
	// println(authKey)
	isValid, extractorDeadline := validateAccessToken(collection, authKey)
	// println(isValid)

	if !isValid {
		return failAuth(object, "Access forbidden")
	}

	// Set the ID extractor deadline, useful for caching valid keys:
	object.Session = &coprocess.SessionState{
		LastUpdated:         time.Now().String(),
		Rate:                50,
		Per:                 10,
		QuotaMax:            int64(0),
		QuotaRenews:         time.Now().Unix(),
		IdExtractorDeadline: extractorDeadline,
		Metadata: map[string]string{ // doesn't work if we don't set metadata
			"token": authKey,
		},
		ApplyPolicies: policiesToApply,
	}

	return object, nil
}

func fatalOnError(err error, msg string) {
	if err != nil {
		logrus.WithError(err).Fatal(msg)
	}
}

func failAuth(object *coprocess.Object, msg string) (*coprocess.Object, error) {
	object.Request.ReturnOverrides.ResponseCode = http.StatusForbidden
	object.Request.ReturnOverrides.ResponseError = msg
	return object, nil
}

// check if accesstoken is valid by comparing if ttl + created is greater than current time
func validateAccessToken(collection *mongo.Collection, token string) (bool, int64) {
	var accesstoken AccessToken
	err := collection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: token}}).Decode(&accesstoken)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the _id %s\n", token)
		return false, 0
	}
	if err != nil {
		logrus.Fatal(err)
	}
	expiry := accesstoken.Ttl + time.Unix(accesstoken.Created.Time().Unix(), 0).Unix() // ttl + created
	if expiry > time.Now().Unix() {
		return true, expiry
	}
	return false, 0 // already expired
}
