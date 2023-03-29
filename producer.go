package vaskafka

import (
    "fmt"
    "strings"

    "github.com/confluentinc/confluent-kafka-go/kafka"
    "github.com/google/uuid"
)

func Produce(eventName, email string) error {
    config, err := NewConfig(eventName)
    if err != nil {
        return err
    }
    emailTrimmed := strings.TrimSpace(email)
    if emailTrimmed == "" {
        panic("singer name must not be blank")
    }

    p, err := kafka.NewProducer(&config.Kafka)
    if err != nil {
        return err
    }

    // Go-routine to handle message delivery reports and
    // possibly other event types (errors, stats, etc)
    go func() {
        for e := range p.Events() {
            switch ev := e.(type) {
            case *kafka.Message:
                if ev.TopicPartition.Error != nil {
                    fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
                } else {
                    fmt.Printf("Produced event to topic %s: key = %-10s value = %s\n",
                        *ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
                }
            }
        }
    }()

    err = p.Produce(&kafka.Message{
        TopicPartition: kafka.TopicPartition{
            Topic: &config.User.CompetitionName, Partition: kafka.PartitionAny,
        },
        Key:   []byte(emailTrimmed),
        Value: []byte(uuid.New().String()),
    }, nil)
    if err != nil {
        return err
    }

    const flushDelay = 15 * 1000
    // Wait for all messages to be delivered
    p.Flush(flushDelay)
    p.Close()

    return nil
}
