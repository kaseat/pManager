package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AddTcsToken adds token to access tcs API
func (db Db) AddTcsToken(token string) error {
	ctx := db.context()
	filter := bson.M{}
	opts := options.Update()
	opts.SetUpsert(true)

	update := bson.M{
		"$set": bson.M{
			"token": token,
		},
	}

	_, err := db.settings.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}

// DeleteTcsToken deletes token to access tcs API
func (db Db) DeleteTcsToken() error {
	ctx := db.context()
	filter := bson.M{}
	delOptions := options.Delete()
	_, err := db.settings.DeleteMany(ctx, filter, delOptions)
	if err != nil {
		return err
	}
	return nil
}

// GetTcsToken finds token to access tcs API
func (db Db) GetTcsToken() (string, error) {
	ctx := db.context()
	filter := bson.M{}
	findOptions := options.FindOne()
	ins := db.settings.FindOne(ctx, filter, findOptions)
	var raw struct {
		Token string `bson:"token"`
	}
	err := ins.Decode(&raw)
	if err != nil {
		return "", err
	}
	return raw.Token, nil
}
