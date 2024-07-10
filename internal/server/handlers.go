package server

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"glam/internal/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"log"
	"net/http"
	"time"
)

func getProfile(ctx *gin.Context, client *mongo.Client) {
	userId := ctx.Param("id")
	collection := client.Database("glam").Collection("users")

	user := db.User{}

	objId, _ := primitive.ObjectIDFromHex(userId)
	filter := bson.M{"_id": objId}

	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func getProfiles(ctx *gin.Context, client *mongo.Client) {
	collection := client.Database("glam").Collection("users")

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var users []db.User
	if err = cursor.All(ctx, &users); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func createOrModifyProfile(ctx *gin.Context, client *mongo.Client) {
	var user db.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := client.Database("glam").Collection("users")
	if user.ID.IsZero() {
		if user.Posts == nil {
			user.Posts = []primitive.ObjectID{}
		}
		if user.LikedPosts == nil {
			user.LikedPosts = []primitive.ObjectID{}
		}
		if user.Notifications == nil {
			user.Notifications = []primitive.ObjectID{}
		}
		result, err := collection.InsertOne(context.Background(), user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		user.ID = result.InsertedID.(primitive.ObjectID)
		ctx.JSON(http.StatusCreated, user)
	} else {
		filter := bson.M{"_id": user.ID}
		update := bson.M{
			"$set": user,
		}

		// Выполняем обновление
		result, err := collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"modified_count": result.ModifiedCount})
	}
}

func GetAllPostsByAuthorID(ctx *gin.Context, client *mongo.Client) {
	userId := ctx.Param("id")

	// Получение коллекции "posts"
	collection := client.Database("glam").Collection("posts")

	// Парсинг ObjectID из строки authorID
	objID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Поиск постов с указанным authorID
	filter := bson.M{"user_id": objID}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	// Слайс для хранения найденных постов
	posts := []db.Post{}
	for cursor.Next(context.Background()) {
		var post db.Post
		if err := cursor.Decode(&post); err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		posts = append(posts, post)
	}

	if err := cursor.Err(); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"posts": posts})
}

func GetAllPosts(ctx *gin.Context, client *mongo.Client) {

	// Получение коллекции "posts"
	collection := client.Database("glam").Collection("posts")

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	// Слайс для хранения найденных постов
	posts := []db.Post{}
	for cursor.Next(context.Background()) {
		var post db.Post
		if err := cursor.Decode(&post); err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		posts = append(posts, post)
	}

	if err := cursor.Err(); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"posts": posts})
}

func CreatePost(ctx *gin.Context, client *mongo.Client) {
	var post db.Post

	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	postsCollection := client.Database("glam").Collection("posts")
	usersCollection := client.Database("glam").Collection("users")

	session, err := client.StartSession()
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
	}
	defer session.EndSession(context.TODO())

	transactionOptions := options.Transaction().SetWriteConcern(writeconcern.Majority())

	err = mongo.WithSession(context.TODO(), session, func(sessionContext mongo.SessionContext) error {
		// Начало транзакции
		err := session.StartTransaction(transactionOptions)
		if err != nil {
			return err
		}

		result, err := postsCollection.InsertOne(context.Background(), post)
		if err != nil {
			return err
		}

		post.ID = result.InsertedID.(primitive.ObjectID)

		filter := bson.M{"_id": post.UserID}
		update := bson.M{"$push": bson.M{"posts": post.ID}}
		_, err = usersCollection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			return err
		}

		// Коммит транзакции
		return session.CommitTransaction(sessionContext)
	})

	if err != nil {
		if err := session.AbortTransaction(context.TODO()); err != nil {
			log.Println(err)
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusCreated, post)
	}
}

func GetAllLikedPosts(ctx *gin.Context, client *mongo.Client) {
	userId := ctx.Param("id")

	usersCollection := client.Database("glam").Collection("users")
	postsCollection := client.Database("glam").Collection("posts")

	objID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user db.User
	filter := bson.M{"_id": objID}
	err = usersCollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	likedPosts := []db.Post{}
	for _, postID := range user.LikedPosts {
		var post db.Post
		filter := bson.M{"_id": postID}
		err := postsCollection.FindOne(context.Background(), filter).Decode(&post)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		likedPosts = append(likedPosts, post)
	}

	ctx.JSON(http.StatusOK, gin.H{"liked_posts": likedPosts})
}

func LikePost(ctx *gin.Context, client *mongo.Client) {
	userId := ctx.Param("user_id")
	postId := ctx.Param("post_id")

	// Получение коллекций "users" и "posts"
	usersCollection := client.Database("glam").Collection("users")
	postsCollection := client.Database("glam").Collection("posts")
	notificationCollection := client.Database("glam").Collection("notification")

	// Парсинг ObjectID из строк userID и postID
	userObjID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	postObjID, err := primitive.ObjectIDFromHex(postId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session, err := client.StartSession()
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
	}
	defer session.EndSession(context.TODO())

	transactionOptions := options.Transaction().SetWriteConcern(writeconcern.Majority())

	err = mongo.WithSession(context.TODO(), session, func(sessionContext mongo.SessionContext) error {
		// Начало транзакции
		err := session.StartTransaction(transactionOptions)
		if err != nil {
			return err
		}

		// Проверка, существует ли такой пользователь
		var user db.User
		filter := bson.M{"_id": userObjID}
		err = usersCollection.FindOne(context.Background(), filter).Decode(&user)
		if err != nil {
			return err
		}

		// Проверка, существует ли такой пост
		var post db.Post
		filter = bson.M{"_id": postObjID}
		err = postsCollection.FindOne(context.Background(), filter).Decode(&post)
		if err != nil {
			return err
		}

		// Проверка, лайкал ли пользователь уже этот пост
		for _, likedPostID := range user.LikedPosts {
			if likedPostID == postObjID {
				return nil
			}
		}

		// Обновление количества лайков поста
		update := bson.M{"$inc": bson.M{"likes_count": 1}}
		_, err = postsCollection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			return err
		}

		filter = bson.M{"_id": userObjID}
		update = bson.M{"$push": bson.M{"liked_posts": postObjID}}
		_, err = usersCollection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			return err
		}

		_, err = notificationCollection.InsertOne(context.Background(), db.Notification{
			Type:   "like",
			UserID: userObjID,
			PostId: postObjID,
		})
		if err != nil {
			return err
		}
		// Коммит транзакции
		return session.CommitTransaction(sessionContext)
	})
	if err != nil {
		if err := session.AbortTransaction(context.TODO()); err != nil {
			log.Println(err)
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.Status(http.StatusOK)
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func GetNotifications(ctx *gin.Context, client *mongo.Client) {
	//userId := ctx.Param("id")

	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database("glam").Collection("notification")

	changeStreamOptions := options.ChangeStream().SetFullDocument(options.UpdateLookup)
	changeStream, err := collection.Watch(c, mongo.Pipeline{}, changeStreamOptions)
	if err != nil {
		log.Println(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
	defer changeStream.Close(c)

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	log.Println("NEw connnection")
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	for {
		var changeDocument bson.M
		if changeStream.Next(ctx) {
			if err := changeStream.Decode(&changeDocument); err != nil {
				log.Println(err)
				continue
			}

			// Преобразуем документ в JSON
			jsonData, err := json.Marshal(changeDocument)
			if err != nil {
				log.Println(err)
				continue
			}

			log.Println(string(jsonData))

			// Отправляем данные по WebSocket
			if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
				log.Println(err)
				break
			}
		} else {
			log.Println("Ошибка при получении следующего изменения:", changeStream.Err())
			break
		}
	}
}
