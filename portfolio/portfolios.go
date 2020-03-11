package portfolio

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db database
var cfg Config

// Init portgolio module
func Init(config Config) error {
	cfg = config
	db.context = func() context.Context { return context.Background() }
	clientOptions := options.Client().ApplyURI(cfg.MongoURL)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return err
	}
	err = client.Connect(db.context())
	if err != nil {
		return err
	}
	db.operations = client.Database(cfg.DbName).Collection("Operations")
	db.portfolios = client.Database(cfg.DbName).Collection("Portfolios")
	db.owners = client.Database(cfg.DbName).Collection("Owners")
	db.prices = client.Database(cfg.DbName).Collection("Prices")
	return nil
}

// AddPortfolio adds new potrfolio
func (o *Owner) AddPortfolio(name, description string) (Portfolio, error) {
	portfolio := Portfolio{
		OwnerID:     o.OwnerID,
		Name:        name,
		Description: description}

	res, err := db.portfolios.InsertOne(db.context(), portfolio)
	if err != nil {
		return portfolio, err
	}

	if id, ok := res.InsertedID.(primitive.ObjectID); ok {
		portfolio.PortfolioID = id.Hex()
		return portfolio, nil
	}

	return portfolio, errors.New("Filed convert 'primitive.ObjectID' to 'string'")
}

// GetPortfolio gets operation by id
func (o *Owner) GetPortfolio(portfolioID string) (bool, Portfolio, error) {
	var result Portfolio

	objID, err := primitive.ObjectIDFromHex(portfolioID)
	if err != nil {
		err = errors.New("Invalid portfolio Id")
		return false, result, err
	}

	filter := bson.M{"$and": []interface{}{bson.M{"_id": objID}, bson.M{"oid": o.OwnerID}}}
	findOptions := options.Find()
	ctx := db.context()

	cur, err := db.portfolios.Find(ctx, filter, findOptions)
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

// UpdatePortfolio updates current portfolio
func (p *Portfolio) UpdatePortfolio() (bool, error) {

	objID, err := primitive.ObjectIDFromHex(p.PortfolioID)
	if err != nil {
		err = errors.New("Invalid portfolio Id")
		return false, err
	}

	filter := bson.M{"$and": []interface{}{bson.M{"_id": objID}, bson.M{"oid": p.OwnerID}}}
	update := bson.M{
		"$set": bson.M{
			"name":        p.Name,
			"description": p.Description,
		},
	}
	res, err := db.portfolios.UpdateOne(db.context(), filter, update)
	if err != nil {
		return false, err
	}
	if res.ModifiedCount > 0 {
		return true, nil
	}
	return false, nil
}

// GetAllPortfolios finds all available portfolios at the moment
func (o *Owner) GetAllPortfolios() ([]Portfolio, error) {
	filter := bson.M{"oid": o.OwnerID}
	findOptions := options.Find()

	ctx := db.context()
	cur, err := db.portfolios.Find(ctx, filter, findOptions)
	defer cur.Close(ctx)

	if err != nil {
		return nil, err
	}

	var results []Portfolio

	err = cur.All(ctx, &results)
	if err != nil {
		return nil, err
	}

	return results, err
}

// DeleteAllPortfolios removes all portfolios
func (o *Owner) DeleteAllPortfolios() (bool, error) {
	ctx := db.context()
	filter := bson.M{"oid": o.OwnerID}

	res, err := db.portfolios.DeleteMany(ctx, filter)
	if err != nil {
		return false, err
	}

	if res.DeletedCount != 0 {
		return true, nil
	}

	return false, nil
}

// DeletePortfolio removes portfolio by Id
func (o *Owner) DeletePortfolio(portfolioID string) (bool, error) {
	ctx := db.context()

	objID, err := primitive.ObjectIDFromHex(portfolioID)
	if err != nil {
		err = errors.New("Invalid portfolio Id")
		return false, err
	}

	filter := bson.M{"$and": []interface{}{bson.M{"_id": objID}, bson.M{"oid": o.OwnerID}}}
	res, err := db.portfolios.DeleteOne(ctx, filter)
	if err != nil {
		return false, err
	}

	if res.DeletedCount != 0 {
		return true, err
	}

	return false, err
}

func (p *Portfolio) String() string {
	s := p.PortfolioID + " " + p.Name + " " + p.Description
	return s
}
