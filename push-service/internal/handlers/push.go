package handlers

import (
	"context"
	"net/http"
	"push-service/internal/models"
	"push-service/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PushHandler struct {
	pushService service.PushService
}

func NewPushHandler(pushService service.PushService) *PushHandler {
	return &PushHandler{pushService: pushService}
}

// SendPush godoc
// @Summary Send push notification
// @Description Send a push notification to a user's devices via RabbitMQ queue
// @Tags push
// @Accept json
// @Produce json
// @Param request body models.SendPushRequest true "Push notification request"
// @Success 200 {object} map[string]string "Push notification enqueued successfully"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Failed to send push notification"
// @Router /v1/push/send [post]
func (h *PushHandler) SendPush(c *gin.Context) {
	var req models.SendPushRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Warn("Invalid push request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	if err := h.pushService.SendPush(c.Request.Context(), req); err != nil {
		zap.L().Error("Failed to send push", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to send push notification",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Push notification sent successfully",
		"user_id": req.UserID,
	})
}

// SendBulkPush godoc
// @Summary Send bulk push notifications
// @Description Send push notifications to multiple users via RabbitMQ queue
// @Tags push
// @Accept json
// @Produce json
// @Param request body models.BulkPushRequest true "Bulk push notification request"
// @Success 200 {object} map[string]interface{} "Bulk push notifications enqueued successfully"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Failed to send bulk push notifications"
// @Router /v1/push/send-bulk [post]
func (h *PushHandler) SendBulkPush(c *gin.Context) {
	var req models.BulkPushRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Warn("Invalid bulk push request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	if err := h.pushService.SendBulkPush(c.Request.Context(), req); err != nil {
		zap.L().Error("Failed to send bulk push", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send bulk push notifications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Bulk push notifications sent successfully",
		"user_count": len(req.UserIDs),
	})
}

// GetQueueStats godoc
// @Summary Get queue statistics
// @Description Get statistics for all push notification queues (main, retry, dead letter)
// @Tags queue
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Queue statistics"
// @Failure 500 {object} map[string]string "Failed to get queue statistics"
// @Router /v1/queue/stats [get]
func (h *PushHandler) GetQueueStats(c *gin.Context) {
	stats, err := h.pushService.GetQueueStats(c.Request.Context())
	if err != nil {
		zap.L().Error("Failed to get queue stats", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get queue statistics",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"queues": stats,
	})
}

// TestDirectSend godoc
// @Summary Test direct FCM send
// @Description Send a test push notification directly via FCM (bypasses queue, for testing only)
// @Tags push
// @Accept json
// @Produce json
// @Param request body object true "Direct send request" example({"token":"fcm_token","title":"Test","body":"Test message"})
// @Success 200 {object} map[string]string "FCM test message sent successfully"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "FCM send failed"
// @Router /v1/push/test-direct [post]
func (h *PushHandler) TestDirectSend(c *gin.Context) {
	var req struct {
		Token string         `json:"token" binding:"required"`
		Title string         `json:"title" binding:"required"`
		Body  string         `json:"body" binding:"required"`
		Data  map[string]any `json:"data,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Warn("Invalid direct send request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	zap.L().Info("🔧 Testing FCM direct send",
		zap.String("token", req.Token),
		zap.String("title", req.Title),
	)

	notification := models.PushNotification{
		Title: req.Title,
		Body:  req.Body,
		Data:  req.Data,
	}

	// Use the FCM client directly
	err := h.pushService.(interface {
		SendDirect(ctx context.Context, token string, notification models.PushNotification) error
	}).SendDirect(c.Request.Context(), req.Token, notification)

	if err != nil {
		zap.L().Error("💥 FCM direct send failed",
			zap.String("token", req.Token),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "FCM send failed",
			"details": err.Error(),
		})
		return
	}

	zap.L().Info("✅ FCM direct send successful")
	c.JSON(http.StatusOK, gin.H{
		"message": "FCM test message sent successfully",
	})
}
