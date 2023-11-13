package ratingingester

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	kafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"movieexample.com/rating/pkg/model"
)

func main() {
	fmt.Println("Producing a Kafka producer")
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	const filename = "ratingsdata.json"
	fmt.Println("Reading rating events from file " + filename)
	ratingEvents, err := readRatingEvents(filename)
	if err != nil {
		panic(err)
	}
	const topic = "ratings"

	if err := produceRatingEvents(topic, producer, ratingEvents); err != nil {
		panic(err)
	}
	const timeout = time.Second * 10
	fmt.Println("Waiting " + timeout.String() + " until all events get produced.")
	producer.Flush(int(timeout.Milliseconds()))
}

func readRatingEvents(fileName string) ([]model.RatingEvent, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var ratings []model.RatingEvent
	if err := json.NewDecoder(f).Decode(&ratings); err != nil {
		return nil, err
	}
	return ratings, nil
}

func produceRatingEvents(topic string, producer *kafka.Producer, events []model.RatingEvent) error {
	for _, event := range events {
		encodedEvent, err := json.Marshal(event)
		if err != nil {
			return err
		}
		if err := producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(encodedEvent),
		}, nil); err != nil {
			return err
		}
	}
	return nil
}
