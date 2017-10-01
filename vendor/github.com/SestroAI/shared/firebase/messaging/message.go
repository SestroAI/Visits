package messaging

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/SestroAI/shared/config"
	"time"
)

func GetPubSubClient() (*pubsub.Client, error) {
	return pubsub.NewClient(context.Background(), config.GetGoogleProjectID())
}

func sendPubsubMessage(topicName string, data interface{}) error {
	ctx := context.Background()
	client, err := GetPubSubClient()

	//Create topic if doesn't exist
	topic, _ := client.CreateTopic(ctx, topicName)

	message := &pubsub.Message{
		Data:[]byte(data),
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