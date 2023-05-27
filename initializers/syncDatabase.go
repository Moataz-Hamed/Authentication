package initializers

import "Auth/models"

func SyncDatabase() {
	Db.AutoMigrate(&models.User{})
}
