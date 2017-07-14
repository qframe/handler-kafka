package qhandler_kafka

import (
	"fmt"
	"sync"
	"github.com/zpatrick/go-config"
	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/qnib/qframe-types"
	"github.com/qframe/types/inventory"
)

const (
	version = "0.0.0"
	pluginTyp = qtypes.HANDLER
	pluginPkg = "kafka"
)

type Plugin struct {
    qtypes.Plugin
	producer *kafka.Producer
	mutex sync.Mutex
	deliveryChan chan kafka.Event
}

func New(qChan qtypes.QChan, cfg *config.Config, name string) (Plugin, error) {
	var err error
	p := Plugin{
		Plugin: qtypes.NewNamedPlugin(qChan, cfg, pluginTyp, pluginPkg, name, version),
		deliveryChan: make(chan kafka.Event),
	}
	return p, err
}

// Connect creates a connection to InfluxDB
func (p *Plugin) Connect() {
	var err error
	brokers := p.CfgStringOr("broker", "localhost:9092")
	p.producer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": brokers})
	if err != nil {
		msg := fmt.Sprintf("Failed to create producer: %s\n", err)
		p.Log("error", msg)
	} else {
		msg := fmt.Sprintf("Created Producer ID:%d Name:%s\n", p.MyID, p.Name)
		p.Log("info", msg)
	}
}


// Run fetches everything from the Data channel and flushes it to stdout
func (p *Plugin) Run() {
	p.Log("notice", fmt.Sprintf("Start handler %s, v%s", p.Name, version))
	p.Connect()
	bg := p.QChan.Data.Join()
	/*dims := map[string]string{
		"version": version,
		"plugin": p.Name,
	}*/
	for {
		select {
		case val := <-bg.Read:
			switch val.(type) {
			case qtypes.Metric:
				m := val.(qtypes.Metric)
				if p.StopProcessingMetric(m, false) {
					continue
				}
			case qtypes_inventory.ContainerEvent:
				ce := val.(qtypes_inventory.ContainerEvent)
				p.PushToKafka(ce)
			}
		}
	}
}

func (p *Plugin) PushToKafka(ce qtypes_inventory.ContainerEvent) (err error) {
	value := "Hello Go!"
	topic := p.CfgStringOr("topic", "qframe")
	err = p.producer.Produce(&kafka.Message{TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny}, Value: []byte(value)}, p.deliveryChan)

	e := <-p.deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	} else {
		fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}
	return
}