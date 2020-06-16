package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SavePassword saves password hash to storage
func (db Db) SavePassword(user string, hash string) error {
	ctx := db.context()
	filter := bson.M{"user": user}
	update := bson.M{"$set": bson.M{"hash": hash}}
	opts := options.Update().SetUpsert(true)

	_, err := db.passwords.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}

// GetPassword gets password hash from storage
func (db Db) GetPassword(user string) (string, error) {

	filter := bson.M{"user": user}
	opts := options.FindOne()
	ctx := db.context()

	res := db.passwords.FindOne(ctx, filter, opts)

	if res.Err() == mongo.ErrNoDocuments {
		return "", nil
	}

	var data struct {
		User string `bson:"user"`
		Hash string `bson:"hash"`
	}
	res.Decode(&data)
	return data.Hash, nil
}

// DeletePassword removes password hash from storage
func (db Db) DeletePassword(user string) (bool, error) {
	ctx := db.context()
	filter := bson.M{"user": user}
	opts := options.Delete()

	res, err := db.syncs.DeleteOne(ctx, filter, opts)
	if err != nil {
		return false, err
	}
	return res.DeletedCount == 1, nil
}
