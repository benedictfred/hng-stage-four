package repository

import (
	"context"
	"push-service/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type DeviceRepository interface {
	Create(ctx context.Context, device *models.Device) error
	GetByToken(ctx context.Context, token string) (*models.Device, error)
	GetByUserID(ctx context.Context, userID string) ([]models.Device, error)
	UpdateStatus(ctx context.Context, token string, isActive bool) error
	Delete(ctx context.Context, token string) error
}

type deviceRepo struct {
	db *pgxpool.Pool
}

func NewDeviceRepository(db *pgxpool.Pool) DeviceRepository {
	return &deviceRepo{db: db}
}

func (r *deviceRepo) Create(ctx context.Context, device *models.Device) error {
	query := `
		INSERT INTO devices (user_id, token, platform, is_active)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		device.UserID,
		device.Token,
		device.Platform,
		device.IsActive,
	).Scan(&device.ID, &device.CreatedAt, &device.UpdatedAt)

	if err != nil {
		zap.L().Error("Failed to create device", zap.Error(err))
		return err
	}

	return nil
}

func (r *deviceRepo) GetByToken(ctx context.Context, token string) (*models.Device, error) {
	query := `
		SELECT id, user_id, token, platform, is_active, created_at, updated_at
		FROM devices
		WHERE token = $1 AND is_active = true
	`

	var device models.Device
	err := r.db.QueryRow(ctx, query, token).Scan(
		&device.ID,
		&device.UserID,
		&device.Token,
		&device.Platform,
		&device.IsActive,
		&device.CreatedAt,
		&device.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		zap.L().Error("Failed to get device by token", zap.Error(err))
		return nil, err
	}

	return &device, nil
}

func (r *deviceRepo) GetByUserID(ctx context.Context, userID string) ([]models.Device, error) {
	query := `
		SELECT id, user_id, token, platform, is_active, created_at, updated_at
		FROM devices
		WHERE user_id = $1 AND is_active = true
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		zap.L().Error("Failed to get devices by user ID", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var devices []models.Device
	for rows.Next() {
		var device models.Device
		err := rows.Scan(
			&device.ID,
			&device.UserID,
			&device.Token,
			&device.Platform,
			&device.IsActive,
			&device.CreatedAt,
			&device.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}

	return devices, nil
}

func (r *deviceRepo) UpdateStatus(ctx context.Context, token string, isActive bool) error {
	query := `
		UPDATE devices 
		SET is_active = $1, updated_at = NOW()
		WHERE token = $2
	`

	result, err := r.db.Exec(ctx, query, isActive, token)
	if err != nil {
		zap.L().Error("Failed to update device status", zap.Error(err))
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *deviceRepo) Delete(ctx context.Context, token string) error {
	query := `DELETE FROM devices WHERE token = $1`

	result, err := r.db.Exec(ctx, query, token)
	if err != nil {
		zap.L().Error("Failed to delete device", zap.Error(err))
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}