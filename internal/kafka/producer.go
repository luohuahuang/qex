package kafka

import (
	"github.com/Shopify/sarama"
	"log"
	"os"
)

func InitProducer(kafkaConn string) (sarama.SyncProducer, error) {
	// setup sarama log to stdout
	sarama.Logger = log.New(os.Stdout, "", log.Ltime)

	// producer config
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.NoResponse
	config.Producer.Return.Successes = true
	config.ClientID = "whosyourdaddy"

	// sync producer
	prd, err := sarama.NewSyncProducer([]string{kafkaConn}, config)
	log.Printf("connecting to kafka server: %s", kafkaConn)
	return prd, err
}

// SendWithKey with partition key
func Send(producer sarama.SyncProducer, topic string, data string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		// Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(data),
	}
	log.Printf("sending msg to kafka Topic: %s", topic)
	_, _, err := producer.SendMessage(msg)
	if err != nil {
		log.Println("error publish: ", err.Error())
		return err
	}
	//log.Println("Partition: ", p)
	//log.Println("Offset: ", o)
	log.Println("published msg to kafka for consuming")
	return nil
}
