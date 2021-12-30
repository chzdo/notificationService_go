package models

import (
	"context"
	"encoding/json"
	"notification_service/internals/logger"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func RegisterModel(logs *logger.Logger) (map[string]*DBModels, error) {

	db_uri := os.Getenv("DB_URI")

	// if !ok {
	// 	logs.ErrorLogs.Println("Could not find DB_URI")
	// }

	db_name := os.Getenv("DB_NAME")

	clientOption := options.Client().ApplyURI(db_uri)

	client, err := mongo.Connect(context.TODO(), clientOption)

	if err != nil {
		return nil, err
	}

	//defer client.Disconnect(ctx)

	db := client.Database(db_name)

	modellist := []string{
		"triggers",
		"usermobilesettings",
		"organizationsettings",
		"organizationnotifications",
		"userseennotification",
	}

	var model = make(map[string]*DBModels)

	for _, value := range modellist {

		model[value] = &DBModels{
			name: value,
			db:   db,
		}
	}

	return model, nil

}

type DBModels struct {
	db   *mongo.Database
	name string
}

func (m DBModels) Insert(data interface{}) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()

	num, err := m.db.Collection(m.name).CountDocuments(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	dataToBson, err := bson.Marshal(data)

	if err != nil {
		return nil, err
	}

	bsonToInsert := make(map[string]interface{})

	err = bson.Unmarshal(dataToBson, &bsonToInsert)

	if err != nil {
		return nil, err
	}
	bsonToInsert["id"] = int(num) + 1
	bsonToInsert["isActive"] = true
	bsonToInsert["isDeleted"] = false
	bsonToInsert["createdOn"] = time.Now()
	bsonToInsert["updatedOn"] = time.Now()

	_, err = m.db.Collection(m.name).InsertOne(ctx, bsonToInsert)

	if err != nil {
		return nil, err
	}

	return bsonToInsert, nil
}

func (m DBModels) Get(query map[string]interface{}) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()

	var modelToReturn = map[string]interface{}{}

	err := m.db.Collection(m.name).FindOne(ctx, query).Decode(&modelToReturn)

	if err != nil {
		return nil, err
	}

	byteToRetutn, err := json.Marshal(modelToReturn)

	if err != nil {
		return nil, err
	}
	return byteToRetutn, nil

}

func (m DBModels) GetAll(query interface{}) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()

	opts := options.Find()

	opts.SetLimit(100)

	c, err := m.db.Collection(m.name).Find(ctx, query, opts)

	if err != nil {
		return nil, err
	}
	var modelToReturn = make([]map[string]interface{}, 0)

	if err = c.All(ctx, &modelToReturn); err != nil {
		return nil, err
	}

	byteToRetutn, err := json.Marshal(modelToReturn)

	if err != nil {
		return nil, err
	}
	return byteToRetutn, nil

}

func (m DBModels) UpdateOrInsert(query interface{}, updateBody map[string]interface{}) (interface{}, error) {

	_, err := m.Get(query.(map[string]interface{}))

	if err != nil {

		//doest exist
		return m.Insert(updateBody)

	} else {
		//exist

		return m.Update(query, updateBody)

	}

}

func (m DBModels) Update(query interface{}, updateBody map[string]interface{}) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()

	delete(updateBody, "id")
	delete(updateBody, "_id")
	delete(updateBody, "isActive")
	delete(updateBody, "isDeleted")
	delete(updateBody, "createdOn")

	update := map[string]map[string]interface{}{
		"$set": updateBody,
		"$currentDate": {
			"updatedOn": true,
		},
	}
	opts := options.Update()

	return m.db.Collection(m.name).UpdateOne(ctx, query, update, opts)

}

func (m DBModels) Delete(query interface{}) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()

	update := map[string]map[string]interface{}{
		"$set": {
			"isActive":  false,
			"isDeleted": true,
		},
		"$currentDate": {
			"updatedOn": true,
		},
	}
	opts := options.Update()

	return m.db.Collection(m.name).UpdateOne(ctx, query, update, opts)

}
