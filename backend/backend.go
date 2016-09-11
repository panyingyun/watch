package backend

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/eclipse/paho.mqtt.golang"
)

//Lora Gateway and Application Message Receive
type LoraMsg struct {
	Msg   string
	Topic string
}

// Backend implements a MQTT pub-sub backend.
type Backend struct {
	conn   mqtt.Client
	rxChan chan LoraMsg
	topics map[string]struct{}
	mutex  sync.RWMutex
}

// NewBackend creates a new Backend.
func NewBackend(server, username, password string) (*Backend, error) {
	b := Backend{
		rxChan: make(chan LoraMsg),
		topics: make(map[string]struct{}),
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(server)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetOnConnectHandler(b.onConnected)
	opts.SetConnectionLostHandler(b.onConnectionLost)

	log.Infof(" [MQTT] connecting to mqtt server %v", server)
	b.conn = mqtt.NewClient(opts)
	if token := b.conn.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return &b, nil
}

// Close closes the backend.
func (b *Backend) Close() {
	b.conn.Disconnect(250) // wait 250 milisec to complete pending actions
	log.Info(" [MQTT] Disconnect mqtt server")
}

// RxDataChan returns the TabRxData channel.
func (b *Backend) RxDataChan() chan LoraMsg {
	return b.rxChan
}

// Subscribe RxData
func (b *Backend) SubscribeTopic(topic string) error {
	defer b.mutex.Unlock()
	b.mutex.Lock()

	log.Infof(" [MQTT] subscribing to topic %v", topic)
	if token := b.conn.Subscribe(topic, 0, b.rxDataHandler); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	b.topics[topic] = struct{}{}
	return nil
}

// UnSubscribe RxData
func (b *Backend) UnSubscribeTopic(topic string) error {
	defer b.mutex.Unlock()
	b.mutex.Lock()

	log.Infof(" [MQTT] unsubscribing from topic is %v", topic)
	if token := b.conn.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	delete(b.topics, topic)
	return nil
}

func (b *Backend) rxDataHandler(c mqtt.Client, msg mqtt.Message) {
	var rxData LoraMsg
	rxData.Msg = string(msg.Payload())
	rxData.Topic = msg.Topic()
	b.rxChan <- rxData
}

func (b *Backend) onConnected(c mqtt.Client) {
	log.Info(" [MQTT] ", "onConnected to mqtt broker")
	if len(b.topics) == 0 {
		log.Info(" [MQTT] ", "there is no topic here!")
		return
	}
	defer b.mutex.RUnlock()
	b.mutex.RLock()

	for topic := range b.topics {
		b.conn.Subscribe(topic, 2, b.rxDataHandler)
	}
}

func (b *Backend) onConnectionLost(c mqtt.Client, reason error) {
	log.Errorf(" [MQTT] onConnectionLost error %v", reason)
}
