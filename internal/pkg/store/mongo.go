package store

import (
	"context"
	"fmt"
	"macrotrack/internal/pkg/types"
	"time"

	"github.com/google/uuid"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoDB struct {
	client *mongo.Client
}

func (s *mongoDB) Init() error {

	var err error

	credential := options.Credential{
		Username: "mongo",
		Password: "mongo",
	}

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017").SetAuth(credential)

	s.client, err = mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		fmt.Println(err)
		return err
	}

	if err := s.client.Ping(context.TODO(), readpref.Primary()); err != nil {
		fmt.Println(err)
		return err
	}

	//macrosCollection := s.client.Database("macrosdb").Collection("macros")

	return nil
}

func (s *mongoDB) Open() error {
	return nil
}

// func (s *mongoDB) Create(m types.Macro) (uuid.NullUUID, error) {
func (s *mongoDB) Create(m types.Macro) (string, error) {

	/*
		type Macro struct {
			Carbs   int
			Fat     int
			Protein int
			Alcohol int
			Date    string
		}
	*/

	var retString string

	// use uuid for id
	u, err := uuid.NewUUID()
	if err != nil {
		return retString, err
	}

	macrosCollection := s.client.Database("macrosdb").Collection("macros")

	mb := &types.Macro_bson{ID: u, Carbs: m.Carbs, Fat: m.Fat, Protein: m.Protein, Date: time.Now()}

	result, err := macrosCollection.InsertOne(context.TODO(), mb)

	if err != nil {
		return retString, err
	}

	//
	//https://www.mongodb.com/community/forums/t/decode-uuid-binary-format-to-bson-interface/165655
	pb, ok := result.InsertedID.(primitive.Binary)

	if ok {
		retUUID, err := uuid.FromBytes(pb.Data)
		if err == nil {

			if retUUID == u {
				retString = retUUID.String()
			}

		}
	}

	return retString, err

}

func (s *mongoDB) Read(_uuid uuid.UUID) (*types.Macro, error) {

	macrosCollection := s.client.Database("macrosdb").Collection("macros")

	filter := bson.M{"_id": _uuid}

	var mb types.Macro_bson

	err := macrosCollection.FindOne(context.TODO(), filter).Decode(&mb)

	if err != nil {
		return nil, err
	}

	return &types.Macro{Carbs: mb.Carbs, Protein: mb.Protein, Fat: mb.Fat, Alcohol: mb.Alcohol, Date: mb.Date.GoString()}, nil
}

func (s *mongoDB) Update() error {
	return nil
}
func (s *mongoDB) Delete(_uuid uuid.UUID) error {

	macrosCollection := s.client.Database("macrosdb").Collection("macros")

	filter := bson.M{"_id": _uuid}

	result, err := macrosCollection.DeleteOne(context.TODO(), filter)

	if err != nil {
		return err
	}

	fmt.Println(result)

	return nil
}
