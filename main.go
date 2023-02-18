package main

import (
	"fmt"
	"log"
	"os"

	"github.com/faruqfadhil/venue-api/core/module"
	"github.com/faruqfadhil/venue-api/handler"
	"github.com/faruqfadhil/venue-api/pkg/api"
	venueRepo "github.com/faruqfadhil/venue-api/repository/venue"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("unable to load env, err: %v", err)
	}
	db := conn()
	repo := venueRepo.New(db)
	usecase := module.New(repo)
	hdlr := handler.New(usecase)
	middlewareSvc := api.NewMiddlewareService(usecase)
	router := gin.Default()
	v1 := router.Group("/v1")
	{
		v1.GET("/city", hdlr.GetCities)
		v1.POST("/register", hdlr.Register)
		v1.POST("/login", hdlr.Login)
	}
	usingAuth := router.Group("/v1")
	usingAuth.Use(middlewareSvc.AuthenticateRequest())
	{
		// usingAuth.GET("/city", hdlr.GetCities)
	}

	router.Run(fmt.Sprintf(":%s", os.Getenv("GIN_PORT")))
}

func conn() *gorm.DB {
	defaultParams := "charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_PORT"), os.Getenv("MYSQL_DATABASE"), defaultParams)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error when try to connect db: %v", err)
	}
	return db
}
