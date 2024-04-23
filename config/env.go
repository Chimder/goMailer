package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvVars struct {
	GOOGLE_CLIENT_ID     string
	GOOGLE_CLIENT_SECRET string
}

func LoadEnv() EnvVars {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	// redis_url := os.Getenv("REDIS_URL")
	// db_url := os.Getenv("DB_URL")

	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	// redis_db := os.Getenv("REDIS_DB")
	// parsed_redis_db, err := strconv.Atoi(redis_db)
	// if err != nil {
	// 	panic("cannot parse redis DB number")
	// }

	return EnvVars{
		GOOGLE_CLIENT_ID:     googleClientId,
		GOOGLE_CLIENT_SECRET: googleClientSecret,
	}
}
