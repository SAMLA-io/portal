package auth

import (
	"context"
	"fmt"
	"log"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/user"
)

func init() {
	globalCtx = context.Background()
	user, err := getUser("user_2Y2222222222222222222222")

	if err != nil {

		log.Fatal(err)

	}

	fmt.Println(user)
}

var globalCtx context.Context

func getUser(userId string) (*clerk.User, error) {
	userDetails, err := user.Get(globalCtx, userId)
	return userDetails, err
}
