package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Quote struct {
	Name        string `json:"name"`
	CatchPhrase string `json:"catchPhrase"`
}

func main() {

	sess := session.Must(session.NewSession())
	creds := stscreds.NewCredentials(sess, "arn:aws:iam::123456789012:role/user-dynamo-db-assume")
	dynamo := dynamodb.New(sess, &aws.Config{Credentials: creds, Region: aws.String("eu-central-1")})

	tableName := "quotes"

	params := &dynamodb.ScanInput{
		TableName: &tableName,
	}

	dynamo_result, err := dynamo.Scan(params)
	if err != nil {
		fmt.Println(err.Error())
	}

	var quotes []Quote

	for _, i := range dynamo_result.Items {
		quote := Quote{}

		err = dynamodbattribute.UnmarshalMap(i, &quote)
		if err != nil {
			fmt.Println("Got error unmarshalling:")
			fmt.Println(err)
		}
		quotes = append(quotes, quote)

	}

	data, _ := json.Marshal(quotes)
	fmt.Println(string(data))
}
