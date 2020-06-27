package mongo

import (
	"errors"
	"time"

	"github.com/kaseat/pManager/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/oauth2"
)

// AddUser saves user and password hash to storage
func (db Db) AddUser(login, email, hash string) (string, error) {
	ctx := db.context()
	filter := bson.M{"login": login}
	optsFind := options.FindOne()
	found := db.users.FindOne(ctx, filter, optsFind)
	if found.Err() == mongo.ErrNoDocuments {
		opts := options.InsertOne()
		doc := bson.M{"login": login, "hash": hash}
		if email != "" {
			doc["email"] = email
		}
		res, err := db.users.InsertOne(ctx, doc, opts)
		if err != nil {
			return "", err
		}
		return res.InsertedID.(primitive.ObjectID).Hex(), nil
	} else if found.Err() != nil {
		return "", found.Err()
	}
	return "", errors.New("User with this login already exists")
}

// UpdateUser updates user info
func (db Db) UpdateUser(login string, user models.User) (bool, error) {
	ctx := db.context()
	filter := bson.M{"login": login}
	optsFind := options.FindOne()
	found := db.users.FindOne(ctx, filter, optsFind)
	if found.Err() == mongo.ErrNoDocuments {
		return false, nil
	} else if found.Err() != nil {
		return false, found.Err()
	}

	opts := options.Update()
	update := bson.M{"$set": bson.M{"login": user.Login, "email": user.Email}}
	res, err := db.users.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return false, err
	}
	return res.ModifiedCount == 1, nil
}

// AddUserToken adds oauth2 token to user
func (db Db) AddUserToken(state string, token *oauth2.Token) error {

	doc := bson.M{
		"access_token":  token.AccessToken,
		"token_type":    token.TokenType,
		"refresh_token": token.RefreshToken,
		"expiry":        token.Expiry.Format(time.RFC3339Nano),
	}

	filter := bson.M{"state": state}
	update := bson.M{"$set": bson.M{"token": doc}}
	ctx := db.context()

	res := db.users.FindOneAndUpdate(ctx, filter, update)
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}

// GetUserToken gets user's oauth2 token
func (db Db) GetUserToken(login string) (oauth2.Token, error) {

	filter := bson.M{"login": login}
	opts := options.FindOne()
	ctx := db.context()

	res := db.users.FindOne(ctx, filter, opts)
	if res.Err() != nil {
		return oauth2.Token{}, res.Err()
	}

	var data struct {
		Token token `bson:"token"`
	}

	err := res.Decode(&data)
	if err != nil {
		return oauth2.Token{}, err
	}

	t, _ := time.Parse(time.RFC3339Nano, data.Token.Expiry)
	tok := oauth2.Token{
		AccessToken:  data.Token.AccessToken,
		TokenType:    data.Token.TokenType,
		RefreshToken: data.Token.RefreshToken,
		Expiry:       t,
	}

	return tok, nil
}

// AddUserState adds state to user
func (db Db) AddUserState(login string, state string) error {

	filter := bson.M{"login": login}
	update := bson.M{"$set": bson.M{"state": state}}
	ctx := db.context()

	res := db.users.FindOneAndUpdate(ctx, filter, update)
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}

// GetUserState gets user's state
func (db Db) GetUserState(login string) (string, error) {

	filter := bson.M{"login": login}
	opts := options.FindOne()
	ctx := db.context()

	res := db.users.FindOne(ctx, filter, opts)
	if res.Err() != nil {
		return "", res.Err()
	}

	var data struct {
		State string `bson:"state"`
	}

	err := res.Decode(&data)
	if err != nil {
		return "", err
	}

	return data.State, nil
}

// GetUserByLogin gets user by login
func (db Db) GetUserByLogin(login string) (models.User, error) {
	result := models.User{}
	filter := bson.M{"login": login}
	opts := options.FindOne()
	ctx := db.context()

	res := db.users.FindOne(ctx, filter, opts)
	if res.Err() == mongo.ErrNoDocuments {
		return result, res.Err()
	}

	var data struct {
		ID    primitive.ObjectID `bson:"_id"`
		Login string             `bson:"login"`
		Email string             `bson:"email"`
	}
	res.Decode(&data)

	result.UserID = data.ID.Hex()
	result.Login = data.Login
	result.Email = data.Email
	return result, nil
}

// GetUserPassword gets password hash from storage
func (db Db) GetUserPassword(login string) (string, error) {
	filter := bson.M{"login": login}
	opts := options.FindOne()
	ctx := db.context()

	res := db.users.FindOne(ctx, filter, opts)
	if res.Err() == mongo.ErrNoDocuments {
		return "", nil
	}

	var data struct {
		Hash string `bson:"hash"`
	}
	res.Decode(&data)
	return data.Hash, nil
}

// UpdateUserPassword updates user password
func (db Db) UpdateUserPassword(login, hash string) (bool, error) {
	ctx := db.context()
	filter := bson.M{"login": login}
	optsFind := options.FindOne()
	found := db.users.FindOne(ctx, filter, optsFind)
	if found.Err() == mongo.ErrNoDocuments {
		return false, nil
	} else if found.Err() != nil {
		return false, found.Err()
	}

	opts := options.Update()
	update := bson.M{"$set": bson.M{"hash": hash}}
	res, err := db.users.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return false, err
	}
	return res.ModifiedCount == 1, nil
}

// DeleteUser removes password hash from storage
func (db Db) DeleteUser(login string) (bool, error) {
	ctx := db.context()
	filter := bson.M{"login": login}
	opts := options.Delete()

	res, err := db.users.DeleteOne(ctx, filter, opts)
	if err != nil {
		return false, err
	}
	return res.DeletedCount == 1, nil
}
