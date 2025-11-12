package service

import (
	"context"
	"fmt"
	"push-service/internal/config"
	"push-service/internal/models"
	"push-service/internal/platform/fcm"
	"push-service/internal/repository"

	"go.uber.org/zap"
)

type DeviceService interface {
	RegisterDevice(ctx context.Context, req models.CreateDeviceRequest) (*models.DeviceResponse, error)
	UnregisterDevice(ctx context.Context, token string) error
	GetUserDevices(ctx context.Context, userID string) ([]models.DeviceResponse, error)
}

type deviceService struct {
	deviceRepo repository.DeviceRepository
	fcmClient  fcm.FCMClient
	cfg        *config.Config
}

func NewDeviceService(deviceRepo repository.DeviceRepository, fcmClient fcm.FCMClient, cfg *config.Config) DeviceService {
	return &deviceService{
		deviceRepo: deviceRepo,
		fcmClient:  fcmClient,
		cfg:        cfg,
	}
}

func (s *deviceService) RegisterDevice(ctx context.Context, req models.CreateDeviceRequest) (*models.DeviceResponse, error) {
	// Validate token if validation is enabled
	if s.cfg != nil && s.cfg.Queue.Validation.Enabled && s.fcmClient != nil {
		if err := s.fcmClient.ValidateToken(ctx, req.Token); err != nil {
			zap.L().Warn("Token validation failed during device registration",
				zap.String("user_id", req.UserID),
				zap.String("platform", req.Platform),
				zap.String("token", maskToken(req.Token)),
				zap.Error(err),
			)
			return nil, fmt.Errorf("token validation failed: %w", err)
		}
		zap.L().Debug("Token validated successfully",
			zap.String("user_id", req.UserID),
			zap.String("platform", req.Platform),
		)
	}

	// Check if device already exists
	existingDevice, err := s.deviceRepo.GetByToken(ctx, req.Token)
	if err != nil {
		return nil, err
	}

	if existingDevice != nil {
		// Update existing device
		if err := s.deviceRepo.UpdateStatus(ctx, req.Token, true); err != nil {
			return nil, err
		}
		return &models.DeviceResponse{
			ID:       existingDevice.ID,
			UserID:   existingDevice.UserID,
			Token:    existingDevice.Token,
			Platform: existingDevice.Platform,
			IsActive: true,
		}, nil
	}

	// Create new device
	device := &models.Device{
		UserID:   req.UserID,
		Token:    req.Token,
		Platform: req.Platform,
		IsActive: true,
	}

	if err := s.deviceRepo.Create(ctx, device); err != nil {
		return nil, err
	}

	zap.L().Info("Device registered successfully",
		zap.String("user_id", req.UserID),
		zap.String("platform", req.Platform),
	)

	return &models.DeviceResponse{
		ID:       device.ID,
		UserID:   device.UserID,
		Token:    device.Token,
		Platform: device.Platform,
		IsActive: device.IsActive,
	}, nil
}

// maskToken masks a token for logging
func maskToken(token string) string {
	if len(token) <= 20 {
		return "***"
	}
	return token[:10] + "..." + token[len(token)-10:]
}

func (s *deviceService) UnregisterDevice(ctx context.Context, token string) error {
	// Soft delete by setting is_active to false
	err := s.deviceRepo.UpdateStatus(ctx, token, false)
	if err != nil {
		zap.L().Error("Failed to unregister device", 
			zap.String("token", token), 
			zap.Error(err),
		)
		return err
	}

	zap.L().Info("Device unregistered successfully", zap.String("token", token))
	return nil
}

func (s *deviceService) GetUserDevices(ctx context.Context, userID string) ([]models.DeviceResponse, error) {
	devices, err := s.deviceRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]models.DeviceResponse, len(devices))
	for i, device := range devices {
		responses[i] = models.DeviceResponse{
			ID:       device.ID,
			UserID:   device.UserID,
			Token:    device.Token,
			Platform: device.Platform,
			IsActive: device.IsActive,
		}
	}

	return responses, nil
}