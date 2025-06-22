package main

import (
	"fmt"
	"log"
	"nucleus/auth"
)

func main() {
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

	// users, err := auth.ListUsers()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// for _, user := range users {
	// 	fmt.Println(user.ID)
	// }

	user, err := auth.GetUser("user_2ysC4ZWuAK2UPR84855DhG0w5v1")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(user.EmailAddresses[0].EmailAddress)
	fmt.Println(*user.FirstName)
	fmt.Println(*user.LastName)
}
