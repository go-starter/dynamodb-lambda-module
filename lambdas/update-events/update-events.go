package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodb_types "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type EventService struct {
	ctx            context.Context
	logger         *log.Logger
	dynamodbClient DynamodbClient

	eventsDDB string
}

type DynamodbClient interface {
	UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
}

// Incoming events schema
type EventSchema struct {
	EventId      string `json:"eventId"`
	EventDetails string `json:"eventDetails"`
}

func main() {

	// Load config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Error loading AWS config: %s", err.Error())
	}
	// load common clients that can be served in multiple instances of lambda
	service := EventService{
		logger:         log.New(os.Stdout, "", log.LstdFlags),
		dynamodbClient: dynamodb.NewFromConfig(cfg),
		eventsDDB:      os.Getenv("EVENTS_DDB_TABLE"),
	}

	lambda.Start(service.handler)
}

func (svc *EventService) handler(ctx context.Context, event events.CloudWatchEvent) error {
	svc.ctx = ctx

	// Parse the event data
	var data EventSchema
	err := json.Unmarshal(event.Detail, &data)
	if err != nil {
		return err
	}

	// Update the item in the table
	_, err = svc.dynamodbClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(svc.eventsDDB),
		Key: map[string]dynamodb_types.AttributeValue{
			"eventId": &dynamodb_types.AttributeValueMemberS{Value: data.EventId},
		},
		UpdateExpression: aws.String("set eventDetails = :eventDetails"),
		ExpressionAttributeValues: map[string]dynamodb_types.AttributeValue{
			":eventDetails": &dynamodb_types.AttributeValueMemberS{Value: data.EventDetails},
		},
	})
	if err != nil {
		return err
	}

	svc.logger.Printf("Updated person with ID %s\n", data.EventId)
	return nil
}
