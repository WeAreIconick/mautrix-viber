package utils_test

import (
	"fmt"
	"github.com/example/mautrix-viber/internal/utils"
)

func ExampleValidateMatrixUserID() {
	// Valid Matrix user ID
	err := utils.ValidateMatrixUserID("@alice:matrix.org")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Valid user ID")
	}

	// Invalid Matrix user ID (missing @)
	err = utils.ValidateMatrixUserID("alice:matrix.org")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	// Output:
	// Valid user ID
	// Error: invalid Matrix user ID format: alice:matrix.org
}

func ExampleValidateMatrixRoomID() {
	// Valid Matrix room ID
	err := utils.ValidateMatrixRoomID("!abc123:matrix.org")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Valid room ID")
	}

	// Invalid Matrix room ID (missing !)
	err = utils.ValidateMatrixRoomID("#abc123:matrix.org")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	// Output:
	// Valid room ID
	// Error: invalid Matrix room ID format: #abc123:matrix.org
}

func ExampleValidateURL() {
	// Valid HTTPS URL
	err := utils.ValidateURL("https://example.com/webhook")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Valid URL")
	}

	// Invalid URL (no scheme)
	err = utils.ValidateURL("example.com")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	// Output:
	// Valid URL
	// Error: URL must use http or https scheme
}
