package portfolio

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AddOwner adds new potrfolio owner
func AddOwner(login, firstName, lastName string) (Owner, error) {
	owner := Owner{
		Login:     login,
		FirstName: firstName,
		LastName:  lastName,
	}

	res, err := db.owners.InsertOne(db.context(), owner)
	if err != nil {
		return owner, err
	}

	if id, ok := res.InsertedID.(primitive.ObjectID); ok {
		owner.OwnerID = id.Hex()
		return owner, nil
	}

	return owner, errors.New("Filed convert 'primitive.ObjectID' to 'string'")
}

// DeleteOwner removes owner by Id
func DeleteOwner(ownerID string) (bool, error) {
	ctx := db.context()

	objID, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		err = errors.New("Invalid portfolio Id")
		return false, err
	}

	filter := bson.M{"_id": objID}
	res, err := db.owners.DeleteOne(ctx, filter)
	if err != nil {
		return false, err
	}

	if res.DeletedCount != 0 {
		return true, err
	}

	return false, err
}

// GetOwnerByLogin gets owner by login
func GetOwnerByLogin(login string) (bool, Owner, error) {
	var result Owner

	filter := bson.M{"login": login}
	findOptions := options.Find()
	ctx := db.context()

	cur, err := db.owners.Find(ctx, filter, findOptions)
	defer cur.Close(ctx)
	if err != nil {
		return false, result, err
	}

	hasResult := cur.TryNext(ctx)
	if !hasResult {
		return false, result, nil
	}

	err = cur.Decode(&result)
	if err != nil {
		return false, result, err
	}
	return true, result, nil
}
