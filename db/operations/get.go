package operations

import (
	"context"
	"project/db"

	"go.mongodb.org/mongo-driver/bson"
)

func Get(ctx context.Context, telegramID int, filter bson.M) (map[string]interface{}, error) {
	client := db.GetClient()
	collection := client.Database("db").Collection("users")
	
	var result map[string]interface{}
	
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	
	return result, nil
}

func GetAllUsers(ctx context.Context) ([]map[string]interface{}, error) {
	client := db.GetClient()
	collection := client.Database("db").Collection("users")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []map[string]interface{}
	for cursor.Next(ctx) {
		var user map[string]interface{}
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
