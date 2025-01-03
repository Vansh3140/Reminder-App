package database

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
	"os"
	"time"
)

// initTLS initializes a custom TLS configuration for secure database connections.
// It loads the CA certificate from the environment variable and registers a custom TLS config with MySQL driver.
func initTLS() error {
	rootCertPool := x509.NewCertPool()
	// Append Aiven CA certificate to the root certificate pool
	if ok := rootCertPool.AppendCertsFromPEM([]byte(AivenCA)); !ok {
		return fmt.Errorf("failed to append Aiven CA certificate")
	}

	tlsConfig := &tls.Config{
		RootCAs: rootCertPool,
	}

	// Register the TLS configuration with the MySQL driver
	return mysql.RegisterTLSConfig("custom", tlsConfig)
}

// AivenCA holds the database's CA certificate loaded from an environment variable
const AivenCA = os.Getenv("CERTIFICATE")

// Connect establishes a connection to the MySQL database, configures connection settings,
// and ensures the required tables are created.
func Connect() (*sql.DB, error) {
	// Initialize TLS for secure database connections
	if err := initTLS(); err != nil {
		log.Fatalf("Failed to initialize TLS: %v", err)
	}

	// Retrieve database credentials from environment variables
	dsn := os.Getenv("DB_CREDS")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
		return nil, err
	}

	// Configure database connection settings
	db.SetConnMaxLifetime(time.Minute * 3) // Maximum connection lifetime
	db.SetMaxOpenConns(10)                 // Maximum number of open connections
	db.SetMaxIdleConns(10)                 // Maximum number of idle connections

	// Create the users table if it does not exist
	createUserSQL := `CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL
	);`
	_, err = db.Exec(createUserSQL)
	if err != nil {
		log.Fatal("Error creating users table: ", err)
	}

	// Create the events table with a foreign key reference to the users table
	createTableSQL := `CREATE TABLE IF NOT EXISTS events (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		message TEXT NOT NULL,
		date VARCHAR(255) NOT NULL,
		user_id INT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		UNIQUE (name, user_id)
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Error creating events table: ", err)
	}

	return db, nil
}
