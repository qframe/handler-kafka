package main

import (
	"log"
	"os"
	"sync"
	"github.com/zpatrick/go-config"
	"github.com/codegangsta/cli"

	"github.com/qnib/qframe-types"
	"github.com/qframe/handler-kafka"
	"github.com/qframe/collector-docker-events"
	"github.com/qframe/types/qchannel"
)

func check_err(pname string, err error) {
	if err != nil {
		log.Printf("[EE] Failed to create %s plugin: %s", pname, err.Error())
		os.Exit(1)
	}
}

func NewConfig(kv map[string]string) *config.Config {
	return config.NewConfig([]config.Provider{config.NewStatic(kv)})
}

func Run(ctx *cli.Context) {
	// Create conf
	log.Printf("[II] Start Version: %s", ctx.App.Version)

	kv := map[string]string{
		"log.level": "debug",
		"log.only-plugins": "kafka",
		"handler.log.inputs": "docker-events",
		"handler.kafka.inputs": "docker-events",
	}

	cfg := config.NewConfig([]config.Provider{config.NewStatic(kv)})
	qChan := qtypes_qchannel.NewQChan()
	qChan.Broadcast()
	//////// Handlers
	phk, err := qhandler_kafka.New(qChan, cfg, "kafka")
	check_err(phk.Name, err)
	go phk.Run()
	//////// Collectors
	// Docker Logs collector
	pcde, err := qcollector_docker_events.New(qChan, cfg, "docker-events")
	check_err(pcde.Name, err)
	go pcde.Run()
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

func main() {
	app := cli.NewApp()
	app.Name = "kafka example agent"
	app.Usage = "kafka [options]"
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Value: "qframe.yml",
			Usage: "Config file, will overwrite flag default if present.",
		},
	}
	app.Action = Run
	app.Run(os.Args)
}
