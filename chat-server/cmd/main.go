package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/Shopify/sarama"
	"github.com/segmentio/kafka-go"

	kafka_service "soa.mafia-game/kafka-help"
)

var chats map[string]bool
var workers map[string]func(topic string, kafkaBrokers []string)
var mut sync.Mutex
var chatLeft = make(map[string]int)

func messagesHandler(topic string, brokers []string) {
	log.Printf("Listen to %v", topic)
	defer log.Printf("Stop listening to %v", topic)

	config := sarama.NewConfig()
	admin, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		log.Printf("Failed to create cluster admin %v", err)
		return
	}
	defer admin.Close()

	producer, err := kafka_service.GetNewProducer("kafka1:19092")
	if err != nil {
		log.Printf("Failed to create producer %v", err)
		return
	}
	defer producer.Close()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   brokers,
		Topic:     topic,
		Partition: int(0),
	})
	defer reader.Close()
	reader.SetOffset(0)
	for {
		msg, _ := reader.ReadMessage(context.Background())

		chatId := string(msg.Key)
		user_login := topic
		message := string(msg.Value)

		if message == "exit" {
			err = admin.DeleteTopic(user_login)
			if err != nil {
				log.Printf("Failed to delete topic %s: %v", user_login, err)
			}
			mut.Lock()
			chatLeft[chatId] += 1
			log.Printf("Exit %v already exited: %v", user_login, chatLeft[chatId])
			if chatLeft[chatId] == 4 {
				err = admin.DeleteTopic(chatId)
				if err != nil {
					log.Printf("Failed to delete topic %s: %v", chatId, err)
				}
			}
			delete(chats, chatId)
			delete(workers, topic)
			mut.Unlock()
			return
		}
		mut.Lock()
		if _, exists := chats[chatId]; !exists {
			err = admin.CreateTopic(chatId, &sarama.TopicDetail{NumPartitions: 4, ReplicationFactor: 1}, false)
			if err != nil {
				log.Printf("Failed to create topic %s: %v", chatId, err)
			}
		}
		for i := 0; i < 4; i++ {
			err = kafka_service.Produce(user_login, message, chatId, int32(i), producer)
			if err != nil {
				log.Printf("Failed to produce %s to %s: %v", message, chatId, err)
			}
		}
		mut.Unlock()
	}
}

func process() {
	chats := make(map[string]bool)
	kafkaServers := os.Getenv("KAFKA_BROKER_URL")
	kafkaBrokers := strings.Split(kafkaServers, ",")

	config := sarama.NewConfig()
	admin, err := sarama.NewClusterAdmin(kafkaBrokers, config)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := admin.Close(); err != nil {
			panic(err)
		}
	}()
	log.Print("Chat Server successfuly started")
	chat_prefix := "chat-"
	workers = make(map[string]func(string, []string))
	for {
		topics, err := admin.ListTopics()
		if err != nil {
			panic(err)
		}
		for topic := range topics {
			_, exist := workers[topic]
			if len(topic) >= len(chat_prefix) && topic[:len(chat_prefix)] == chat_prefix {
				chats[topic] = true
				continue
			}
			if !exist {
				log.Printf("Hi %v, open connection to chat!", topic)
				workers[topic] = messagesHandler
				go workers[topic](topic, kafkaBrokers)
			}
		}
	}
}

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)
	go process()
	<-stop
	log.Printf("Shuting down")
}
