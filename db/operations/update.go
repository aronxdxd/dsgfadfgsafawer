package operations

import (
	"context"
	"project/db"

	"go.mongodb.org/mongo-driver/bson"
)

func UpdateUserField(ctx context.Context, telegramID int, field string, value interface{}) error {
	client := db.GetClient()
	collection := client.Database("db").Collection("users")
	filter := bson.M{"user_id": telegramID}
	update := bson.M{"$set": bson.M{field: value}}

	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

func IncrementUserField(ctx context.Context, telegramID int, field string, increment interface{}) error {
	client := db.GetClient()
	collection := client.Database("db").Collection("users")
	filter := bson.M{"user_id": telegramID}
	update := bson.M{"$inc": bson.M{field: increment}}

	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}
