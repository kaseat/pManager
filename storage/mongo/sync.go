package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SaveLastUpdateTime saves last date when specified provider made sync
func (db Db) SaveLastUpdateTime(provider string, date time.Time) error {
	ctx := db.context()
	filter := bson.M{"provider": provider}
	update := bson.M{"$set": bson.M{"date": date}}
	opts := options.Update().SetUpsert(true)

	_, err := db.syncs.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}

// ClearLastUpdateTime removes last date when specified provider made sync
func (db Db) ClearLastUpdateTime(provider string) error {
	ctx := db.context()
	filter := bson.M{"provider": provider}
	opts := options.Delete()

	_, err := db.syncs.DeleteOne(ctx, filter, opts)
	if err != nil {
		return err
	}
	return nil
}

// GetLastUpdateTime receives last date when specified provider made sync
func (db Db) GetLastUpdateTime(provider string) (time.Time, error) {
	filter := bson.M{"provider": provider}
	opts := options.FindOne()
	ctx := db.context()

	res := db.syncs.FindOne(ctx, filter, opts)
	if res.Err() == mongo.ErrNoDocuments {
		return time.Time{}, nil
	}
	var data struct {
		Provider string    `bson:"provider"`
		Date     time.Time `bson:"date"`
	}
	res.Decode(&data)
	return data.Date, nil
}
