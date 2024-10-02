package recovery

import (
	"context"
	"log"
	"time"

	"project/db"
	"project/db/operations"
	"project/global"

	"go.mongodb.org/mongo-driver/bson"
)
var (
	energyRecoveryInterval time.Duration
	energyRecoveryAmount   int
)

func init() {
	config, err := global.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	energyRecoveryInterval = time.Duration(config.EnergyRecoveryInterval) * time.Minute
	energyRecoveryAmount = config.EnergyRecoveryAmount
}

func StartEnergyRecovery() {
	ticker := time.NewTicker(energyRecoveryInterval)
	go func() {
		for range ticker.C {
			log.Println("Recovering energy")
			go recoverEnergy()
		}
	}()
}

func recoverEnergy() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := db.GetClient()
	collection := client.Database("db").Collection("users")

	filter := bson.M{"$expr": bson.M{"$lt": bson.A{"$energy", "$energy_max"}}}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Printf("Error finding users: %v", err)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user struct {
			UserID    int `bson:"user_id"`
			Energy    int32              `bson:"energy"`
			EnergyMax int32              `bson:"energy_max"`
		}

		if err := cursor.Decode(&user); err != nil {
			log.Printf("Error decoding user: %v", err)
			continue
		}

		energyToAdd := calculateEnergyToAdd(user.Energy, user.EnergyMax)

		if err := updateUserEnergy(ctx, user.UserID, energyToAdd); err != nil {
			log.Printf("Error updating user energy: %v", err)
		}
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
	}
}

func calculateEnergyToAdd(currentEnergy, maxEnergy int32) int32 {
	energyToAdd := int32(energyRecoveryAmount)
	if currentEnergy+energyToAdd > maxEnergy {
		energyToAdd = maxEnergy - currentEnergy
	}
	return energyToAdd
}

func updateUserEnergy(ctx context.Context, userID int, energyToAdd int32) error {
	return operations.IncrementUserField(ctx, userID, "energy", energyToAdd)
}
