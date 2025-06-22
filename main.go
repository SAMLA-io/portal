package main

import (
	"fmt"
	"log"
	"nucleus/auth"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	clerkAPIKey := os.Getenv("CLERK_SECRET_KEY")
	if clerkAPIKey == "" {
		log.Fatal("CLERK_SECRET_KEY environment variable is required")
	}

	/* Create a new user */
	// var password string = "@D()JASNC@S23D"
	// userParams := &user.CreateParams{
	// 	EmailAddresses: &[]string{"test@test.com"},
	// 	Password:       &password,
	// }

	// user, err := auth.CreateUser(userParams)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(user)

	users, err := auth.ListUsers()
	if err != nil {
		log.Fatal(err)
	}

	for _, user := range users {
		fmt.Println(user.ID)
	}
}
