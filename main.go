package main

import (
	"database/sql"
	"github.com/Vansh3140/Reminder-App/database"
	"github.com/Vansh3140/Reminder-App/handlers"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Application version
const version = "1.0.0"

// Secret key for signing JWT tokens
var secretKey = []byte(os.Getenv("SECRET_KEY"))

// Credentials struct to parse login and signup requests
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	// Connect to the database
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	defer db.Close()

	// Initialize the Fiber app with the specified configuration
	app := fiber.New(fiber.Config{
		AppName: version,
	})

	// Middleware for logging HTTP requests
	app.Use(logger.New())

	// Public routes for login and signup
	app.Post("/login", func(c *fiber.Ctx) error {
		return login(c, db)
	})
	app.Post("/signup", func(c *fiber.Ctx) error {
		return signup(c, db)
	})

	// Protected API routes using JWT middleware
	api := app.Group("/api/v1")
	api.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: secretKey},
	}))

	// Event management routes (protected)
	api.Post("/event", func(c *fiber.Ctx) error {
		return handlers.CreateEvent(c, db)
	})
	api.Get("/event/:name", func(c *fiber.Ctx) error {
		return handlers.GetEvent(c, db)
	})
	api.Put("/event/:name", func(c *fiber.Ctx) error {
		return handlers.UpdateEvent(c, db)
	})
	api.Delete("/event/:name", func(c *fiber.Ctx) error {
		return handlers.DeleteEvent(c, db)
	})

	// Graceful shutdown setup
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGTSTP)

	// Start the server in a goroutine
	go func() {
		if err := app.Listen(":8080"); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Wait for a termination signal
	<-stop
	log.Println("Received shutdown signal, shutting down...")

	// Shutdown the server gracefully
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}

	log.Println("Server shutdown successfully")
}

// login handles user authentication and JWT generation
func login(c *fiber.Ctx, db *sql.DB) error {
	var creds Credentials
	if err := c.BodyParser(&creds); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	var userID int
	var storedPassword string

	// Query the database for user credentials
	err := db.QueryRow("SELECT id, password FROM users WHERE username = ?", creds.Username).Scan(&userID, &storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(400).JSON(fiber.Map{
				"status":  "error",
				"message": "No user with the given credentials exists",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// Compare the provided password with the stored hash
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(creds.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid username or password"})
	}

	// Generate and return a JWT token
	return jwtSigner(c, creds.Username)
}

// signup handles new user registration
func signup(c *fiber.Ctx, db *sql.DB) error {
	var creds Credentials
	if err := c.BodyParser(&creds); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Hash the user's password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// Insert the new user into the database
	insertQuery, err := db.Prepare("INSERT INTO users (username, password) VALUES (?, ?)")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	_, err = insertQuery.Exec(creds.Username, hashedPassword)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// Generate and return a JWT token
	return jwtSigner(c, creds.Username)
}

// jwtSigner generates a JWT token for a given username
func jwtSigner(c *fiber.Ctx, username string) error {
	// Create and sign a JWT token with user claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 365)), // Token expires in 1 year
	})

	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	// Return the signed JWT token
	return c.JSON(fiber.Map{"token": signedToken})
}
