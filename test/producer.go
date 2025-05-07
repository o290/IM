package main

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

func SendTopic(topic string, msg []byte) {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP("127.0.0.1:9092"),
		Topic:                  topic,
		Balancer:               &kafka.Hash{},
		WriteTimeout:           10 * time.Second,
		RequiredAcks:           kafka.RequireNone,
		AllowAutoTopicCreation: true,
	}
	defer writer.Close()

	err := writer.WriteMessages(
		context.Background(),
		kafka.Message{Value: msg},
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("topic:%s 消息发送成功 \n", topic)

}

func main() {
	fmt.Println("111")
	SendTopic("test_topic", []byte("枫枫"))
	SendTopic("test_topic", []byte("知道"))
}
