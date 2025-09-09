package main

import (
	"context"
	"fmt"
	"log"
	"os"
	server "rest_with_mongo/http"
	docs_handlers "rest_with_mongo/http/handlers/docshandlers"
	user_handlers "rest_with_mongo/http/handlers/userhandlers"
	"rest_with_mongo/repository/kafka"
	"rest_with_mongo/repository/mongodb"
	docs_repo "rest_with_mongo/repository/mongodb/docs"
	user_repo "rest_with_mongo/repository/mongodb/users"
	docs_service "rest_with_mongo/services/docs"
	user_services "rest_with_mongo/services/users"
	"rest_with_mongo/utils/jwt"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"
)

func InitializeServer(db *mongo.Database, logger *zap.Logger, producer *kafka.Producer) error{
	userRepo := user_repo.NewUserRepository(db)
	jwtGen := &jwt.ApxJwt{SecretKey: "abasdfadfavadga"}
	userService := user_services.NewUserService(userRepo, jwtGen, producer)
	userHanlder := user_handlers.NewUserHandler(userService, jwtGen) 
	docsRepo := docs_repo.NewDocsRepo(db)
	docsService := docs_service.NewDocsService(docsRepo)
	docsHandler := docs_handlers.NewDocsHandler(docsService)
	

	
	server := server.NewServer(logger, userHanlder, docsHandler, jwtGen)
	return server.Start(":8080")
}


func main() {
	fmt.Println("Welcome to REST API with MongoDB")
	err := godotenv.Load()
	KAFKA_BROKER:= os.Getenv("KAFKA_BROKER")
    if err != nil {
    		log.Fatalf("Error loading .env file: %v", err)
	}
	producer := kafka.NewProducer([]string{KAFKA_BROKER}, "user_registration_emails")
	defer producer.Close()
	logger,err := zap.NewProduction()
	
	ctx , cancel := context.WithTimeout(context.Background(), 10* time.Second) 
	defer cancel()

	db, err := mongodb.Connect(ctx)
	if( err !=nil){
		panic(err)
	}
	
	if err := InitializeServer(db, logger, producer); err != nil {
		logger.Fatal("server failed to start", zap.Error(err))
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
		
}