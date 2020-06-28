package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AddUserLastUpdateTime saves last date when specified provider made sync
func (db Db) AddUserLastUpdateTime(login string, provider string, date time.Time) error {
	ctx := db.context()
	filter := bson.M{"login": login}
	update := bson.M{"$set": bson.M{"lastSync": bson.M{"date": date, "provider": provider}}}
	opts := options.Update()

	_, err := db.users.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}

// DeleteUserLastUpdateTime removes last date when specified provider made sync
func (db Db) DeleteUserLastUpdateTime(login string, provider string) error {
	ctx := db.context()
	filter := bson.M{"login": login}
	update := bson.M{"$unset": bson.M{"lastSync": ""}}
	opts := options.Update()

	_, err := db.users.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}

// GetUserLastUpdateTime receives last date when specified provider made sync
func (db Db) GetUserLastUpdateTime(login string, provider string) (time.Time, error) {
	filter := bson.M{"$and": []interface{}{bson.M{"login": login}, bson.M{"provider": provider}}}
	opts := options.FindOne()
	ctx := db.context()

	res := db.users.FindOne(ctx, filter, opts)
	if res.Err() == mongo.ErrNoDocuments {
		return time.Time{}, nil
	}
	var data struct {
		LastSync struct {
			Provider string    `bson:"provider"`
			Date     time.Time `bson:"date"`
		} `bson:"lastSync"`
	}
	res.Decode(&data)
	return data.LastSync.Date, nil
}
