package db

import "go.mongodb.org/mongo-driver/bson/primitive"

type Notification struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Type   string             `bson:"type,omitempty" json:"type"`
	UserID primitive.ObjectID `bson:"user_id,omitempty" json:"user_id"`
	PostId primitive.ObjectID `bson:"post_id,omitempty" json:"post_id"`
}
