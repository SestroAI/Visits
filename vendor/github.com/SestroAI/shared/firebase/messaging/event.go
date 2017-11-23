package messaging

import "time"

import(
		"golang.org/x/net/context"
		"cloud.google.com/go/pubsub"

		"github.com/SestroAI/shared/logger"
)


func GetSubscription(topicName string, subscriptionName string) (*pubsub.Subscription, error){
	client, err := GetPubSubClient()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	//Create Topic if not exists
	topic := client.Topic(topicName)
	ok, err := topic.Exists(ctx)
	if err != nil {
		return nil, err
	}
	if !ok{
		//topic does not exists
		topic, err = client.CreateTopic(ctx, topicName)
		if err != nil {
			logger.Infof(err.Error())
			return nil, err
		}
	}

	//Get or Create Subscrption
	sub := client.Subscription(subscriptionName)
	ok, err = sub.Exists(ctx)
	if err != nil {
		return nil, err
	}
	if !ok {
		sub, err = client.CreateSubscription(
			ctx,
			subscriptionName,
			pubsub.SubscriptionConfig{
				Topic: topic, //client.Topic(VisitEndEventname),
			})
		sub.ReceiveSettings.MaxExtension = 30 * time.Second
	}
	return sub, err
}

func ListenEvent(topicName string, subscriptionName string, callbackFunc func(context.Context, *pubsub.Message)){
	sub, err := GetSubscription(topicName, subscriptionName)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	err = sub.Receive(ctx, callbackFunc)
	if err != nil {
		logger.Errorf("Unable to receive messages for visit end error = %s", err.Error())
		panic(err)
	}
}