package config

import (
    "fmt"
    "log"
    "net/url"
    "os"
    "time"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
)


var JWTSecret []byte

func LoadJWTSecret() {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        log.Fatal("JWT_SECRET is not set in environment")
    }
    JWTSecret = []byte(secret)
}

var DB *gorm.DB

func ConnectDB() {
    dbUser := os.Getenv("POSTGRES_USER") 
    dbPass := os.Getenv("POSTGRES_PASSWORD")
    dbName := os.Getenv("POSTGRES_DB")
    dbPort := os.Getenv("PG_PORT")
    dbHost := os.Getenv("POSTGRES_HOST")

    // Validate required environment variables
    if dbHost == "" {
        log.Fatal("POSTGRES_HOST is required in .env or environment variables")
    }
    if dbUser == "" {
        log.Fatal("POSTGRES_USER is required in .env or environment variables")
    }
    if dbPass == "" {
        log.Fatal("POSTGRES_PASSWORD is required in .env or environment variables")
    }
    if dbName == "" {
        log.Fatal("POSTGRES_DB is required in .env or environment variables")
    }
    if dbPort == "" {
        dbPort = "5432" // Default PostgreSQL port
    }

	log.Printf("Connecting to DB -> host=%s port=%s user=%s dbname=%s", dbHost, dbPort, dbUser, dbName)

    // URL encode password and other components to handle special characters like @, #, etc.
    // Using connection URI format which handles special characters more reliably
    encodedUser := url.QueryEscape(dbUser)
    encodedPassword := url.QueryEscape(dbPass)
    encodedDBName := url.QueryEscape(dbName)
    encodedHost := url.QueryEscape(dbHost)

    // Build connection URI with connection timeout for remote databases
    // Format: postgres://user:password@host:port/dbname?sslmode=disable&connect_timeout=10
    dsn := fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=10",
        encodedUser, encodedPassword, encodedHost, dbPort, encodedDBName,
    )

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to DB:", err)
    }

    // Configure connection pool for better performance and scalability
    sqlDB, err := db.DB()
    if err != nil {
        log.Fatal("Failed to get underlying sql.DB:", err)
    }

    // SetMaxIdleConns sets the maximum number of connections in the idle connection pool
    sqlDB.SetMaxIdleConns(10)
    
    // SetMaxOpenConns sets the maximum number of open connections to the database
    sqlDB.SetMaxOpenConns(100)
    
    // SetConnMaxLifetime sets the maximum amount of time a connection may be reused
    sqlDB.SetConnMaxLifetime(time.Hour)
    
    // Set connection timeout for establishing new connections
    sqlDB.SetConnMaxIdleTime(5 * time.Minute)

    DB = db
    log.Println("Database connection pool configured successfully")
}

func AutoMigrate() {
    DB.AutoMigrate(&models.Role{}, &models.User{})
}
