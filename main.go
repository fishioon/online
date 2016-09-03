package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/valyala/fasthttp"
)

var (
	addr = flag.String("addr", ":8080", "TCP address to listen to")
)

func handle(ctx *fasthttp.RequestCtx) {
	log.Printf("%s %s\n", ctx.Method(), ctx.Path())
	switch string(ctx.Path()) {
	case "/wx/":
		wechatHandler(ctx)
	case "/joke/":
		jokeHandler(ctx)
	case "/admin/":
		adminHandler(ctx)
	default:
		ctx.Error("not found", fasthttp.StatusNotFound)
	}
}

func main() {
	rand.Seed(time.Now().Unix())
	port := flag.String("port", "8080", "port number, default: 8080")
	logpath := flag.String("log", "", "")

	flag.Parse()
	if *logpath != "" {
		f, err := os.OpenFile(*logpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()
		log.SetOutput(f)
	}

	jokes_num := ReloadJokes()
	log.Printf("load jokes num:%d\n", jokes_num)

	menu_num := ReloadMenu()
	log.Printf("load menu num:%d\n", menu_num)

	addr := fmt.Sprintf(":%s", *port)
	if err := fasthttp.ListenAndServe(addr, handle); err != nil {
		log.Fatal("Error in ListenAndServe: %s", err)
	}
}
