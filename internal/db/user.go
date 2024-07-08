package db

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name          string               `bson:"name" json:"name"`
	Avatar        string               `bson:"avatar" json:"avatar"`
	Posts         []primitive.ObjectID `bson:"posts" json:"posts"`
	LikedPosts    []primitive.ObjectID `bson:"liked_posts" json:"liked_posts"`
	Notifications []primitive.ObjectID `bson:"notifications" json:"notifications"`
}
