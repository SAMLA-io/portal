package auth

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"portal/schema"
)

type Subscription = schema.Subscription

var globalCtx context.Context

func init() {
	globalCtx = context.Background()
}

// VerifyUserPermissions verifies that the user has the necessary
// permissions to access the resource (has bought the product)
func VerifyUserPermissions(jwtToken string, productID string) (bool, error) {
	nucleusURL := os.Getenv("NUCLEUS_URL")

	req, err := http.NewRequest("GET", nucleusURL+"/user/subscriptions", nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("Authorization", "Bearer "+jwtToken)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return false, err
	}

	var subscriptions []Subscription
	err = json.Unmarshal(body, &subscriptions)
	if err != nil {
		return false, err
	}

	var valid = false
	for _, subscription := range subscriptions {
		if subscription.Status == "active" && subscription.ProductID == productID {
			valid = true
		}
	}

	return valid, nil
}
