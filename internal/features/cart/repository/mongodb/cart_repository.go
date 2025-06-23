package mongodb

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/cart/domain/entity"
)

const (
	cartCollection = "carts"
)

// CartRepository implements the repository.CartRepository interface for MongoDB
type CartRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewCartRepository creates a new MongoDB cart repository
func NewCartRepository(db *mongo.Database) *CartRepository {
	return &CartRepository{
		db:         db,
		collection: db.Collection(cartCollection),
	}
}

// Create saves a new cart to the database
func (r *CartRepository) Create(ctx context.Context, cart *entity.Cart) error {
	_, err := r.collection.InsertOne(ctx, cart)
	return err
}

// GetByID retrieves a cart by its ID
func (r *CartRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Cart, error) {
	var cart entity.Cart
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&cart)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

// GetByUserID retrieves a cart by user ID
func (r *CartRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.Cart, error) {
	var cart entity.Cart
	err := r.collection.FindOne(ctx, bson.M{
		"user_id":    userID,
		"expires_at": bson.M{"$gt": time.Now()},
	}).Decode(&cart)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

// Update updates an existing cart in the database
func (r *CartRepository) Update(ctx context.Context, cart *entity.Cart) error {
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": cart.ID}, cart)
	return err
}

// Delete removes a cart from the database
func (r *CartRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// DeleteExpired removes all expired carts from the database
func (r *CartRepository) DeleteExpired(ctx context.Context) error {
	_, err := r.collection.DeleteMany(ctx, bson.M{
		"expires_at": bson.M{"$lt": time.Now()},
	})
	return err
}

// Setup creates necessary indexes for the cart collection
func (r *CartRepository) Setup(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "expires_at", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "expires_at", Value: 1}},
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}
