package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
)

func main() {

	err := godotenv.Load()
    if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
    }

	KAFKA_BROKER := os.Getenv("KAFKA_BROKER")


	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{KAFKA_BROKER},
		Topic: "user_registration_emails",
		GroupID: "email_service",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	fmt.Println("EMail worker started, listening to kafka")
	for{
		m, err := r.ReadMessage(context.Background())
		if err!=nil{
			log.Printf("error reading kafka message: %v", err)
			continue
		}

		email := string(m.Value)
		fmt.Printf("sending email for: %s\n", email)
		sendEmail(string(m.Value))
	}
}



func sendEmail(payload string) {
    log.Printf("📨 Sending email with payload: %s", payload)
    time.Sleep(2 * time.Second) 
    log.Println("✅ Email sent successfully!")
}

