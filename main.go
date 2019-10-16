package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"sender/app"
)

func main() {
	go func() {
		if err := http.ListenAndServe("0.0.0.0:6060", nil);err != nil{
			log.Println(err)
		}
	}()
	app.Run()
}
