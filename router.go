/**
 * REST API router
 * Rosbit Xu
 */
package main

import (
	"github.com/rosbit/http-helper"
	"cdc-gateway/common-endpoints"
	"cdc-gateway/cache"
	"cdc-gateway/conf"
	"net/http"
	"fmt"
)

func StartService() error {
	cache.StartCleaningThread()

	api := helper.NewHelper(helper.WithLogger("cloud-direct-connect"))

	serviceConf := gwconf.ServiceConf
	for i, _ := range serviceConf.Apps {
		service := &serviceConf.Apps[i]
		gwconf.AddCDCConf(service)
	}

	commonEndpoints := &serviceConf.CommonEndpoints
	api.Get(commonEndpoints.HealthCheck, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "OK\n")
	})
	api.POST(commonEndpoints.MakeRequest, ce.MakeRequest)
	api.POST(commonEndpoints.ParseResponse, ce.ParseResponse)
	api.POST(commonEndpoints.GetTransInfo, ce.GetTransInfo)

	listenParam := fmt.Sprintf("%s:%d", serviceConf.ListenHost, serviceConf.ListenPort)
	fmt.Printf("%v\n", http.ListenAndServe(listenParam, api))
	return nil
}

