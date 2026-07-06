package mqtt

import (
	"context"
	"fmt"
	"log/slog"

	paho "github.com/eclipse/paho.mqtt.golang"

	"github.com/kbball/same-message-to-mesh/backend/internal/domain/entity"
	portsvc "github.com/kbball/same-message-to-mesh/backend/internal/domain/port/service"
)

// Publisher sends plain-text alert messages to an MQTT broker.
type Publisher struct {
	client paho.Client
	topic  string
}

var _ portsvc.AlertPublisher = (*Publisher)(nil)

// New connects to the MQTT broker described by cfg and returns a Publisher.
func New(cfg entity.MQTTConfig) (*Publisher, error) {
	opts := paho.NewClientOptions().
		AddBroker(fmt.Sprintf("tcp://%s:%d", cfg.Host, cfg.Port)).
		SetClientID("same-mesh-publisher").
		SetCleanSession(true).
		SetAutoReconnect(true).
		SetOnConnectHandler(func(_ paho.Client) {
			slog.Info("MQTT connected", "host", cfg.Host, "port", cfg.Port, "topic", cfg.PublishTopic)
		}).
		SetConnectionLostHandler(func(_ paho.Client, err error) {
			slog.Warn("MQTT connection lost", "error", err)
		})

	client := paho.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("connecting to MQTT broker: %w", token.Error())
	}

	return &Publisher{client: client, topic: cfg.PublishTopic}, nil
}

// Publish sends the formatted alert message to the configured MQTT topic.
func (p *Publisher) Publish(_ context.Context, _ entity.SAMEAlert, message string) error {
	token := p.client.Publish(p.topic, 1, false, message)
	token.Wait()
	if err := token.Error(); err != nil {
		return fmt.Errorf("publishing to MQTT: %w", err)
	}
	slog.Info("alert published to MQTT", "topic", p.topic, "message", message)
	return nil
}

// Close disconnects from the MQTT broker.
func (p *Publisher) Close() {
	p.client.Disconnect(250)
}
