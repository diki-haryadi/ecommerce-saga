package dynamodb

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/database"
)

type dynamoDB struct {
	config *database.Config
	client *dynamodb.Client
	mu     sync.RWMutex
}

// NewDynamoDB creates a new DynamoDB instance
func NewDynamoDB(config *database.Config) database.Database {
	return &dynamoDB{
		config: config,
	}
}

func (db *dynamoDB) Connect(ctx context.Context) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.client != nil {
		return nil
	}

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(db.config.AdditionalOpts["region"]),
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				if db.config.Host != "" {
					return aws.Endpoint{
						URL: fmt.Sprintf("http://%s:%d", db.config.Host, db.config.Port),
					}, nil
				}
				return aws.Endpoint{}, &aws.EndpointNotFoundError{}
			},
		)),
	)
	if err != nil {
		return fmt.Errorf("error loading aws config: %w", err)
	}

	// Create DynamoDB client
	db.client = dynamodb.NewFromConfig(cfg)
	return nil
}

func (db *dynamoDB) Disconnect(ctx context.Context) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.client = nil
	return nil
}

func (db *dynamoDB) Ping(ctx context.Context) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if db.client == nil {
		return fmt.Errorf("database not connected")
	}

	// Try to list tables as a ping mechanism
	_, err := db.client.ListTables(ctx, &dynamodb.ListTablesInput{
		Limit: aws.Int32(1),
	})
	return err
}

func (db *dynamoDB) IsConnected() bool {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.client != nil
}

func (db *dynamoDB) GetConfig() *database.Config {
	return db.config
}

func (db *dynamoDB) Begin(ctx context.Context) error {
	// DynamoDB uses implicit transactions
	return nil
}

func (db *dynamoDB) Commit(ctx context.Context) error {
	// DynamoDB uses implicit transactions
	return nil
}

func (db *dynamoDB) Rollback(ctx context.Context) error {
	// DynamoDB uses implicit transactions
	return nil
}

func (db *dynamoDB) Execute(ctx context.Context, operation string, args ...interface{}) (interface{}, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if db.client == nil {
		return nil, fmt.Errorf("database not connected")
	}

	if len(args) == 0 {
		return nil, fmt.Errorf("no arguments provided for operation")
	}

	switch operation {
	case "PutItem":
		input, ok := args[0].(*dynamodb.PutItemInput)
		if !ok {
			return nil, fmt.Errorf("invalid input type for PutItem")
		}
		return db.client.PutItem(ctx, input)

	case "GetItem":
		input, ok := args[0].(*dynamodb.GetItemInput)
		if !ok {
			return nil, fmt.Errorf("invalid input type for GetItem")
		}
		return db.client.GetItem(ctx, input)

	case "UpdateItem":
		input, ok := args[0].(*dynamodb.UpdateItemInput)
		if !ok {
			return nil, fmt.Errorf("invalid input type for UpdateItem")
		}
		return db.client.UpdateItem(ctx, input)

	case "DeleteItem":
		input, ok := args[0].(*dynamodb.DeleteItemInput)
		if !ok {
			return nil, fmt.Errorf("invalid input type for DeleteItem")
		}
		return db.client.DeleteItem(ctx, input)

	case "Query":
		input, ok := args[0].(*dynamodb.QueryInput)
		if !ok {
			return nil, fmt.Errorf("invalid input type for Query")
		}
		return db.client.Query(ctx, input)

	case "Scan":
		input, ok := args[0].(*dynamodb.ScanInput)
		if !ok {
			return nil, fmt.Errorf("invalid input type for Scan")
		}
		return db.client.Scan(ctx, input)

	default:
		return nil, fmt.Errorf("unsupported operation: %s", operation)
	}
}

func (db *dynamoDB) Query(ctx context.Context, tableName string, args ...interface{}) (interface{}, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if db.client == nil {
		return nil, fmt.Errorf("database not connected")
	}

	if len(args) < 2 {
		return nil, fmt.Errorf("insufficient arguments for query")
	}

	// Extract key condition expression and attribute values
	keyCondition, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("key condition must be a string")
	}

	attrValues, ok := args[1].(map[string]types.AttributeValue)
	if !ok {
		return nil, fmt.Errorf("attribute values must be a map[string]types.AttributeValue")
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		KeyConditionExpression:    aws.String(keyCondition),
		ExpressionAttributeValues: attrValues,
	}

	// Add optional parameters if provided
	if len(args) > 2 {
		if filterExpr, ok := args[2].(string); ok {
			input.FilterExpression = aws.String(filterExpr)
		}
	}

	return db.client.Query(ctx, input)
}
