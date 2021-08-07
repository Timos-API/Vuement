package persistence

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ComponentPersistor struct {
	c *mongo.Collection
}

type Component struct {
	ComponentID string               `json:"id" bson:"_id"`
	Name        string               `json:"name" bson:"name" keep:"update,create,omitempty" validate:"required,gt=2"`
	Image       string               `json:"image" bson:"image" keep:"update,create,omitempty" validate:"omitempty,url"`
	Children    []primitive.ObjectID `json:"children" bson:"children" keep:"update,create" validate:"required"`
	IsChild     bool                 `json:"isChild" bson:"isChild" keep:"update,create,omitempty"`
	Props       []ComponentProp      `json:"props" bson:"props" keep:"update,create" validate:"dive,required"`
}

type ComponentProp struct {
	Name        string `json:"name" bson:"name" keep:"update,create,omitempty" validate:"required,gt=2"`
	Value       string `json:"value" bson:"value" keep:"update,create,omitempty" validate:"required,gt=2"`
	Description string `json:"description" bson:"description" keep:"update,create,omitempty" validate:"required,gt=2"`
	Type        string `json:"type" bson:"type" keep:"update,create,omitempty" validate:"required,gt=2"`
}

var (
	ErrInsertError     = errors.New("something ubiquitous happened")
	ErrInvalidObjectID = errors.New("invalid ObjectID")
	ErrNothingDeleted  = errors.New("nothing has been deleted")
)

func NewComponentPersistor(c *mongo.Collection) *ComponentPersistor {
	return &ComponentPersistor{c}
}

func (p *ComponentPersistor) Create(ctx context.Context, component interface{}) (*Component, error) {
	res, err := p.c.InsertOne(ctx, component)

	if err != nil {
		return nil, err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		return p.GetById(ctx, oid.Hex())
	}

	return nil, ErrInsertError
}

func (p *ComponentPersistor) Update(ctx context.Context, id string, component interface{}) (*Component, error) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	res := p.c.FindOneAndUpdate(ctx, bson.M{"_id": oid}, bson.M{"$set": component}, options.FindOneAndUpdate().SetReturnDocument(options.After))

	if res.Err() != nil {
		return nil, res.Err()
	}

	var comp Component
	err = res.Decode(&comp)

	if err != nil {
		return nil, err
	}

	return &comp, nil
}

func (p *ComponentPersistor) Delete(ctx context.Context, id string) (bool, error) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return false, err
	}

	res, err := p.c.DeleteOne(ctx, bson.M{"_id": oid})

	if err != nil {
		return false, err
	}

	if res.DeletedCount == 0 {
		return false, ErrNothingDeleted
	}

	return true, nil
}

func (p *ComponentPersistor) GetById(ctx context.Context, id string) (*Component, error) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	res := p.c.FindOne(ctx, bson.M{"_id": oid})

	if res.Err() != nil {
		return nil, res.Err()
	}

	var comp Component
	err = res.Decode(&comp)

	if err != nil {
		return nil, err
	}

	return &comp, nil
}

func (p *ComponentPersistor) GetAll(ctx context.Context) (*[]Component, error) {
	cursor, err := p.c.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	components := []Component{}

	for cursor.Next(ctx) {
		var comp Component
		err := cursor.Decode(&comp)

		if err != nil {
			return nil, err
		}

		components = append(components, comp)
	}

	if cursor.Err() != nil {
		return nil, cursor.Err()
	}

	return &components, nil
}
