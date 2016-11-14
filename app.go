package main

import (
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/Financial-Times/coco-traffic-analyser/analyser"
	log "github.com/Sirupsen/logrus"
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

	app.Action = func() {
		a := analyser.New(*iface)
		a.Start()
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
