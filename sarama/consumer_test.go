package sarama

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"log"
	"testing"
)

func TestConsumer(t *testing.T) {

	cfg := sarama.NewConfig()
	consumerGroup, err := sarama.NewConsumerGroup(addrs,
		"test_group", cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer consumerGroup.Close()

	consumerGroup.Consume(context.Background(),
		[]string{"test_topic"}, testConsumerGroupHandler{})

}

type testConsumerGroupHandler struct {
}

func (t testConsumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	//TODO implement me
	panic("implement me")
}

func (t testConsumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	//TODO implement me
	panic("implement me")
}

func (t testConsumerGroupHandler) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {

	messages := claim.Messages()
	for message := range messages {
		var bizMsg MyBizMsg
		err := json.Unmarshal(message.Value, &bizMsg)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println(bizMsg)
		session.MarkMessage(message, "")
	}
	return nil
}

type MyBizMsg struct {
	Name string
}
