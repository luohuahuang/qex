package kafka

import (
	"log"
	"time"

	"github.com/Shopify/sarama"
	"github.com/wvanbergen/kafka/consumergroup"
)

type Consumer struct {
	ZookeeperConn string
	Topic         string
	ConsumerGroup *consumergroup.ConsumerGroup
}

func New(topic string, zookeeper string) *Consumer {
	k := &Consumer{
		ZookeeperConn: zookeeper,
		Topic:         topic,
		ConsumerGroup: nil,
	}

	cgroup := "zgroup"
	config := consumergroup.NewConfig()
	config.Offsets.Initial = sarama.OffsetNewest
	config.Offsets.ProcessingTimeout = 3 * time.Second

	cg, err := consumergroup.JoinConsumerGroup(cgroup, []string{k.Topic}, []string{k.ZookeeperConn}, config)
	if err != nil {
		log.Panic(err.Error())
	}
	k.ConsumerGroup = cg

	return k
}
