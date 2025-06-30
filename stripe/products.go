package stripe

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/product"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	stripe.Key = os.Getenv("STRIPE_KEY")
}

func GetProduct(productId string) (*stripe.Product, error) {
	product, err := product.Get(productId, nil)
	if err != nil {
		return nil, err
	}

	return product, nil
}
