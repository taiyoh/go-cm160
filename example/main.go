package main

import (
	"github.com/taiyoh/go-cm160"
	"log"
	"time"
)

func main() {

	config := LoadConfig()

	sigCh := BuildSigWatcher()

	client := NewMkrClient(config.Mackerel, config.Name)

	device := cm160.Open()
	defer device.Close()

	log.Println("device initialized")

	go func() {
		for {
			select {
			case <-sigCh:
				log.Println("stop running")
				device.Stop()
			default:
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()

	for {
		if record := device.Read(); record != nil {
			if record.IsLive {
				log.Printf("live record amps: %#v\n", record.Amps)
			} else {
				log.Printf("not live at %d-%02d-%02d %02d:%02d amps: %#v\n", record.Year, record.Month, record.Day, record.Hour, record.Minute, record.Amps)
			}
			client.post(record)
		}
		if client.IsRunning() == false {
			break
		}
	}

	device.Run()

	log.Println("exit process")
}
