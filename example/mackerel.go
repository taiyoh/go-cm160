package main

import (
	"github.com/mackerelio/mackerel-client-go"
	"log"
)

type MkrClient struct {
	Client *mackerel.Client
	Name   string
	HostId string
}

func NewMkrClient(config MackerelConfig, name string) *MkrClient {
	client := mackerel.NewClient(config.Apikey)
	return &MkrClient{Client: client, HostId: config.Hostid, Name: name}
}

func (m *MkrClient) post(amps float32, t int64) {
	err := m.Client.PostHostMetricValuesByHostID(m.HostId, []*mackerel.MetricValue{
		&mackerel.MetricValue{
			Name:  m.Name,
			Time:  t,
			Value: amps,
		},
	})
	if err != nil {
		log.Printf("mackerel post error: %#v", err)
	}
}
