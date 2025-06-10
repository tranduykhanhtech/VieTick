package config

import (
    "fmt"
    "log"
    "os"

    "github.com/joho/godotenv"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

var DB *gorm.DB

// DatabaseConfig holds the database configuration
type DatabaseConfig struct {
    Host     string
    Port     string
    Username string
    Password string
    Database string
}

// LoadDatabaseConfig loads database configuration from environment variables
func LoadDatabaseConfig() (*DatabaseConfig, error) {
    // Chỉ load .env file khi môi trường là development
    if os.Getenv("ENV") == "development" {
        if err := godotenv.Load(); err != nil {
            log.Printf("Error loading .env file: %v", err)
            // Continue without .env if not found, rely on system env vars
        }
    }

    config := &DatabaseConfig{
        Host:     getEnv("DB_HOST", "localhost"),
        Port:     getEnv("DB_PORT", "3306"),
        Username: getEnv("DB_USERNAME", "root"),
        Password: getEnv("DB_PASSWORD", ""),
        Database: getEnv("DB_DATABASE", "vietick"),
    }

    if config.Username == "" || config.Password == "" || config.Database == "" {
        return nil, fmt.Errorf("database configuration is incomplete")
    }

    return config, nil
}

// InitDB initializes the database connection
func InitDB() error {
    config, err := LoadDatabaseConfig()
    if err != nil {
        return fmt.Errorf("failed to load database config: %v", err)
    }

    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&tls=true",
        config.Username,
        config.Password,
        config.Host,
        config.Port,
        config.Database,
    )

    log.Printf("Connecting to database at %s:%s...", config.Host, config.Port)
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return fmt.Errorf("failed to connect to database: %v", err)
    }
    log.Printf("Successfully connected to database %s", config.Database)

    // Set connection pool settings
    sqlDB, err := db.DB()
    if err != nil {
        return fmt.Errorf("failed to get database instance: %v", err)
    }

    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)

    // Test connection
    if err := sqlDB.Ping(); err != nil {
        return fmt.Errorf("failed to ping database: %v", err)
    }
    log.Println("Database connection test successful")

    DB = db
    return nil
}

// CloseDB closes the database connection
func CloseDB() error {
    if DB != nil {
        sqlDB, err := DB.DB()
        if err != nil {
            return fmt.Errorf("failed to get database instance: %v", err)
        }
        if err := sqlDB.Close(); err != nil {
            return fmt.Errorf("failed to close database connection: %v", err)
        }
        log.Println("Database connection closed")
    }
    return nil
}

func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
} 