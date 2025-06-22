/*
Written by Juan Pablo Gutierrez
Package auth provides a wrapper for the Clerk API.
*/
package auth

import (
	"context"
	"log"
	"os"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/joho/godotenv"
)

var globalCtx context.Context

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	globalCtx = context.Background()
	clerkAPIKey := os.Getenv("CLERK_SECRET_KEY")
	if clerkAPIKey == "" {
		log.Fatal("CLERK_SECRET_KEY environment variable is required")
	}

	clerk.SetKey(clerkAPIKey)
}

func CreateUser(userParams *user.CreateParams) (*clerk.User, error) {
	newUser, err := user.Create(globalCtx, userParams)

	return newUser, err
}

func GetUser(userId string) (*clerk.User, error) {
	userDetails, err := user.Get(globalCtx, userId)
	return userDetails, err
}

func ListUsers() ([]*clerk.User, error) {
	users, err := user.List(globalCtx, &user.ListParams{})
	return users.Users, err
}

func DeleteUser(userId string) error {
	_, err := user.Delete(globalCtx, userId)
	return err
}
