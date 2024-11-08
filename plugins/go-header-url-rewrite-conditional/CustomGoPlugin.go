package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/TykTechnologies/tyk/apidef"
	"github.com/TykTechnologies/tyk/ctx"
	"github.com/TykTechnologies/tyk/log"
	"github.com/TykTechnologies/tyk/regexp"
)

var logger = log.Get()

func setUpstream(rw http.ResponseWriter, r *http.Request) {
	upstream := "http://httpbingo.org"
	apiDef := ctx.GetDefinition(r)
	apiDef.Proxy.EnableLoadBalancing = true
	apiDef.Proxy.StructuredTargetList = apidef.NewHostListFromList([]string{upstream})
}

func init() {
	logger.Info("--- Go custom plugin v4 init success! ---- ")
}
