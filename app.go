package main

import (
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	"time"

	"github.com/Financial-Times/coco-traffic-analyser/analyser"
	"github.com/Financial-Times/coco-traffic-analyser/resources"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/jawher/mow.cli"
)

func init() {
	f := &log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339Nano,
	}

	log.SetFormatter(f)
}

func main() {
	app := cli.App("coco-traffic-analyser", "Analayse communication traffic between CoCo microservices")

	iface := app.String(cli.StringOpt{
		Name:   "interface",
		Value:  "eth0",
		Desc:   "The network interface where to liste the traffic",
		EnvVar: "INTERFACE",
	})

	port := app.Int(cli.IntOpt{
		Name:   "port",
		Value:  8080,
		Desc:   "application port",
		EnvVar: "PORT",
	})

	app.Action = func() {
		a := analyser.New(*iface)

		go serve(":"+strconv.Itoa(*port), a)

		a.Start()

	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}

func serve(host string, a *analyser.StandardAnalyser) {
	h := resources.NewAnalyserHandler(a)
	r := mux.NewRouter()
	r.HandleFunc("/analyser/traffic-graph", h.ServeTrafficGraph).Methods("GET")
	err := http.ListenAndServe(host, r)
	if err != nil {
		log.Fatalf("Unable to start server: %v", err)
	}
}
