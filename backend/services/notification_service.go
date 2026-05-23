package services

import (
	"supplierhub-backend/config"
	"supplierhub-backend/models"

	"gorm.io/gorm"
)

func notificationDB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return config.DB
}

func CreateNotification(tx *gorm.DB, notification models.Notification) (models.Notification, error) {
	db := notificationDB(tx)
	if err := db.AutoMigrate(&models.Notification{}); err != nil {
		return notification, err
	}
	if err := db.Create(&notification).Error; err != nil {
		return notification, err
	}
	return notification, nil
}

func CreateRoleNotifications(tx *gorm.DB, role models.Role, notification models.Notification) ([]models.Notification, error) {
	db := notificationDB(tx)

	var users []models.User
	if err := db.Where("role = ?", role).Find(&users).Error; err != nil {
		return nil, err
	}

	notifications := make([]models.Notification, 0, len(users))
	for _, user := range users {
		item := notification
		item.ID = ""
		item.UserID = user.ID
		item.Role = string(role)

		created, err := CreateNotification(tx, item)
		if err != nil {
			return notifications, err
		}
		notifications = append(notifications, created)
	}

	return notifications, nil
}
