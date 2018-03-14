package main

import (
	"encoding/json"

	"time"

	"github.com/qiangxue/fasthttp-routing"
)

func printObject(v interface{}) string {
	res2B, _ := json.Marshal(v)
	return string(res2B)
}

func setServerHeaders(ctx *routing.Context) {
	ctx.Response.Header.Set("Server", ServerName)
	ctx.Response.Header.Set("X-Powered-By", PoweredBy+build)

	//ctx.Response.Header.Set("Content-Type", "text/plain")
}

func sendMetrics(service string, statusCode int, startedAt time.Time, timeFinished time.Duration, host string, urlPath string, client_ip string, user_agent string) {
	if config.Clickhouse.IsEnabled {
		insertStatsToCH(service, statusCode, startedAt, timeFinished, host, urlPath, client_ip, user_agent)
	}
	Logger.Debugln(service, statusCode, startedAt.String(), timeFinished.Seconds(), urlPath, client_ip, user_agent)
}
