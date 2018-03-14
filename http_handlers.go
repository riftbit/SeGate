package main

import (
	"time"

	"strings"

	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

func OptionsHandler(ctx *routing.Context) error {
	ctx.Response.SetBody([]byte(`{"error": "this method not allowed"}`))
	ctx.Response.SetStatusCode(405)
	setServerHeaders(ctx)
	ctx.Abort()
	return nil
}

func StartHandler(ctx *routing.Context) error {

	nodeData, ok := NodesList[ctx.Param("node")]
	if ok != true {
		ctx.Response.SetBody([]byte(`{"error": "node not found"}`))
		ctx.Response.SetStatusCode(404)
		setServerHeaders(ctx)
		ctx.Abort()
		return nil
	}
	ctx.Set("TimeStarted", time.Now())
	ctx.Set("NodeData", nodeData)
	ctx.Set("UrlParams", strings.TrimPrefix(string(ctx.Path()), "/"+string(ctx.Param("node"))))
	return nil
}

func FinishHandler(ctx *routing.Context) error {
	startedAt := ctx.Get("TimeStarted").(time.Time)
	NodeData := ctx.Get("NodeData").(nodeElement)
	timeFinished := time.Since(startedAt)
	urlPath := ctx.Get("UrlParams").(string)

	ipAddress := ctx.RemoteIP().String()

	ctx.Request.Header.VisitAll(func(key, value []byte) {
		if string(key) == "X-Forwarded-For" {
			ipAddress = string(value)
		}
	})

	go sendMetrics(NodeData.NodeName, ctx.Response.StatusCode(), startedAt, timeFinished, NodeData.NodeURL, urlPath, ipAddress, string(ctx.Request.Header.UserAgent()))
	return nil
}

func ProxyHandler(ctx *routing.Context) error {

	nodeData := ctx.Get("NodeData").(nodeElement)

	req := fasthttp.AcquireRequest()

	ctx.Request.Header.VisitAll(func(key, value []byte) {
		req.Header.Set(string(key), string(value))
	})

	urlpart := ctx.Get("UrlParams").(string)
	req.SetRequestURI(nodeData.NodeURL + urlpart)

	req.Header.SetMethod("POST")

	crypto, _ := NewCryptoService(nodeData.AESKey)

	if Logger.Level.String() == "debug" {
		messdecv, err := crypto.EncryptBase64(ctx.Request.Body())
		if err != nil {
			Logger.Errorln(err)
		}
		Logger.Debugln("Encrypted request:", string(messdecv))
		ctx.Request.SetBody(messdecv)
	}

	messdec, err := crypto.DecryptBase64(ctx.Request.Body())
	if err != nil {
		Logger.Errorln(err)
	}
	Logger.Debugln("Decrypted request:", string(messdec))

	req.SetBody(messdec)

	resp := fasthttp.AcquireResponse()
	client := &fasthttp.Client{}
	client.Do(req, resp)

	resp.Header.VisitAll(func(key, value []byte) {
		ctx.Response.Header.Set(string(key), string(value))
	})

	mess, err := crypto.EncryptBase64(resp.Body())
	if err != nil {
		Logger.Errorln(err)
	}

	ctx.Response.SetBody(mess)
	ctx.Response.SetStatusCode(resp.StatusCode())
	ctx.Response.Header.SetContentType("text/plain")
	setServerHeaders(ctx)

	return nil
}
