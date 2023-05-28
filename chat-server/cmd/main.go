package main

import (
	// "errors"
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	// "time"

	"github.com/Shopify/sarama"
	"github.com/segmentio/kafka-go"

	kafka_service "soa.mafia-game/kafka-help"
)

var chats map[string]bool
var workers map[string]func(string)
var mut sync.Mutex
var chatLeft = make(map[string]int)

func messagesHandler(topic string) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// указываем адреса брокеров Kafka
	brokers := []string{"kafka1:19092"}

	// // создаем консьюмера
	// consumer, err := sarama.NewConsumer(brokers, config)
	// if err != nil {
	//     panic(err)
	// }
	// defer func() {
	//     if err := consumer.Close(); err != nil {
	//         panic(err)
	//     }
	// }()
	// // создаем партиционный консьюмер
	// consumerPartition, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	// if err != nil {
	//     // panic(err)
	//     return
	// }
	// defer func() {
	//     if err := consumerPartition.Close(); err != nil {
	//         panic(err)
	//     }
	// }()

	admin, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		// panic(err)
		return
	}
	defer func() {
		if err := admin.Close(); err != nil {
			panic(err)
		}
	}()

	producer, err := kafka_service.GetNewProducer("kafka1:19092")
	if err != nil {
		log.Printf("Failed to create producer %v", err)
		return
	}
    defer producer.Close()
	// читаем сообщения из топика
	maxOffset := 100000 // bad, but still =(

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   brokers,
		Topic:     topic,
		Partition: int(0),
	})
	defer reader.Close()
	reader.SetOffset(0)
	log.Printf("Listen to %v", topic)
	defer log.Printf("Stop listening to %v", topic)
	for {
		msg, _ := reader.ReadMessage(context.Background())

		chatId := string(msg.Key)
		user_login := topic
		message := string(msg.Value)

		if message == "exit" {
			// delete from maps
			// delete topic from kafka
			_ = admin.DeleteTopic(user_login)
			// _ = admin.DeleteTopic(chatId)
			partOffsets := map[int32]int64{}
			for i := 0; i < 4; i++ {
				partOffsets[int32(i)] = int64(maxOffset)
			}
			mut.Lock()
			chatLeft[chatId] += 1
			log.Printf("Exit %v already exited: %v", user_login, chatLeft[chatId])
			if chatLeft[chatId] == 4 {
				err := admin.DeleteTopic(chatId) // NOT WORKIN WELL
				if err != nil {
					log.Printf("Failed to delete topic: %v", err)
				}
			}
			delete(chats, chatId)
			delete(workers, topic)
			mut.Unlock()
			return
		}
		mut.Lock()
		if _, exists := chats[chatId]; !exists {
			_ = admin.CreateTopic(chatId, &sarama.TopicDetail{NumPartitions: 4, ReplicationFactor: 1}, false)
			// if err != nil && !errors.Is(sarama.ErrTopicAlreadyExists, err) {
			//     log.Fatalf("Failed to open chat: %v; %v", err, errors.Is(sarama.ErrTopicAlreadyExists, err))
			// }
		}
		for i := 0; i < 4; i++ {
			kafka_service.Produce(user_login, message, chatId, int32(i), producer)
		}
		mut.Unlock()
	}
}
func process() {
	chats := make(map[string]bool)
	config := sarama.NewConfig()

	kafkaBrokers := []string{"kafka1:19092"} // from env

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
	workers = make(map[string]func(string))
	for {
		topics, err := admin.ListTopics()
		if err != nil {
			panic(err)
		}
		// log.Printf("!!!! all topics: %v", topics)
		for topic := range topics {
			_, exist := workers[topic]
			if len(topic) >= len(chat_prefix) && topic[:len(chat_prefix)] == chat_prefix {
				chats[topic] = true
				continue
			}
			if !exist {
				log.Printf("Hi %v, open connection to chat!", topic)
				workers[topic] = messagesHandler
				go workers[topic](topic)
			}
		}
		// time.Sleep(10 * time.Millisecond)
	}
}
func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)
	go process()
	<-stop
	log.Printf("Shuting down")
	// config := sarama.NewConfig()

	// kafkaBrokers := []string{"kafka1:19092"} // from env

	// admin, _ := sarama.NewClusterAdmin(kafkaBrokers, config)
	// for topics := range chats {
	//     admin.DeleteRecords()
	// }
}
