package mongo

import (
	"time"

	"github.com/kaseat/pManager/models/provider"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AddUserLastUpdateTime saves last date when specified provider made sync
func (db Db) AddUserLastUpdateTime(login string, provider provider.Type, date time.Time) error {
	ctx := db.context()
	filter := bson.M{"login": login}
	update := bson.M{"$push": bson.M{"lastSync": bson.M{"date": date.Format(time.RFC3339Nano), "provider": provider}}}
	opts := options.Update()

	_, err := db.users.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}

// DeleteUserLastUpdateTime removes last date when specified provider made sync
func (db Db) DeleteUserLastUpdateTime(login string, provider provider.Type) error {
	ctx := db.context()
	filter := bson.M{"login": login}
	update := bson.M{"$pull": bson.M{"lastSync": bson.M{"provider": provider}}}
	opts := options.Update()

	_, err := db.users.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}

// GetUserLastUpdateTime receives last date when specified provider made sync
func (db Db) GetUserLastUpdateTime(login string, provider provider.Type) (time.Time, error) {
	filter := bson.M{"login": login}
	opts := options.FindOne()
	ctx := db.context()

	res := db.users.FindOne(ctx, filter, opts)
	if res.Err() == mongo.ErrNoDocuments {
		return time.Time{}, nil
	}
	var data struct {
		LastSync []struct {
			Provider string `bson:"provider"`
			Date     string `bson:"date"`
		} `bson:"lastSync"`
	}

	res.Decode(&data)
	result := time.Time{}
	for _, it := range data.LastSync {
		if it.Provider == string(provider) {
			t, _ := time.Parse(time.RFC3339Nano, it.Date)
			result = t
			break
		}
	}

	return result, nil
}
