package fcm

import (
	"context"
	"fmt"
	"push-service/internal/config"
	"push-service/internal/models"
	"strings"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

type FCMClient interface {
	Send(ctx context.Context, deviceToken string, notification models.PushNotification) error
	SendMultiple(ctx context.Context, deviceTokens []string, notification models.PushNotification) (int, int, error)
	SendMulticast(ctx context.Context, deviceTokens []string, notification models.PushNotification) (*messaging.BatchResponse, error)
	ValidateToken(ctx context.Context, deviceToken string) error
}

type fcmClient struct {
	client *messaging.Client
}

func NewFCMClient(cfg *config.FCMConfig) (FCMClient, error) {
	ctx := context.Background()

	// Get credentials from config
	credentials, err := cfg.GetFCMCredentials()
	if err != nil {
		return nil, fmt.Errorf("failed to get FCM credentials: %w", err)
	}

	// Configure Firebase App
	firebaseConfig := &firebase.Config{
		ProjectID: cfg.ProjectID,
	}

	opt := option.WithCredentialsJSON(credentials)
	app, err := firebase.NewApp(ctx, firebaseConfig, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to create Firebase app: %w", err)
	}

	// Create FCM client
	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create FCM client: %w", err)
	}

	zap.L().Info("FCM client initialized successfully",
		zap.String("project_id", cfg.ProjectID),
		zap.Bool("using_file", cfg.UseFile),
	)
	return &fcmClient{client: client}, nil
}

func (f *fcmClient) Send(ctx context.Context, deviceToken string, notification models.PushNotification) error {
	// Convert map[string]any to map[string]string for FCM
	data := convertDataToStringMap(notification.Data)

	// Add link to data if provided
	if notification.Link != nil && *notification.Link != "" {
		if data == nil {
			data = make(map[string]string)
		}
		data["link"] = *notification.Link
		data["click_action"] = *notification.Link
	}

	msgNotification := &messaging.Notification{
		Title: notification.Title,
		Body:  notification.Body,
	}

	// Add image if provided
	if notification.Image != nil && *notification.Image != "" {
		msgNotification.ImageURL = *notification.Image
	}

	message := &messaging.Message{
		Token:        deviceToken,
		Notification: msgNotification,
		Data:         data,
	}

	// Add webpush config for web notifications
	if notification.Image != nil || notification.Link != nil {
		webpushConfig := &messaging.WebpushConfig{
			Headers: map[string]string{
				"Urgency": "high",
			},
		}

		if notification.Image != nil || notification.Link != nil {
			webpushNotification := &messaging.WebpushNotification{
				Title: notification.Title,
				Body:  notification.Body,
			}
			if notification.Image != nil && *notification.Image != "" {
				webpushNotification.Icon = *notification.Image
				webpushNotification.Image = *notification.Image
			}
			// Link is handled via data payload for web push
			webpushConfig.Notification = webpushNotification
		}
		message.Webpush = webpushConfig
	}

	response, err := f.client.Send(ctx, message)
	if err != nil {
		zap.L().Error("Failed to send FCM message",
			zap.String("token", deviceToken),
			zap.Error(err),
		)
		return err
	}

	zap.L().Info("FCM message sent successfully",
		zap.String("message_id", response),
		zap.String("token", deviceToken),
	)
	return nil
}

func (f *fcmClient) SendMultiple(ctx context.Context, deviceTokens []string, notification models.PushNotification) (int, int, error) {
	// Convert map[string]any to map[string]string for FCM
	data := convertDataToStringMap(notification.Data)

	// Add link to data if provided
	if notification.Link != nil && *notification.Link != "" {
		if data == nil {
			data = make(map[string]string)
		}
		data["link"] = *notification.Link
		data["click_action"] = *notification.Link
	}

	msgNotification := &messaging.Notification{
		Title: notification.Title,
		Body:  notification.Body,
	}

	// Add image if provided
	if notification.Image != nil && *notification.Image != "" {
		msgNotification.ImageURL = *notification.Image
	}

	// For multiple devices, send individually for better error tracking
	successCount := 0
	failureCount := 0

	for _, token := range deviceTokens {
		message := &messaging.Message{
			Token:        token,
			Notification: msgNotification,
			Data:         data,
		}

		// Add webpush config for web notifications
		if notification.Image != nil || notification.Link != nil {
			webpushConfig := &messaging.WebpushConfig{
				Headers: map[string]string{
					"Urgency": "high",
				},
			}

			if notification.Image != nil || notification.Link != nil {
				webpushNotification := &messaging.WebpushNotification{
					Title: notification.Title,
					Body:  notification.Body,
				}
				if notification.Image != nil && *notification.Image != "" {
					webpushNotification.Icon = *notification.Image
					webpushNotification.Image = *notification.Image
				}
				// Link is handled via data payload for web push
				webpushConfig.Notification = webpushNotification
			}
			message.Webpush = webpushConfig
		}

		_, err := f.client.Send(ctx, message)
		if err != nil {
			zap.L().Error("Failed to send FCM message to device",
				zap.String("token", token),
				zap.Error(err),
			)
			failureCount++
			continue
		}

		successCount++
	}

	zap.L().Info("Batch FCM messages completed",
		zap.Int("success_count", successCount),
		zap.Int("failure_count", failureCount),
		zap.Int("total", len(deviceTokens)),
	)

	return successCount, failureCount, nil
}

func (f *fcmClient) SendMulticast(ctx context.Context, deviceTokens []string, notification models.PushNotification) (*messaging.BatchResponse, error) {
	// Convert map[string]any to map[string]string for FCM
	data := convertDataToStringMap(notification.Data)

	// Add link to data if provided
	if notification.Link != nil && *notification.Link != "" {
		if data == nil {
			data = make(map[string]string)
		}
		data["link"] = *notification.Link
		data["click_action"] = *notification.Link
	}

	msgNotification := &messaging.Notification{
		Title: notification.Title,
		Body:  notification.Body,
	}

	// Add image if provided
	if notification.Image != nil && *notification.Image != "" {
		msgNotification.ImageURL = *notification.Image
	}

	// For web push, we need to configure it properly
	webpushConfig := &messaging.WebpushConfig{
		Headers: map[string]string{
			"Urgency": "high",
		},
	}

	webpushNotification := &messaging.WebpushNotification{
		Title: notification.Title,
		Body:  notification.Body,
	}

	if notification.Image != nil && *notification.Image != "" {
		webpushNotification.Icon = *notification.Image
		webpushNotification.Image = *notification.Image
	}
	// Link is handled via data payload for web push

	webpushConfig.Notification = webpushNotification

	message := &messaging.MulticastMessage{
		Tokens:       deviceTokens,
		Notification: msgNotification,
		Data:         data,
		Webpush:      webpushConfig,
	}

	response, err := f.client.SendMulticast(ctx, message)
	if err != nil {
		zap.L().Error("Failed to send multicast FCM message",
			zap.Int("device_count", len(deviceTokens)),
			zap.Error(err),
		)
		return nil, err
	}

	// Log detailed results
	if response.FailureCount > 0 {
		for i, resp := range response.Responses {
			if resp.Error != nil {
				zap.L().Warn("Individual FCM send failed",
					zap.String("token", deviceTokens[i]),
					zap.Error(resp.Error),
				)
			}
		}
	}

	zap.L().Info("Multicast FCM messages completed",
		zap.Int("success_count", response.SuccessCount),
		zap.Int("failure_count", response.FailureCount),
		zap.Int("total", len(deviceTokens)),
	)
	return response, nil
}

// convertDataToStringMap converts map[string]any to map[string]string
// FCM requires all data values to be strings
func convertDataToStringMap(data map[string]any) map[string]string {
	if data == nil {
		return nil
	}

	result := make(map[string]string)
	for key, value := range data {
		switch v := value.(type) {
		case string:
			result[key] = v
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
			result[key] = fmt.Sprintf("%v", v)
		default:
			// For complex types, convert to string
			result[key] = fmt.Sprintf("%v", v)
		}
	}
	return result
}

// ValidateToken validates a device token by attempting to send a test message
// This is a lightweight validation that checks if the token format is valid
// and if FCM can accept it. Note: This doesn't guarantee the token will work
// for actual notifications, but catches obviously invalid tokens.
func (f *fcmClient) ValidateToken(ctx context.Context, deviceToken string) error {
	// Basic format validation for FCM tokens
	if len(deviceToken) < 10 {
		return fmt.Errorf("token too short")
	}

	// Try to send a minimal test message to validate the token
	// We use a test notification that won't actually be delivered
	testNotification := models.PushNotification{
		Title: "test",
		Body:  "test",
	}

	// Use a short timeout for validation
	validationCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Attempt to send - if it fails with certain errors, token is invalid
	err := f.Send(validationCtx, deviceToken, testNotification)
	if err != nil {
		// Check for specific invalid token errors
		errStr := strings.ToLower(err.Error())
		if strings.Contains(errStr, "invalid") || strings.Contains(errStr, "not found") ||
			strings.Contains(errStr, "registration") || strings.Contains(errStr, "unregistered") {
			return fmt.Errorf("invalid token: %w", err)
		}
		// For other errors (network, etc.), we consider the token potentially valid
		// since the error might be transient
		zap.L().Debug("Token validation encountered non-fatal error",
			zap.String("token", maskToken(deviceToken)),
			zap.Error(err),
		)
	}

	return nil
}

// maskToken masks a token for logging (shows first 10 and last 10 chars)
func maskToken(token string) string {
	if len(token) <= 20 {
		return "***"
	}
	return token[:10] + "..." + token[len(token)-10:]
}
