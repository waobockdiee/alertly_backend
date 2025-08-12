package notifications

import (
	"alertly/internal/auth"
	"alertly/internal/database"
	"alertly/internal/response"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SaveDeviceToken(c *gin.Context) {

	var accountID int64
	var err error

	var req struct {
		DeviceToken string `json:"deviceToken" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error bindjson: %v", err)
		response.Send(c, http.StatusBadRequest, true, "error.", nil)
		return
	}

	accountID, err = auth.GetUserFromContext(c)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "Unauthorized", nil)
		return
	}

	repo := NewRepository(database.DB)
	if err := repo.SaveDeviceToken(accountID, req.DeviceToken); err != nil {
		log.Printf("Error bindjson: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "could not save device token", nil)
		return
	}
	response.Send(c, http.StatusOK, false, "Success", nil)
}

func DeleteDeviceToken(c *gin.Context) {

	var accountID int64
	var err error

	var req struct {
		DeviceToken string `json:"deviceToken" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error bindjson: %v", err)
		response.Send(c, http.StatusBadRequest, true, "error.", nil)
		return
	}

	accountID, err = auth.GetUserFromContext(c)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "Unauthorized", nil)
		return
	}

	repo := NewRepository(database.DB)
	if err := repo.DeleteDeviceToken(accountID, req.DeviceToken); err != nil {
		log.Printf("Error bindjson: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "could not delete device token", nil)
		return
	}

	log.Printf("no error deleting device token: TOKEN: %v ACCOUNT_ID: %v", req.DeviceToken, accountID)
	response.Send(c, http.StatusOK, false, "Success", nil)
}

// GetNotifications obtiene las notificaciones del usuario
func GetNotifications(c *gin.Context) {
	accountID, err := auth.GetUserFromContext(c)
	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "Unauthorized", nil)
		return
	}

	// Obtener parámetros de paginación
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	repo := NewRepository(database.DB)
	notifications, err := repo.GetNotifications(accountID, limit, offset)
	if err != nil {
		log.Printf("Error getting notifications: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error getting notifications", nil)
		return
	}

	response.Send(c, http.StatusOK, false, "Success", notifications)
}

// GetUnreadCount obtiene el conteo de notificaciones no leídas
func GetUnreadCount(c *gin.Context) {
	accountID, err := auth.GetUserFromContext(c)
	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "Unauthorized", nil)
		return
	}

	repo := NewRepository(database.DB)
	count, err := repo.GetUnreadCount(accountID)
	if err != nil {
		log.Printf("Error getting unread count: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error getting unread count", nil)
		return
	}

	response.Send(c, http.StatusOK, false, "Success", map[string]interface{}{
		"unread_count": count,
	})
}

// MarkAsRead marca una notificación como leída
func MarkAsRead(c *gin.Context) {
	accountID, err := auth.GetUserFromContext(c)
	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "Unauthorized", nil)
		return
	}

	var req struct {
		NotificationID int64 `json:"notification_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error bindjson: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Invalid request", nil)
		return
	}

	repo := NewRepository(database.DB)
	err = repo.MarkAsRead(accountID, req.NotificationID)
	if err != nil {
		log.Printf("Error marking as read: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error marking as read", nil)
		return
	}

	response.Send(c, http.StatusOK, false, "Success", nil)
}

// MarkAllAsRead marca todas las notificaciones como leídas
func MarkAllAsRead(c *gin.Context) {
	accountID, err := auth.GetUserFromContext(c)
	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "Unauthorized", nil)
		return
	}

	repo := NewRepository(database.DB)
	err = repo.MarkAllAsRead(accountID)
	if err != nil {
		log.Printf("Error marking all as read: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error marking all as read", nil)
		return
	}

	response.Send(c, http.StatusOK, false, "Success", nil)
}

// DeleteNotification elimina una notificación
func DeleteNotification(c *gin.Context) {
	accountID, err := auth.GetUserFromContext(c)
	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "Unauthorized", nil)
		return
	}

	var req struct {
		NotificationID int64 `json:"notification_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error bindjson: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Invalid request", nil)
		return
	}

	repo := NewRepository(database.DB)
	err = repo.DeleteNotification(accountID, req.NotificationID)
	if err != nil {
		log.Printf("Error deleting notification: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error deleting notification", nil)
		return
	}

	response.Send(c, http.StatusOK, false, "Success", nil)
}
