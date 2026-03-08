package delivery

import (
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/jbl1108/github.com/jbl1108/goSecret/usecases/datamodel"
	"github.com/jbl1108/github.com/jbl1108/goSecret/usecases/ports/input"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTClient struct {
	client       mqtt.Client
	topic        string
	inputUsecase input.KeyValInputPort
}

func NewMQTTClient(broker string, username string, password string, topic string, inputUsecase input.KeyValInputPort) *MQTTClient {
	opts := mqtt.NewClientOptions().AddBroker(broker)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetDefaultPublishHandler(defaultMessageHandler)
	opts.SetAutoReconnect(true)
	opts.SetConnectionLostHandler(mqtt.DefaultConnectionLostHandler)
	client := mqtt.NewClient(opts)
	return &MQTTClient{client: client, topic: topic, inputUsecase: inputUsecase}
}

func defaultMessageHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received message on topic: %s, payload: %s", msg.Topic(), string(msg.Payload()))
}

func (m *MQTTClient) Connect() {
	reader := m.client.OptionsReader()
	log.Printf("Connecting to MQTT broker at: %v", reader.Servers())
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	log.Printf("Connecting to topic: %v", m.topic)
	m.client.Subscribe(m.topic, 1, m.messageHandler)
}

func (m *MQTTClient) Disconnect() {
	m.client.Disconnect(250)
}

func (m *MQTTClient) messageHandler(client mqtt.Client, msg mqtt.Message) {
	var message datamodel.Message
	if err := json.Unmarshal(msg.Payload(), &message); err != nil {
		log.Printf("failed to parse message: %v", err)
		return
	}
	log.Printf("Received message: %+v", message)
	prefix, err := m.getTopic(msg.Topic())
	if err != nil {
		log.Printf("Error getting prefix: %v", err)
		return
	}
	message.Topic = prefix
	err = m.inputUsecase.SetKey(message)
	if err != nil {
		log.Printf("Error handling key value: %v", err)
	}

}
func (*MQTTClient) getTopic(topic string) (string, error) {
	parts := strings.Split(topic, "/")
	if len(parts) < 2 {
		return "", errors.New("Invalid topic format, expected 'keyvalue/{bucket}' got: " + topic)
	}
	if parts[0] != "keyvalue" {
		return "", errors.New("Invalid topic format, expected 'keyvalue/{bucket}' got: " + topic)
	}
	prefix := parts[1]
	return prefix, nil
}
