package main

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	"github.com/open-policy-agent/opa/storage/inmem"

	"github.com/TykTechnologies/tyk/ctx"
	"github.com/TykTechnologies/tyk/log"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/storage"
	"github.com/spf13/viper"
)

var (
	opaLogger = log.Get()
)

// Compile Plugin: go build -buildmode=plugin -o ./tyk_go_plugins/authorize_opa/authorize_opa.so ./tyk_go_plugins/authorize_opa/authorize_opa.go
// Compile Gateway: go install -tags 'goplugin'

// PostAuthOpa doesn't do much at the moment, but it's supposed to query OPA to determine if the application
// is allowed to access an endpoint on behalf of a user
func PostAuthOpa(w http.ResponseWriter, r *http.Request) {
	opaLogger.Infof("oauth2.scope: %#v", strings.Split(r.Header.Get("X-Tyk-Plugin-OAuth2Introspect-Scope"), " "))
	opaLogger.Info("requestMethod: ", r.Method)
	opaLogger.Info("requestPath: ", r.URL.String())

	conf, err := getOPAFromConfigData(ctx.GetDefinition(r).ConfigData)
	if err != nil {
		opaLogger.Errorf("problem loading config data: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}

	regoInput := map[string]interface{}{
		"scopes": strings.Split(r.Header.Get("X-Tyk-Plugin-OAuth2Introspect-Scope"), " "),
		"method": r.Method,
		"path":   r.URL.String(),
	}

	regoInstance := rego.New(
		rego.Query("data.todos.allow"),
		rego.Compiler(conf.Policy),
		rego.Store(conf.DB),
		rego.Input(regoInput))

	resultSet, _ := regoInstance.Eval(context.TODO())

	valid := resultSet[0].Expressions[0].Value.(bool)
	if !valid {
		// dont allow the request through
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
		return
	}
}

type opaConf struct {
	DB     storage.Store
	Policy *ast.Compiler
}

func getOPAFromConfigData(configData map[string]interface{}) (*opaConf, error) {
	conf := &opaConf{}

	opaViper := viper.New()
	if err := opaViper.MergeConfigMap(configData); err != nil {
		return nil, err
	}

	policyB64 := opaViper.GetString("opa.policy")
	opaLogger.Debugf("opa.policy %s", policyB64)
	if policyB64 == "" {
		return conf, errors.New("opa.policy not found")
	}

	db := opaViper.GetStringMap("opa.db")
	opaLogger.Debugf("opa.db %#v", db)
	if db == nil {
		return conf, errors.New("opa.db not found")
	}

	policyBytes, err := base64.StdEncoding.DecodeString(policyB64)
	if err != nil {
		return conf, errors.New("policy not in b64 format")
	}

	conf.Policy, err = ast.CompileModules(map[string]string{
		opaViper.GetString("opa.module_name") + ".rego": string(policyBytes),
	})
	if err != nil {
		return conf, err
	}
	conf.DB = inmem.NewFromObject(db)

	return conf, nil
}

func main() {}
