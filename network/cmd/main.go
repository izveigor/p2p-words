package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/izveigor/p2p-words/network/pkg/config"
	"github.com/izveigor/p2p-words/network/pkg/p2p"
)

func main() {
	p2p.ReadFiles()
	go p2p.InitHTTPServiceServer(config.Config.P2PSvcUrl)
	time.Sleep(time.Millisecond * 2000)
	go p2p.StartServer()
	time.Sleep(time.Millisecond * 2000)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		p2p.StopServer()
		os.Exit(1)
	}()

	for {
	}
}
