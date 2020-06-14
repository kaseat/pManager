package mongo

import (
	"fmt"

	"github.com/kaseat/pManager/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AddPortfolio adds new potrfolio
func (db Db) AddPortfolio(userID string, p models.Portfolio) (string, error) {
	uid, err := db.findUser(userID)

	if err != nil {
		return "", err
	}
	if uid.IsZero() {
		return "", fmt.Errorf("No user found with %s Id", userID)
	}

	doc := bson.M{
		"uid":  uid,
		"name": p.Name,
		"desc": p.Description,
	}

	ctx := db.context()
	opts := options.InsertOne()
	res, err := db.portfolios.InsertOne(ctx, doc, opts)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

// GetPortfolio gets operation by id
func (db Db) GetPortfolio(userID string, portfolioID string) (models.Portfolio, error) {
	var result models.Portfolio

	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return result, fmt.Errorf("Could not decode user Id (%s). Internal error : %s", userID, err)
	}
	pid, err := primitive.ObjectIDFromHex(portfolioID)
	if err != nil {
		return result, fmt.Errorf("Could not decode portfolio Id (%s). Internal error : %s", portfolioID, err)
	}

	filter := bson.M{"$and": []interface{}{bson.M{"_id": pid}, bson.M{"uid": uid}}}
	findOptions := options.FindOne()
	ctx := db.context()

	r := db.portfolios.FindOne(ctx, filter, findOptions)

	if r.Err() != nil {
		return result, r.Err()
	}
	var transferObj struct {
		ID   primitive.ObjectID `bson:"_id"`
		UID  primitive.ObjectID `bson:"uid"`
		Name string             `bson:"name"`
		Desc string             `bson:"desc"`
	}

	r.Decode(&transferObj)

	result.Description = transferObj.Desc
	result.Name = transferObj.Name
	result.PortfolioID = portfolioID
	result.UserID = userID
	return result, nil
}

// UpdatePortfolio updates portfolio with provided uid and pid
func (db Db) UpdatePortfolio(userID string, portfolioID string, p models.Portfolio) (bool, error) {
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, fmt.Errorf("Could not decode user Id (%s). Internal error : %s", userID, err)
	}
	pid, err := primitive.ObjectIDFromHex(portfolioID)
	if err != nil {
		return false, fmt.Errorf("Could not decode portfolio Id (%s). Internal error : %s", portfolioID, err)
	}
	ctx := db.context()
	filter := bson.M{"$and": []interface{}{bson.M{"_id": pid}, bson.M{"uid": uid}}}
	update := bson.M{
		"$set": bson.M{
			"name":        p.Name,
			"description": p.Description,
		},
	}

	res, err := db.portfolios.UpdateOne(ctx, filter, update)
	if err != nil {
		return false, err
	}
	if res.ModifiedCount > 0 {
		return true, nil
	}
	return false, nil
}

// Checks if user with specified _id exists. Then needs to be checked on .IsZero()
func (db Db) findUser(uid string) (primitive.ObjectID, error) {
	id, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		return id, err
	}

	ctx := db.context()
	filter := bson.M{"_id": id}
	opts := options.FindOne()

	r := db.users.FindOne(ctx, filter, opts)
	var result struct {
		ID primitive.ObjectID `bson:"_id"`
	}

	if r.Err() != nil {
		return result.ID, nil
	}

	r.Decode(&result)
	return result.ID, nil
}
