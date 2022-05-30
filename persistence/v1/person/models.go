package person

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type SearchPerson struct {
	Tags []string
}

type NewPerson struct {
	Name            string
	Type            string
	Role            *string
	ContactDuration *string
	Tags            []string
}

type Person struct {
	Id              string
	Name            string
	Type            string
	Role            *string
	ContactDuration *string
	Tags            []string
	UpdatedAt       time.Time
	CreatedAt       time.Time
}

type newPerson struct {
	Name            string    `bson:"name"`
	Type            string    `bson:"type"`
	Role            *string   `bson:"role"`
	ContactDuration *string   `bson:"contractDuration"`
	Tags            []string  `bson:"tags" example:"c#,c++"`
	UpdatedAt       time.Time `bson:"updatedAt"`
	CreatedAt       time.Time `bson:"createdAt"`
}

type updatePerson struct {
	Id              primitive.ObjectID `bson:"_id"`
	Name            string             `bson:"name"`
	Type            string             `bson:"type"`
	Role            *string            `bson:"role"`
	ContactDuration *string            `bson:"contractDuration"`
	Tags            []string           `bson:"tags" example:"c#,c++"`
	UpdatedAt       time.Time          `bson:"updatedAt"`
	CreatedAt       time.Time          `bson:"createdAt"`
}

type person struct {
	Id              primitive.ObjectID `bson:"_id"`
	Name            string             `bson:"name"`
	Type            string             `bson:"type"`
	Role            *string            `bson:"role"`
	ContactDuration *string            `bson:"contractDuration"`
	Tags            []string           `bson:"tags" example:"c#,c++"`
	UpdatedAt       time.Time          `bson:"updatedAt"`
	CreatedAt       time.Time          `bson:"createdAt"`
}
