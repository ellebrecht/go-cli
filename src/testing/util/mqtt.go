package util

import (
	"crypto/tls"
	"time"

	mqtt "geeny/mqtt"
)

type MockMQTT struct {
	Error error
}

func (m MockMQTT) ValidateTopic(topic string, contentTypeID string, thingID string) error {
	return m.Error
}
func (m MockMQTT) Connect() error {
	return m.Error
}
func (m MockMQTT) Subscribe(topic string) error {
	return m.Error
}
func (m MockMQTT) SendData(data string, topic string, number int, interval time.Duration) {
}
func (m MockMQTT) SendDataForever(data string, topic string, interval time.Duration) {
}
func (m MockMQTT) Unsubscribe(topic string) error {
	return m.Error
}
func (m MockMQTT) SetEndpoint(endpoint string) {
}
func (m MockMQTT) SetClientID(clientID string) {
}
func (m MockMQTT) SetMessageHandler(messageHandler mqtt.MQTTMessageHandler) {
}
func (m MockMQTT) SetConnectHandler(connectHandler mqtt.MQTTOnConnectHandler) {
}
func (m MockMQTT) SetConnectionLostHandler(connectHandler mqtt.MQTTConnectionLostHandler) {
}
func (m MockMQTT) SetKeepAlive(keepAlive time.Duration) {
}
func (m MockMQTT) SetCertificate(cert tls.Certificate) {
}
