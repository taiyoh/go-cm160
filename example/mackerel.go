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

func NewMkrClient(apikey, hostid, name string) *MkrClient {
	client := mackerel.NewClient(apikey)
	loc, _ := time.LoadLocation("Asia/Tokyo")
	return &MkrClient{Client: client, HostId: hostid, Loc: loc, Name: name}
}

func (self *MkrClient) post(record *cm160.Record) {
	t := time.Date(record.Year, time.Month(record.Month), record.Day, record.Hour, record.Minute, time.Now().Second(), 0, self.Loc)
	err := self.Client.PostHostMetricValuesByHostID(self.HostId, []*mackerel.MetricValue{
		&mackerel.MetricValue{
			Name:  self.Name,
			Time:  t.Unix(),
			Value: record.Amps,
		},
	})
	if err != nil {
		log.Printf("mackerel error: %#v", err)
	}
}
