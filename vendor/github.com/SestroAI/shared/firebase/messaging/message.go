package messaging

import (
	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
	"context"
	"github.com/SestroAI/shared/config"
	"time"
	"bytes"
	"encoding/gob"
)

func GetPubSubClient() (*pubsub.Client, error) {
	return pubsub.NewClient(
		context.Background(),
		config.GetGoogleProjectID(),
		option.WithCredentialsFile(config.ServiceAccountKeyPath),
	)
}

func getBytes(data interface{}) ([]byte, error){
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func sendPubsubMessage(topicName string, data interface{}) error {
	ctx := context.Background()
	client, err := GetPubSubClient()

	//Create topic if doesn't exist
	topic, _ := client.CreateTopic(ctx, topicName)

	messageData, err := getBytes(data)
	if err != nil {
		return err
	}
	message := &pubsub.Message{
		Data:messageData,
	}

	_, err = topic.Publish(ctx, message).Get(ctx)
	return err
}

type EventMessage struct {
	Timestamp time.Time `json:"timestamp"`
	Data interface{} `json:"data"`
	Meta *EventMeta `json:"meta"`
	TopicName string `json:"topicName"`
}

type EventMeta struct {
	GoogleProjectId string `json:"googleProjectId"`
	ServiceName string `json:"serviceName"`
}

func SendMessage(topicName string, data interface{}) error {
	event := EventMessage{
		Timestamp:time.Now(),
		Meta:&EventMeta{
			GoogleProjectId:config.GetGoogleProjectID(),
			ServiceName:config.ServiceName,
		},
		Data: data,
		TopicName:topicName,
	}

	return sendPubsubMessage(topicName, event)
}