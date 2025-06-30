package auth

import (
	"context"
	"encoding/json"
	"log"
	"portal/stripe"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/user"
)

var globalCtx context.Context

func init() {
	globalCtx = context.Background()
}

func getUser(userId string) (*clerk.User, error) {
	userDetails, err := user.Get(globalCtx, userId)
	return userDetails, err
}

// VerifyUserPermissions verifies that the user has the necessary
// permissions to access the resource (has bought the product)
func VerifyUserPermissions(userId string) (bool, error) {
	userDetails, err := getUser(userId)
	if err != nil {
		return false, err
	}

	jsonData, err := userDetails.PublicMetadata.MarshalJSON()
	if err != nil {
		return false, err
	}

	var metadata map[string]interface{}
	err = json.Unmarshal(jsonData, &metadata)
	if err != nil {
		return false, err
	}

	productsId := metadata["stripe"].(map[string]interface{})["products_id"].([]interface{})

	// TODO, this only checks if the user's product exists in Stripe (by comparing the user's product id with the product id in Stripe)
	// It does not check if the user has access to this specific product (by checking the user's subscription status or by providing the desired product id and comparing it with the user's product id)
	for _, productId := range productsId {
		product, err := stripe.GetProduct(productId.(string))
		if err != nil {
			return false, err
		}

		if product.ID == productId {
			log.Println("Product found, user has access to this product")
		} else {
			log.Println("Product not found, user does not have access to this product")
		}
	}

	return true, nil
}
