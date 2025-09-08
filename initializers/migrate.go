package initializers

import (
	"go-chat-app/models"
	"log"
)

func Migrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Chat{})
	if err != nil {
		log.Fatal("failed to migrate:", err)
	}
	log.Println("Database migrated successfully")
}
