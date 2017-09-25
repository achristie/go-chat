package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/tcp"
)

const (
	configKey = "CHAT"
)

func init() {
	os.Setenv("CHAT_HOST", ":6000")
}

func main() {

	if err := cfg.Init(cfg.EnvProvider{Namespace: configKey}); err != nil {
		fmt.Println("error initializing configuration system", err)
		os.Exit(1)
	}

	log.Println("Configuration\n", cfg.Log())

	host := cfg.MustString("HOST")

	cfg := tcp.Config{
		NetType:     "tcp4",
		Addr:        host,
		ConnHandler: connHandler{},
		ReqHandler:  reqHandler{},
		RespHandler: respHandler{},

		OptEvent: tcp.OptEvent{
			Event: Event,
		},
	}

	t, err := tcp.New("Sample", cfg)
	if err != nil {
		log.Printf("main : %s", err)
		return
	}

	if err := t.Start(); err != nil {
		log.Printf("main : %s", err)
		return
	}
	defer t.Stop()

	log.Printf("main : Waiting for data on: %s", t.Addr())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
}

var evtTypes = []string{
	"unknown",
	"Accept",
	"Join",
	"Read",
	"Remove",
	"Drop",
	"Groom",
}

var typTypes = []string{
	"unknown",
	"Error",
	"Info",
	"Trigger",
}

func Event(evt, typ int, ipAddress string, format string, a ...interface{}) {
	log.Printf("****> EVENT : IP [ %s ] : EVT[%s] TYP[%s] : %s",
		ipAddress, evtTypes[evt], typTypes[typ], fmt.Sprintf(format, a...))
}
