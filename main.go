package main

import (
	"context"
	"log"
	"os"

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

func HandleLambdaEvent(ctx context.Context, getContactRequest ContactRequest) (User, error) {
	cfg, err := config.LoadDefaultConfig(ctx, func(opts *config.LoadOptions) error {
		opts.Region = os.Getenv("AWS_REGION")
		return nil
	})
	if err != nil {
		log.Printf("error loading dynamo configuration: %v", err)
		return User{}, err
	}
	svc := dynamodb.NewFromConfig(cfg)
	userRetrieved, err := retrieveContact(ctx, svc, getContactRequest.ID)
	if err != nil {
		log.Printf("error getting user information: %v", err)
		return User{}, err
	}
	return userRetrieved, nil
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
