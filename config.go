package vaskafka

import (
    _ "embed"
    "fmt"
    "strings"

    "github.com/confluentinc/confluent-kafka-go/kafka"
)

type Config struct {
    Kafka kafka.ConfigMap
    User  UserConfigType
}

type UserConfigType struct {
    CompetitionName string
    GroupID         string
}

//go:embed kafka.properties
var kafkaProperties []byte

func NewConfig(competitionName string) (Config, error) {
    kconfig := make(map[string]kafka.ConfigValue)

    lines := strings.Split(string(kafkaProperties), "\n")
    for _, line := range lines {
        if !strings.HasPrefix(line, "#") && strings.TrimSpace(line) != "" {
            kv := strings.Split(line, "=")
            parameter := strings.TrimSpace(kv[0])
            value := strings.TrimSpace(kv[1])
            kconfig[parameter] = value
        }
    }
    // kconfig["bootstrap.servers"] = "localhost:9092"

    competitionNameTrimmed := strings.TrimSpace(competitionName)
    var empty Config
    if competitionNameTrimmed == "" {
        return empty, fmt.Errorf("competition name must not be blank")
    }

    return Config{
        Kafka: kconfig,
        User: UserConfigType{
            CompetitionName: competitionNameTrimmed,
            GroupID:         "masked-singer",
        },
    }, nil
}
