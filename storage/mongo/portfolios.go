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

// GetPortfolios gets all portfolio fpvie user Id
func (db Db) GetPortfolios(userID string) ([]models.Portfolio, error) {
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("Could not decode user Id (%s). Internal error : %s", userID, err)
	}

	filter := bson.M{"uid": uid}
	findOptions := options.Find()
	ctx := db.context()

	cur, err := db.portfolios.Find(ctx, filter, findOptions)

	var transferObj []struct {
		ID   primitive.ObjectID `bson:"_id"`
		UID  primitive.ObjectID `bson:"uid"`
		Name string             `bson:"name"`
		Desc string             `bson:"desc"`
	}

	cur.All(ctx, &transferObj)

	result := make([]models.Portfolio, len(transferObj))
	for i, obj := range transferObj {
		result[i] = models.Portfolio{
			PortfolioID: obj.ID.Hex(),
			UserID:      obj.UID.Hex(),
			Name:        obj.Name,
			Description: obj.Desc,
		}
	}
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

// DeletePortfolio removes portfolio by Id
// Also removes all operations associated with this portfolio
func (db Db) DeletePortfolio(userID string, portfolioID string) (bool, error) {
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
	opts := options.Delete()

	res, err := db.portfolios.DeleteOne(ctx, filter, opts)
	if err != nil {
		return false, err
	}
	if res.DeletedCount >= 1 {
		return true, nil
	}

	filter = bson.M{"pid": pid}

	_, err = db.operations.DeleteMany(ctx, filter, opts)
	if err != nil {
		return false, err
	}

	return false, nil
}

// DeletePortfolios removes all portfolios for provided user
func (db Db) DeletePortfolios(userID string) (int64, error) {
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return 0, fmt.Errorf("Could not decode portfolio Id (%s). Internal error : %s", userID, err)
	}

	ctx := db.context()
	filter := bson.M{"uid": uid}
	opts := options.Delete()

	ps, err := db.GetPortfolios(userID)
	if err != nil {
		return 0, err
	}

	res, err := db.portfolios.DeleteMany(ctx, filter, opts)
	if err != nil {
		return 0, err
	}

	pids := make([]primitive.ObjectID, len(ps))
	for i, p := range ps {
		id, _ := primitive.ObjectIDFromHex(p.PortfolioID)
		pids[i] = id
	}

	filter = bson.M{"pid": bson.M{"$in": pids}}
	opts = options.Delete()

	_, err = db.operations.DeleteMany(ctx, filter, opts)
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
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
