module tyk-grpc-plugin

go 1.16

replace github.com/jensneuse/graphql-go-tools => github.com/TykTechnologies/graphql-go-tools v1.6.2-0.20221026084245-1fc4f5ca74bb

require (
	github.com/TykTechnologies/tyk v1.9.2-0.20221129200023-767c8b336fa5
	github.com/sirupsen/logrus v1.9.0
	go.mongodb.org/mongo-driver v1.11.1
	google.golang.org/grpc v1.51.0
)
