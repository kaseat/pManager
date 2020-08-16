package mongo

import (
	"time"

	"github.com/kaseat/pManager/models"
	"github.com/kaseat/pManager/models/currency"
	"github.com/kaseat/pManager/models/instrument"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AddInstruments saves instruments info into a storage
func (db Db) AddInstruments(instr []models.Instrument) error {
	if len(instr) == 0 {
		return nil
	}
	docs := make([]interface{}, len(instr))
	for i, p := range instr {
		doc := bson.M{
			"ticker": p.Ticker,
			"name":   p.Name,
			"curr":   p.Currency,
			"type":   p.Type,
		}
		if !p.PriceUptdTime.IsZero() {
			doc["lut"] = p.PriceUptdTime
		}
		if p.ISIN != "" {
			doc["isin"] = p.ISIN
		}
		if p.FIGI != "" {
			doc["figi"] = p.FIGI
		}
		docs[i] = doc
	}

	ctx := db.context()
	opts := options.InsertMany()
	_, err := db.instruments.InsertMany(ctx, docs, opts)
	if err != nil {
		return err
	}

	return nil
}

// SetInstrumentPriceUptdTime sets time instrument prise was updated
func (db Db) SetInstrumentPriceUptdTime(isin string, updTime time.Time) (bool, error) {
	ctx := db.context()
	filter := bson.M{"isin": isin}
	opts := options.Update()

	update := bson.M{
		"$set": bson.M{
			"lut": updTime,
		},
	}

	u, err := db.instruments.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return false, err
	}

	if u.ModifiedCount == 1 {
		return true, nil
	}

	return false, nil
}

// ClearInstrumentPriceUptdTime clears time instrument prise was updated
func (db Db) ClearInstrumentPriceUptdTime(isin string) (bool, error) {
	ctx := db.context()
	filter := bson.M{"isin": isin}
	opts := options.Update()

	update := bson.M{
		"$unset": bson.M{
			"lut": "",
		},
	}

	u, err := db.instruments.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return false, err
	}

	if u.ModifiedCount == 1 {
		return true, nil
	}

	return false, nil
}

// ClearAllInstrumentPriceUptdTime clears time instrument prise was updated (for all)
func (db Db) ClearAllInstrumentPriceUptdTime() (bool, error) {
	ctx := db.context()
	filter := bson.M{}
	opts := options.Update()

	update := bson.M{
		"$unset": bson.M{
			"lut": "",
		},
	}

	u, err := db.instruments.UpdateMany(ctx, filter, update, opts)
	if err != nil {
		return false, err
	}

	if u.ModifiedCount != 0 {
		return true, nil
	}

	return false, nil
}

// GetInstruments finds instruments depending on input prameters
func (db Db) GetInstruments(key string, value string) ([]models.Instrument, error) {
	filter := bson.M{key: value}
	findOptions := options.Find()
	return db.getInstruments(filter, findOptions)
}

// GetAllInstruments finds all instruments
func (db Db) GetAllInstruments() ([]models.Instrument, error) {
	filter := bson.M{}
	findOptions := options.Find()
	return db.getInstruments(filter, findOptions)
}

// DeleteInstruments removes instruments depending on input prameters
func (db Db) DeleteInstruments(key string, value string) (int64, error) {
	filter := bson.M{key: value}
	delOptions := options.Delete()
	return db.delInstruments(filter, delOptions)
}

// DeleteAllInstruments removes all instruments from storage
func (db Db) DeleteAllInstruments() (int64, error) {
	filter := bson.M{}
	delOptions := options.Delete()
	return db.delInstruments(filter, delOptions)
}

func (db Db) delInstruments(filter primitive.M, delOptions *options.DeleteOptions) (int64, error) {
	ctx := db.context()
	del, err := db.instruments.DeleteMany(ctx, filter, delOptions)
	if err != nil {
		return 0, err
	}
	return del.DeletedCount, nil
}

func (db Db) getInstruments(filter primitive.M, findOptions *options.FindOptions) ([]models.Instrument, error) {
	ctx := db.context()
	ins, err := db.instruments.Find(ctx, filter, findOptions)
	defer ins.Close(ctx)

	if err != nil {
		return nil, err
	}

	var raw []struct {
		ISIN          string    `bson:"isin"`
		FIGI          string    `bson:"figi"`
		Ticker        string    `bson:"ticker"`
		Name          string    `bson:"name"`
		Currency      string    `bson:"curr"`
		Type          string    `bson:"type"`
		PriceUptdTime time.Time `bson:"lut"`
	}

	err = ins.All(ctx, &raw)
	if err != nil {
		return nil, err
	}

	if raw == nil {
		return []models.Instrument{}, nil
	}
	results := make([]models.Instrument, len(raw))

	for i, item := range raw {
		data := models.Instrument{
			ISIN:          item.ISIN,
			FIGI:          item.FIGI,
			Ticker:        item.Ticker,
			Name:          item.Name,
			Currency:      currency.Type(item.Currency),
			Type:          instrument.Type(item.Type),
			PriceUptdTime: item.PriceUptdTime,
		}
		results[i] = data
	}

	return results, err
}
