package mqtt

import (
	"encoding/json"
	"errors"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	"time"
)

type MqttClient struct {
	client mqtt.Client
}

func NewMqttClient(serverUrl string) *MqttClient {
	opts := mqtt.NewClientOptions().AddBroker(serverUrl).SetClientID("raman")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetConnectTimeout(5*time.Second)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return &MqttClient{
		client: c,
	}
}

func (m *MqttClient) PublishSpectrumData(spectrum []float64) (err error) {
	if !m.client.IsConnected() {
		return errors.New("mqtt not connected")
	}

	bytes, err := json.Marshal(spectrum)
	if err != nil {
		return
	}

	logrus.Infof("About to send %d bytes of data to raman/data", len(bytes))
	publish := m.client.Publish("raman/data", 1, false, bytes)
	err = publish.Error()
	if err != nil {
		return
	}
	return 
}
