package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"example.com/gem"
)

func main() {
	flag.Usage = func() {
		fmt.Println("run the demo application from the gem-spaas team")
		os.Exit(2)
	}
	listen := flag.String("a", ":8888", "http address")
	flag.Parse()

	g := gem.New()

	http.Handle("/productionplan", gem.GetPlan(g))
	http.Handle("/ws", gem.GetWebsocket(g))
	if err := http.ListenAndServe(*listen, nil); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
