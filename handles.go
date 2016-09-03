package main

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

func wechatHandler(ctx *fasthttp.RequestCtx) {
	str := ""
	if string(ctx.Method()) == "POST" {
		str = WeProcessMsg(string(ctx.Request.Body()))
	} else {
		str = WeCheckSign(string(ctx.FormValue("signature")),
			string(ctx.FormValue("timestamp")),
			string(ctx.FormValue("nonce")),
			string(ctx.FormValue("echostr")))
	}
	fmt.Fprintf(ctx, str)
}

func jokeHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, RandomJoke())
}

func adminHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "no admin")
}
