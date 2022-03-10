package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/luohuahuang/qex/monitor"
	"log"
	"os"
)

func InitProducer(kafkaConn string) (sarama.SyncProducer, error) {
	sarama.Logger = log.New(os.Stdout, "", log.Ltime)

	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.NoResponse
	config.Producer.Return.Successes = true
	config.ClientID = "randomclientid"

	prd, err := sarama.NewSyncProducer([]string{kafkaConn}, config)
	return prd, err
}

func Send(producer sarama.SyncProducer, topic string, data string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(data),
	}
	_, _, err := producer.SendMessage(msg)
	if err != nil {
		monitor.SendAlert(err)
		return err
	}
	return nil
}
