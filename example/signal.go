package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func BuildSigWatcher() chan bool {
	// signal trap
	sigchan := make(chan bool)
	go func() {
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		caughtsig := <-ch
		log.Printf("signal caught: %#v exiting\n", caughtsig)
		sigchan <- true
	}()
	log.Println("signal handling is ready")
	return sigchan
}
