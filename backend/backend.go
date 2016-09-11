package backend

import (
	"sync"

	loralog "github.com/panyingyun/wath/loralog"

	"github.com/eclipse/paho.mqtt.golang"
)

// Backend implements a MQTT pub-sub backend.
type Backend struct {
	conn   mqtt.Client
	rxChan chan models.TabRxData
	mutex  sync.RWMutex
	topic  string
}

// NewBackend creates a new Backend.
func NewBackend(server, username, password string) (*Backend, error) {
	b := Backend{
		rxChan: make(chan models.TabRxData),
		topic:  "",
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(server)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetOnConnectHandler(b.onConnected)
	opts.SetConnectionLostHandler(b.onConnectionLost)

	loralog.Infof("backend/mqttpubsub: connecting to mqtt broker %v", server)
	b.conn = mqtt.NewClient(opts)
	if token := b.conn.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return &b, nil
}

// Close closes the backend.
func (b *Backend) Close() {
	b.conn.Disconnect(250) // wait 250 milisec to complete pending actions
	loralog.Info("backend/mqttpubsub: Disconnect mqtt broker")
}

// RxDataChan returns the TabRxData channel.
func (b *Backend) RxDataChan() chan models.TabRxData {
	return b.rxChan
}

// Subscribe RxData
func (b *Backend) SubscribeTopic(topic string) error {
	defer b.mutex.Unlock()
	b.mutex.Lock()

	loralog.Infof("backend/mqttpubsub: subscribing to topic %v", topic)
	if token := b.conn.Subscribe(topic, 0, b.rxDataHandler); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	b.topic = topic
	return nil
}

// UnSubscribe RxData
func (b *Backend) UnSubscribeTopic(topic string) error {
	defer b.mutex.Unlock()
	b.mutex.Lock()

	loralog.Infof("backend/mqttpubsub: unsubscribing from topic %v", topic)
	if token := b.conn.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	b.topic = ""
	return nil
}

func (b *Backend) rxDataHandler(c mqtt.Client, msg mqtt.Message) {
	loralog.Infof("backend/mqttpubsub: packet received from topic %v", msg.Topic())
	var rxData models.TabRxData
	if err := rxData.Unmarshal(msg.Payload(), msg.Topic()); err != nil {
		loralog.Errorf("backend/mqttpubsub: decode rxData error: %s", err)
		return
	}
	b.rxChan <- rxData
}

func (b *Backend) onConnected(c mqtt.Client) {
	if b.topic == "" {
		return
	}
	defer b.mutex.RUnlock()
	b.mutex.RLock()

	loralog.Info("backend/mqttpubsub: onConnected to mqtt broker")
	if token := b.conn.Subscribe(b.topic, 0, b.rxDataHandler); token.Wait() && token.Error() != nil {
		return
	}
}

func (b *Backend) onConnectionLost(c mqtt.Client, reason error) {
	loralog.Errorf("backend/mqttpubsub: mqtt onConnectionLost error: %s", reason)
}
