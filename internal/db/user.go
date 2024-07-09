package db

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name          string               `bson:"name" json:"name"`
	Avatar        string               `bson:"avatar" json:"avatar"`
	Posts         []primitive.ObjectID `bson:"posts,omitempty" json:"posts,omitempty"`
	LikedPosts    []primitive.ObjectID `bson:"liked_posts,omitempty" json:"liked_posts,omitempty"`
	Notifications []primitive.ObjectID `bson:"notifications,omitempty" json:"notifications,omitempty"`
}
