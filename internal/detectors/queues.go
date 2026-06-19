package detectors

import (
	"path/filepath"
	"strings"

	"github.com/RevIneX/context/internal/data"
)

var queueSignatures = []struct {
	Sig  string
	Tech string
}{
	{"kafka", "Kafka"},
	{"rabbitmq", "RabbitMQ"},
	{"amqplib", "RabbitMQ"},
	{"amqp", "RabbitMQ"},
	{"sqs", "SQS"},
	{"aws-sdk", "SQS"},
	{"nats", "NATS"},
	{"stan", "NATS"},
	{"pulsar", "Pulsar"},
	{"activemq", "ActiveMQ"},
	{"zeromq", "ZeroMQ"},
	{"zmq", "ZeroMQ"},
	{"mqtt", "MQTT"},
	{"redpanda", "Redpanda"},
}

var queueFiles = map[string]string{
	"activemq.xml": "ActiveMQ",
}

var queueAllowedExts = []string{
	".yaml", ".yml", ".conf", ".xml", ".json", ".properties",
}

type QueueDetector struct{}

func (d *QueueDetector) Name() string {
	return "queue"
}

func (d *QueueDetector) Detect(filePath string) []data.Finding {
	fileName := filepath.Base(filePath)
	ext := filepath.Ext(fileName)

	allowed := false
	for _, allowedExt := range queueAllowedExts {
		if strings.EqualFold(ext, allowedExt) {
			allowed = true
			break
		}
	}
	if !allowed {
		return nil
	}

	lowName := strings.ToLower(fileName)

	if tech, ok := queueFiles[lowName]; ok {
		return []data.Finding{
			data.NewFinding(data.CatQueue, tech, "", filePath, 1, nil),
		}
	}

	for i := range queueSignatures {
		if strings.Contains(lowName, queueSignatures[i].Sig) {
			return []data.Finding{
				data.NewFinding(data.CatQueue, queueSignatures[i].Tech, "", filePath, 1, nil),
			}
		}
	}

	return nil
}
