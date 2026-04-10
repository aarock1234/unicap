// Package main demonstrates solving a captcha with a built-in provider.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aarock1234/unicap"
	// "github.com/aarock1234/unicap/provider/anticaptcha"
	// "github.com/aarock1234/unicap/provider/capsolver"
	"github.com/aarock1234/unicap/provider/twocaptcha"
	"github.com/aarock1234/unicap/tasks"
	"github.com/joho/godotenv"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}

}

func run() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	// Initialize the provider, in this case CapSolver, with your API key
	provider, err := twocaptcha.New(os.Getenv("2CAPTCHA_API_KEY"))
	if err != nil {
		return err
	}

	// Create a client with the provider
	client, err := unicap.New(provider)
	if err != nil {
		return err
	}

	// Set up a context with timeout for the captcha solving operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Solve the ReCaptcha V2 challenge
	// This will automatically poll until the solution is ready
	solution, err := client.Solve(ctx, &tasks.ReCaptchaV3EnterpriseTask{
		WebsiteURL: "https://id.embark.games/id/link?code=SWAZTRGB",
		WebsiteKey: "6LdHrfonAAAAALPlD3Wn6XJr4IEllwuQDDaMkxfs",
		MinScore:   0.9,
	})
	if err != nil {
		return err
	}

	fmt.Printf("token: %s\n", solution.Token)

	return nil
}
