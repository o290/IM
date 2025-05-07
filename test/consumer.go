package main

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

func TopicEvent(topic string) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{"127.0.0.1:9092"},
		Topic:          topic,
		CommitInterval: 1 * time.Second,
		GroupID:        "rec_team",
		StartOffset:    kafka.FirstOffset,
	})
	fmt.Printf("topic:%s kafka消费者正在运行中\n", topic)
	for {
		message, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("读取kafka失败: ", err)
			break
		}
		fmt.Printf("topic=%s, value=%s \n", message.Topic, string(message.Value))
	}
}
func main() {
	go TopicEvent("test_topic")
	select {}
}
