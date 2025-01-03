package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"log"
)

// Events struct defines the structure of an event.
type Events struct {
	Name    string `json:"name"`
	Date    string `json:"date"`
	Message string `json:"message"`
}

// getUserID retrieves the user ID from the database based on the username extracted from JWT claims.
func getUserID(c *fiber.Ctx, db *sql.DB) int {
	user := c.Locals("user") // Extract the decoded JWT claims
	var username string

	// Assert and extract claims from the JWT token
	if token, ok := user.(*jwt.Token); ok {
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return 0
		}

		username = claims["username"].(string)
		log.Println("Authenticated user:", username) // Log the username for debugging purposes
	} else {
		return 0
	}

	var userID int

	// Query the database to fetch user ID for the given username
	err := db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		return 0
	}

	return userID
}

// CreateEvent handles the creation of a new event in the database.
func CreateEvent(c *fiber.Ctx, db *sql.DB) error {
	event := new(Events)
	// Parse the request body into the event struct
	if err := json.Unmarshal(c.Body(), &event); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": string(err.Error()),
		})
	}

	var userID = getUserID(c, db)

	// Prepare and execute the SQL query to insert the event
	insertQuery, err := db.Prepare("INSERT INTO events (name, message, date, user_id) VALUES(?,?,?,?)")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": string(err.Error()),
		})
	}
	defer insertQuery.Close()

	_, err = insertQuery.Exec(event.Name, event.Message, event.Date, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": string(err.Error()),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":     "created",
		"event_name": event.Name,
		"message":    "Event created successfully",
	})
}

// UpdateEvent updates the details of an existing event.
func UpdateEvent(c *fiber.Ctx, db *sql.DB) error {
	eventName := c.Params("name") // Get the event name from URL params

	newEvent := new(Events)

	// Parse the request body into the newEvent struct
	if err := json.Unmarshal(c.Body(), &newEvent); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": string(err.Error()),
		})
	}

	oldEvent := new(Events)

	var id int
	var userID = getUserID(c, db)

	// Fetch the current details of the event
	err := db.QueryRow("SELECT id, name, message, date FROM events WHERE name = ? and user_id = ?", eventName, userID).Scan(&id, &oldEvent.Name, &oldEvent.Message, &oldEvent.Date)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Record not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": string(err.Error()),
		})
	}

	// Update fields if new values are provided
	if newEvent.Name != "" {
		oldEvent.Name = newEvent.Name
	}
	if newEvent.Message != "" {
		oldEvent.Message = newEvent.Message
	}
	if newEvent.Date != "" {
		oldEvent.Date = newEvent.Date
	}

	// Prepare and execute the SQL query to update the event
	updateQuery, err := db.Prepare("UPDATE events SET name = ?, message = ?, date = ? WHERE id = ?")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": string(err.Error()),
		})
	}
	defer updateQuery.Close()

	_, err = updateQuery.Exec(oldEvent.Name, oldEvent.Message, oldEvent.Date, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": string(err.Error()),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":   "updated",
		"event_id": id,
		"message":  "Event updated successfully",
	})
}

// GetEvent retrieves the details of a specific event by name.
func GetEvent(c *fiber.Ctx, db *sql.DB) error {
	eventName := c.Params("name") // Get the event name from URL params

	event := new(Events)

	var id int
	var userID = getUserID(c, db)

	// Query the database to fetch event details
	err := db.QueryRow("SELECT id, name, message, date FROM events WHERE name = ? and user_id = ?", eventName, userID).Scan(&id, &event.Name, &event.Message, &event.Date)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": string(err.Error()),
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": string(err.Error()),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":   "fetched",
		"event_id": id,
		"details":  event,
		"message":  "Event fetched successfully",
	})
}

// DeleteEvent removes an event from the database by name.
func DeleteEvent(c *fiber.Ctx, db *sql.DB) error {
	var userID = getUserID(c, db)

	eventName := c.Params("name") // Get the event name from URL params

	// Prepare and execute the SQL query to delete the event
	deleteQuery, err := db.Prepare("DELETE FROM events WHERE name = ? and user_id = ?")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": string(err.Error()),
		})
	}
	defer deleteQuery.Close()

	result, err := deleteQuery.Exec(eventName, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": string(err.Error()),
		})
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": string(err.Error()),
		})
	}

	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Record not found",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":     "deleted",
		"event_name": eventName,
		"message":    "Event deleted successfully",
	})
}
