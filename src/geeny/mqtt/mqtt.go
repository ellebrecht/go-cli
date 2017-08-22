package data

import (
	"crypto/tls"
	"errors"
	"strings"
	"time"

	util "geeny/util"

	mqttLib "github.com/eclipse/paho.mqtt.golang"
)

// MQTT encapsulates mqtt communication. Currently wrapped around a 3rd party lib. @see https://eclipse.org/paho/clients/golang/, https://github.com/eclipse/paho.mqtt.golang
type MQTT struct {
	client  mqttLib.Client
	options *mqttLib.ClientOptions
}

// MQTTInterface defines MQTT functions
type MQTTInterface interface {
	ValidateTopic(topic string, contentTypeID string, thingID string) error
	Connect() error
	Subscribe(topic string) error
	SendData(data string, topic string, number int, interval time.Duration)
	SendDataForever(data string, topic string, interval time.Duration)
	Unsubscribe(topic string) error
	SetEndpoint(endpoint string)
	SetClientID(clientID string)
	SetMessageHandler(messageHandler MQTTMessageHandler)
	SetConnectHandler(connectHandler MQTTOnConnectHandler)
	SetConnectionLostHandler(connectHandler MQTTConnectionLostHandler)
	SetKeepAlive(keepAlive time.Duration)
	SetCertificate(cert tls.Certificate)
}

// MQTTMessageHandler fires when messages are received
type MQTTMessageHandler func(mqtt MQTT, topic string, payload string)

// MQTTOnConnectHandler fires on connect
type MQTTOnConnectHandler func(mqtt MQTT)

// MQTTConnectionLostHandler fires on connection lost
type MQTTConnectionLostHandler func(mqtt MQTT, err error)

// ValidateTopic returns an error if a topic is invalid
func (m *MQTT) ValidateTopic(topic string, contentTypeID string, thingID string) error {
	topicParts := strings.Split(topic, "/")
	if len(topicParts) != 3 {
		return errors.New("please make sure your topic is in this format: data/:contentTypeID/:thingID")
	}
	if strings.Compare(topicParts[0], "data") != 0 {
		return errors.New("bad topic. the first part of your topic should be data/")
	}
	if strings.Compare(topicParts[1], contentTypeID) != 0 {
		return errors.New("bad content type id in topic. " + topicParts[1] + " doesn't match " + contentTypeID)
	}
	if strings.Compare(topicParts[2], thingID) != 0 {
		return errors.New("bad thing id in topic. " + topicParts[2] + " doesn't match " + thingID)
	}
	if !util.IsValidUUID(topicParts[1]) {
		return errors.New(topicParts[1] + " is not a valid UUID")
	}
	if !util.IsValidUUID(topicParts[2]) {
		return errors.New(topicParts[2] + " is not a valid UUID")
	}
	return nil
}

// Connect to an mqtt endpoint
func (m *MQTT) Connect() error {
	m.client = mqttLib.NewClient(m.getOptions())
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// Subscribe to a topic
func (m *MQTT) Subscribe(topic string) error {
	// subscribe to the topic
	// at a maximum qos of zero, wait for the receipt to confirm the subscription
	if token := m.client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// Unsubscribe from a topic
func (m *MQTT) Unsubscribe(topic string) error {
	if token := m.client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	m.client.Disconnect(250)
	return nil
}

// SendData to a topic
func (m *MQTT) SendData(data string, topic string, number int, interval time.Duration) {
	for i := 0; i < int(number); i++ {
		m.sendData(data, topic)
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
}

// SendDataForever sends data to a topic forever
func (m *MQTT) SendDataForever(data string, topic string, interval time.Duration) {
	for true {
		m.sendData(data, topic)
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
}

// - setters

// SetEndpoint sets mqtt endpoint
func (m *MQTT) SetEndpoint(endpoint string) {
	m.getOptions().AddBroker("tls://" + endpoint + ":8883")
}

// SetClientID sets mqtt clientid
func (m *MQTT) SetClientID(clientID string) {
	m.getOptions().SetClientID("subscribe-" + clientID)
}

// SetMessageHandler sets mqtt message handler
func (m *MQTT) SetMessageHandler(messageHandler MQTTMessageHandler) {
	m.getOptions().SetDefaultPublishHandler(func(client mqttLib.Client, msg mqttLib.Message) {
		messageHandler(*m, string(msg.Topic()), string(msg.Payload()))
	})
}

// SetConnectHandler sets mqtt on connection handler
func (m *MQTT) SetConnectHandler(connectHandler MQTTOnConnectHandler) {
	m.getOptions().SetOnConnectHandler(func(client mqttLib.Client) {
		connectHandler(*m)
	})
}

// SetConnectionLostHandler sets mqtt connection lost handler
func (m *MQTT) SetConnectionLostHandler(connectionLostHandler MQTTConnectionLostHandler) {
	m.getOptions().SetConnectionLostHandler(func(client mqttLib.Client, err error) {
		connectionLostHandler(*m, err)
	})
}

// SetKeepAlive sets mqtt keep alive duration
func (m *MQTT) SetKeepAlive(keepAlive time.Duration) {
	m.getOptions().SetKeepAlive(keepAlive)
}

// SetCertificate sets mqtt certificate
func (m *MQTT) SetCertificate(cert tls.Certificate) {
	m.getOptions().SetTLSConfig(&tls.Config{
		Certificates: []tls.Certificate{cert},
	})
}

// - private

func (m *MQTT) getOptions() *mqttLib.ClientOptions {
	if m.options == nil {
		m.options = mqttLib.NewClientOptions()
	}
	return m.options
}

func (m *MQTT) sendData(data string, topic string) {
	token := m.client.Publish(topic, 0, false, data)
	token.Wait()
}
