package handler

import (
	"context"
	"gbADMIN/config"
	"gbADMIN/schema"
)

type DataBaseHandler struct {
	Mongo config.MongoDB
}

func (db *DataBaseHandler) SetNewLog(log schema.Log) error {
	_, err := db.Mongo.Logger.InsertOne(
		context.TODO(),
		log,
	)
	if err != nil {
		return err
	}
	return nil
}

// // if operation find no document it will return ErrNoDocuments error.
// func (db *DataBaseHandler) IsUserRegistered(id string) error {
// 	_id, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return err
// 	}

// 	result := db.Mongo.Users.FindOne(
// 		context.TODO(),
// 		bson.M{"_id": _id},
// 	)
// 	if result.Err() != nil {
// 		return err
// 	}
// 	return nil
// }
