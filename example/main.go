package main

import (
	"github.com/taiyoh/go-cm160"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const location string = "Asia/Tokyo"

func main() {

	config := LoadConfig("config.toml")

	device := cm160.Open()
	defer device.Close()

	log.Println("device initialized")

	// signal handling
	go func() {
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		caughtsig := <-ch
		log.Printf("signal caught: %#v exiting\n", caughtsig)
		device.Stop()
	}()

	// sender
	loc, _ := time.LoadLocation(location)
	client := NewMkrClient(config.Mackerel, config.Name)
	Send := func(record *cm160.Record) {
		now := time.Now()
		t := time.Date(record.Year, time.Month(record.Month), record.Day, record.Hour, record.Minute, now.Second(), 0, loc)
		if record.IsLive {
			log.Printf("live record amps: %#v\n", record.Amps)
		} else {
			log.Printf("not live at %s amps: %#v\n", t.Format("2006-01-02 15:04"), record.Amps)
		}
		// 10分以内のデータなら送信しよう
		if dur := now.Sub(t); dur.Minutes() < 10.0 {
			client.Post(record.Amps, t.Unix())
		}
	}

	// main loop
	for {
		if record := device.Read(); record != nil {
			Send(record)
		}
		if !device.IsRunning() {
			break
		}
	}

	log.Println("exit process")
}
