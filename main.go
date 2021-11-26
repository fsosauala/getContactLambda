package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	tableName = "contactsFredy"
)

type (
	ContactRequest struct {
		ID string `json:"id"`
	}
	User struct {
		ID        string `dynamodbav:"id" json:"id"`
		Status    string `dynamodbav:"status" json:"status"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}
)

func HandleLambdaEvent(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := req.PathParameters["id"]
	cfg, err := config.LoadDefaultConfig(ctx, func(opts *config.LoadOptions) error {
		opts.Region = os.Getenv("AWS_REGION")
		return nil
	})
	if err != nil {
		log.Printf("error loading dynamo configuration: %v", err)
		return events.APIGatewayProxyResponse{}, err
	}
	svc := dynamodb.NewFromConfig(cfg)
	userRetrieved, err := retrieveContact(ctx, svc, id)
	if err != nil {
		log.Printf("error getting user information: %v", err)
		return events.APIGatewayProxyResponse{}, err
	}
	data, err := json.Marshal(userRetrieved)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	return events.APIGatewayProxyResponse{
		StatusCode:      http.StatusOK,
		Body:            string(data),
		IsBase64Encoded: false,
	}, nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}

func retrieveContact(ctx context.Context, svc *dynamodb.Client, userID string) (User, error) {
	out, err := svc.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: userID,
			},
		},
	})
	if err != nil {
		return User{}, err
	}
	var u User
	err = attributevalue.UnmarshalMap(out.Item, &u)
	if err != nil {
		return User{}, err
	}
	return u, nil
}
