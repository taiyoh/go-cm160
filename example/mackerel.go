package main

import (
	"github.com/mackerelio/mackerel-client-go"
	"github.com/taiyoh/go-cm160"
	"log"
	"time"
)

type MkrClient struct {
	Client *mackerel.Client
	Loc    *time.Location
	Name   string
	HostId string
}

func NewMkrClient(config MackerelConfig, name string) *MkrClient {
	client := mackerel.NewClient(config.Apikey)
	loc, _ := time.LoadLocation("Asia/Tokyo")
	return &MkrClient{Client: client, HostId: config.Hostid, Loc: loc, Name: name}
}

func (self *MkrClient) post(record *cm160.Record) {
	now := time.Now()
	t := time.Date(record.Year, time.Month(record.Month), record.Day, record.Hour, record.Minute, now.Second(), 0, self.Loc)
	if dur := now.Sub(t); dur.Hours() < 24.0 {
		err := self.Client.PostHostMetricValuesByHostID(self.HostId, []*mackerel.MetricValue{
			&mackerel.MetricValue{
				Name:  self.Name,
				Time:  t.Unix(),
				Value: record.Amps,
			},
		})
		if err != nil {
			log.Printf("mackerel post error: %#v", err)
		}
	}
}
