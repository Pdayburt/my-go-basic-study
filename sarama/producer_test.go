package sarama

import (
	"fmt"
	"github.com/IBM/sarama"
	"testing"
)

var addrs = []string{"localhost:9094"}

func TestSyncProducer(t *testing.T) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewHashPartitioner
	producer, err := sarama.NewSyncProducer(addrs, config)
	if err != nil {
		t.Fatal(err)
	}
	partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
		Topic: "test_topic",
		Key:   sarama.StringEncoder("oid"),
		Value: sarama.StringEncoder("hello kafka~2"),
		/*	Headers: []sarama.RecordHeader{
				{
					Key:   []byte("trace_id"),
					Value: []byte("123456"),
				},
			},
			Metadata: "这是metadata",*/
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(partition, offset)

}

func TestAsyncProducer(t *testing.T) {
	config := sarama.NewConfig()
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true
	producer, err := sarama.NewAsyncProducer(addrs, config)
	if err != nil {
		t.Fatal(err)
	}
	msgCh := producer.Input()
	msgCh <- &sarama.ProducerMessage{
		Topic: "test_topic",
		Key:   sarama.StringEncoder("oid"),
		Value: sarama.StringEncoder("hello kafka~3"),
	}
	errors := producer.Errors()
	successes := producer.Successes()

	select {
	case err := <-errors:
		t.Log("error:", err.Error())
	case success := <-successes:
		t.Log("success:", success.Value)
	}

}
