package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// M2xClient : client for this app with M2X
type M2xClient struct {
	APIKey   string
	DeviceID string
}

type m2xParentValue struct {
	Values bmp180values `json:"values"`
}

type bmp180values struct {
	Ampare []m2xvalue `json:"ampare"`
}

type m2xvalue struct {
	Timestamp string `json:"timestamp"`
	Value     string `json:"value"`
}

// NewM2xClient : returns new MkrClient
func NewM2xClient(config M2xConfig) *M2xClient {
	return &M2xClient{DeviceID: config.Deviceid, APIKey: config.Apikey}
}

func (m *M2xClient) doRequest(values m2xParentValue) (*http.Response, error) {
	data, err := json.Marshal(values)
	if err != nil {
		log.Printf("json.Marshal error: %#v. value: %#v", err, values)
		return nil, err
	}
	reqbody := bytes.NewBuffer(data)
	//log.Printf("reqbody: %s", data)
	url := fmt.Sprintf("http://api-m2x.att.com/v2/devices/%s/updates", m.DeviceID)
	request, err := http.NewRequest("POST", url, reqbody)
	request.Header.Set("X-M2X-KEY", m.APIKey)
	request.Header.Set("Content-Type", "application/json")
	request.ContentLength = int64(reqbody.Len())
	client := &http.Client{}
	return client.Do(request)
}

// Post : posts temperature and pressure to m2x
func (m *M2xClient) Post(ampare float32, t time.Time) {
	ts := t.UTC().Format("2006-01-02T15:04:05Z")
	res, err := m.doRequest(m2xParentValue{
		bmp180values{
			[]m2xvalue{m2xvalue{ts, fmt.Sprintf("%f", ampare)}},
		},
	})
	if err != nil {
		log.Printf("m2x Post error: %s", err.Error())
		return
	}
	defer res.Body.Close()
	byteArray, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode != 202 {
		log.Printf("m2x Post error: %s", string(byteArray))
		return
	}
	//log.Printf("success: %s", string(byteArray))
}
