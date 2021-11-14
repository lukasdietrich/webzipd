package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/lukasdietrich/webzipd/internal/namespace"
	"github.com/lukasdietrich/webzipd/internal/render"
	"github.com/lukasdietrich/webzipd/internal/routing"
)

func main() {
	var (
		foldername string
		routemode  string
		address    string
	)

	flag.StringVar(
		&foldername,
		"folder",
		".",
		"Path to the folder containing content zips.")

	flag.StringVar(
		&routemode,
		"mode",
		"hostname",
		"Mode to use for routing. One of 'hostname' or 'path'.")

	flag.StringVar(
		&address,
		"address",
		":8080",
		"Address to listen on.")

	flag.Parse()

	index, err := namespace.OpenIndex(foldername)
	if err != nil {
		log.Fatalf("could not create index: %v", err)
	}

	renderer := render.NewRenderer(index)
	mux, err := routing.NewMux(renderer, routemode)
	if err != nil {
		log.Fatalf("could not create mux: %v", err)
	}

	if err := http.ListenAndServe(address, mux); err != nil {
		log.Fatalf("webzipd stopped with an error: %v", err)
	}
}
