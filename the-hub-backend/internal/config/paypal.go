package config

import (
	"context"
	"os"

	"github.com/plutov/paypal/v4"
)

var paypalClient *paypal.Client

// InitPayPal initializes the PayPal client
func InitPayPal() error {
	clientID := os.Getenv("PAYPAL_CLIENT_ID")
	clientSecret := os.Getenv("PAYPAL_CLIENT_SECRET")
	isSandbox := os.Getenv("PAYPAL_SANDBOX") == "true"

	if clientID == "" || clientSecret == "" {
		Logger.Warn("PayPal credentials not configured")
		return nil // Don't fail if PayPal is not configured
	}

	var err error
	paypalClient, err = paypal.NewClient(clientID, clientSecret, paypal.APIBaseSandBox)
	if isSandbox {
		paypalClient, err = paypal.NewClient(clientID, clientSecret, paypal.APIBaseSandBox)
	} else {
		paypalClient, err = paypal.NewClient(clientID, clientSecret, paypal.APIBaseLive)
	}

	if err != nil {
		return err
	}

	// Get access token
	_, err = paypalClient.GetAccessToken(context.Background())
	if err != nil {
		return err
	}

	Logger.Info("PayPal client initialized successfully")
	return nil
}

// GetPayPalClient returns the PayPal client instance
func GetPayPalClient() *paypal.Client {
	return paypalClient
}
