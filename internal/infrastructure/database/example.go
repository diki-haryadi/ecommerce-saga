package database

import (
	"context"
	"log"
	"time"

	as "github.com/aerospike/aerospike-client-go/v6"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func Example() {
	// Create a database configuration using the builder
	postgresConfig := NewConnectionBuilder().
		WithHost("localhost").
		WithPort(5432).
		WithCredentials("user", "pass").
		WithDatabase("myapp").
		WithMaxConnections(20).
		WithConnectTimeout(5*time.Second).
		WithMaxIdleTime(10*time.Minute).
		WithRetryPolicy(3, time.Second).
		WithSSLMode("disable").
		Build()

	// Get the database manager instance
	manager := GetInstance()

	// Add PostgreSQL connection
	if err := manager.AddConnection("postgres_main", TypePostgres, postgresConfig); err != nil {
		log.Fatalf("Error adding PostgreSQL connection: %v", err)
	}

	// Create MySQL configuration
	mysqlConfig := NewConnectionBuilder().
		WithHost("localhost").
		WithPort(3306).
		WithCredentials("user", "pass").
		WithDatabase("myapp").
		WithMaxConnections(20).
		WithOption("parseTime", "true").
		WithOption("multiStatements", "true").
		Build()

	// Add MySQL connection
	if err := manager.AddConnection("mysql_main", TypeMySQL, mysqlConfig); err != nil {
		log.Fatalf("Error adding MySQL connection: %v", err)
	}

	// Create Oracle configuration
	oracleConfig := NewConnectionBuilder().
		WithHost("localhost").
		WithPort(1521).
		WithCredentials("system", "oracle").
		WithDatabase("ORCLPDB1").
		WithMaxConnections(20).
		Build()

	// Add Oracle connection
	if err := manager.AddConnection("oracle_main", TypeOracle, oracleConfig); err != nil {
		log.Fatalf("Error adding Oracle connection: %v", err)
	}

	// Create MongoDB configuration
	mongoConfig := NewConnectionBuilder().
		WithHost("localhost").
		WithPort(27017).
		WithCredentials("user", "pass").
		WithDatabase("myapp").
		WithMaxConnections(20).
		Build()

	// Add MongoDB connection
	if err := manager.AddConnection("mongodb_main", TypeMongoDB, mongoConfig); err != nil {
		log.Fatalf("Error adding MongoDB connection: %v", err)
	}

	// Create DynamoDB configuration
	dynamoConfig := NewConnectionBuilder().
		WithHost("localhost"). // For local DynamoDB
		WithPort(8000).
		WithOption("region", "us-west-2").
		Build()

	// Add DynamoDB connection
	if err := manager.AddConnection("dynamodb_main", TypeDynamoDB, dynamoConfig); err != nil {
		log.Fatalf("Error adding DynamoDB connection: %v", err)
	}

	// Create Cassandra configuration
	cassandraConfig := NewConnectionBuilder().
		WithHost("localhost").
		WithPort(9042).
		WithCredentials("cassandra", "cassandra").
		WithDatabase("mykeyspace").
		WithMaxConnections(20).
		WithOption("consistency", "quorum").
		WithOption("retry_policy", "exponential").
		Build()

	// Add Cassandra connection
	if err := manager.AddConnection("cassandra_main", TypeCassandra, cassandraConfig); err != nil {
		log.Fatalf("Error adding Cassandra connection: %v", err)
	}

	// Create Aerospike configuration
	aerospikeConfig := NewConnectionBuilder().
		WithHost("localhost").
		WithPort(3000).
		WithCredentials("admin", "admin123").
		WithMaxConnections(20).
		WithMaxIdleTime(30 * time.Minute).
		Build()

	// Add Aerospike connection
	if err := manager.AddConnection("aerospike_main", TypeAerospike, aerospikeConfig); err != nil {
		log.Fatalf("Error adding Aerospike connection: %v", err)
	}

	// Example usage for each database type
	ctx := context.Background()

	// PostgreSQL Example
	postgresDB, _ := manager.GetConnection("postgres_main")
	if _, err := postgresDB.Execute(ctx, "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", "John Doe", "john@example.com"); err != nil {
		log.Printf("PostgreSQL error: %v", err)
	}

	// MySQL Example
	mysqlDB, _ := manager.GetConnection("mysql_main")
	if _, err := mysqlDB.Execute(ctx, "INSERT INTO users (name, email) VALUES (?, ?)", "Jane Doe", "jane@example.com"); err != nil {
		log.Printf("MySQL error: %v", err)
	}

	// Oracle Example
	oracleDB, _ := manager.GetConnection("oracle_main")
	if _, err := oracleDB.Execute(ctx, "INSERT INTO users (name, email) VALUES (:1, :2)", "Bob Smith", "bob@example.com"); err != nil {
		log.Printf("Oracle error: %v", err)
	}

	// MongoDB Example
	mongoDB, _ := manager.GetConnection("mongodb_main")
	doc := map[string]interface{}{
		"name":  "Alice Smith",
		"email": "alice@example.com",
	}
	if _, err := mongoDB.Execute(ctx, "users", "insert", doc); err != nil {
		log.Printf("MongoDB error: %v", err)
	}

	// DynamoDB Example
	dynamoDB, _ := manager.GetConnection("dynamodb_main")
	putItem := &types.PutItemInput{
		TableName: &[]string{"users"}[0],
		Item: map[string]types.AttributeValue{
			"id":    &types.AttributeValueMemberS{Value: "1"},
			"name":  &types.AttributeValueMemberS{Value: "Charlie Brown"},
			"email": &types.AttributeValueMemberS{Value: "charlie@example.com"},
		},
	}
	if _, err := dynamoDB.Execute(ctx, "PutItem", putItem); err != nil {
		log.Printf("DynamoDB error: %v", err)
	}

	// Cassandra Example
	cassandraDB, _ := manager.GetConnection("cassandra_main")
	if _, err := cassandraDB.Execute(ctx, "INSERT INTO users (id, name, email) VALUES (uuid(), ?, ?)", "David Wilson", "david@example.com"); err != nil {
		log.Printf("Cassandra error: %v", err)
	}

	// Aerospike Example
	aerospikeDB, _ := manager.GetConnection("aerospike_main")
	bins := as.BinMap{
		"name":  "Eve Johnson",
		"email": "eve@example.com",
	}
	if _, err := aerospikeDB.Execute(ctx, "Put", "test", "users", "user1", bins); err != nil {
		log.Printf("Aerospike error: %v", err)
	}

	// Close all connections when done
	if err := manager.CloseAll(); err != nil {
		log.Printf("Error closing connections: %v", err)
	}
}
