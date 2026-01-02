package main

// Unused import removed

func main() {

	r := SetupRouter()

	// Start the server on port 8080
	r.Run(":8080")
}
