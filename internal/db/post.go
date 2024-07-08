package db

import "go.mongodb.org/mongo-driver/bson/primitive"

type Post struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     primitive.ObjectID `bson:"user_id,omitempty" json:"user_id"`
	Content    string             `bson:"content,omitempty" json:"content"`
	LikesCount int64              `bson:"likes_count,omitempty" json:"likes_count"`
}
