package server

import (
	"errors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"glam/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

func createRouter(client *mongo.Client) *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:63342"},        // Разрешенные источники
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},  // Разрешенные методы
		AllowHeaders:     []string{"Content-Type", "Authorization"}, // Разрешенные заголовки
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.GET("/ping", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	router.GET("/profile/:id", func(context *gin.Context) {
		getProfile(context, client)
	})
	router.GET("/profiles", func(context *gin.Context) {
		getProfiles(context, client)
	})
	router.POST("/profile", func(context *gin.Context) {
		createOrModifyProfile(context, client)
	})

	router.GET("/posts/:id", func(context *gin.Context) {
		GetAllPostsByAuthorID(context, client)
	})
	router.GET("/all_posts", func(context *gin.Context) {
		GetAllPosts(context, client)
	})
	router.POST("/posts", func(context *gin.Context) {
		CreatePost(context, client)
	})
	router.GET("/posts/liked/:id", func(context *gin.Context) {
		GetAllLikedPosts(context, client)
	})
	router.POST("/posts/like/:user_id/:post_id", func(context *gin.Context) {
		LikePost(context, client)
	})

	router.GET("/notifications/:user_id", func(context *gin.Context) {
		GetNotifications(context, client)
	})
	return router
}

func SpawnServer(config config.Config, client *mongo.Client) *http.Server {
	router := createRouter(client)
	server := &http.Server{
		Handler: router,
		Addr:    config.Server.Address,
	}
	go func(s *http.Server) {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}(server)
	return server
}
