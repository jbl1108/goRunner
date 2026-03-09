package delivery

import (
	"encoding/json"
	"log"

	"github.com/jbl1108/goRunner/usecases/datamodel"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTClient struct {
	client mqtt.Client
	topic  string
}

func NewMQTTClient(broker string, username string, password string, topic string) *MQTTClient {
	opts := mqtt.NewClientOptions().AddBroker(broker)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetDefaultPublishHandler(defaultMessageHandler)
	opts.SetAutoReconnect(true)
	opts.SetConnectionLostHandler(mqtt.DefaultConnectionLostHandler)
	client := mqtt.NewClient(opts)
	return &MQTTClient{client: client, topic: topic}
}

func (m *MQTTClient) TrainingUpdated(training datamodel.Training) error {
	payload, err := json.Marshal(training)
	if err != nil {
		return err
	}
	return m.postMessage("updated", payload)
}

func (m *MQTTClient) TrainingAdded(training datamodel.Training) error {
	payload, err := json.Marshal(training)
	if err != nil {
		return err
	}
	return m.postMessage("added", payload)
}

func (m *MQTTClient) TrainingDeleted(uid string) error {
	payload, err := json.Marshal(uid)
	if err != nil {
		return err
	}
	return m.postMessage("deleted", payload)
}

func (m *MQTTClient) postMessage(topicPostfix string, payload []byte) error {
	token := m.client.Publish(m.topic+"/"+topicPostfix, 0, false, payload)
	token.Wait()
	return token.Error()
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
}

func (m *MQTTClient) Disconnect() {
	m.client.Disconnect(250)
}
