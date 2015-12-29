package main

import (
	"github.com/mackerelio/mackerel-client-go"
	"log"
)

// MkrClient client for this app with mackerel
type MkrClient struct {
	Client *mackerel.Client
	Name   string
	HostID string
}

// NewMkrClient returns new MkrClient
func NewMkrClient(config MackerelConfig, name string) *MkrClient {
	client := mackerel.NewClient(config.Apikey)
	return &MkrClient{Client: client, HostID: config.Hostid, Name: name}
}

// Post post ampare to mackerel
func (m *MkrClient) Post(amps float32, t int64) {
	err := m.Client.PostHostMetricValuesByHostID(m.HostID, []*mackerel.MetricValue{
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
