package repository

import (
	"context"
	"database/sql"
	"fmt"

	pbd "github.com/danielmoisa/bolt-app/pkg/proto/driver"
	"github.com/danielmoisa/bolt-app/services/trip-service/internal/domain"
)

type tripRepository struct {
	db *sql.DB
}

func NewTripRepository(db *sql.DB) *tripRepository {
	return &tripRepository{db: db}
}

func (r *tripRepository) CreateTrip(ctx context.Context, trip *domain.TripModel) (*domain.TripModel, error) {
	query := `
        INSERT INTO trips (user_id, status, created_at, updated_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		trip.UserID,
		trip.Status,
		trip.CreatedAt,
		trip.UpdatedAt,
	).Scan(&trip.ID)

	if err != nil {
		return nil, fmt.Errorf("failed to create trip: %w", err)
	}

	return trip, nil
}

func (r *tripRepository) GetTripByID(ctx context.Context, id string) (*domain.TripModel, error) {
	query := `
        SELECT id, user_id, status, created_at, updated_at
        FROM trips 
        WHERE id = $1`

	var trip domain.TripModel

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&trip.ID,
		&trip.UserID,
		&trip.Status,
		&trip.CreatedAt,
		&trip.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("trip not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get trip: %w", err)
	}

	return &trip, nil
}

func (r *tripRepository) UpdateTrip(ctx context.Context, tripID string, status string, driver *pbd.Driver) error {
	var query string
	var args []interface{}

	if driver != nil {
		query = `
            UPDATE trips 
            SET status = $1, driver_id = $2, driver_name = $3, driver_car_plate = $4, updated_at = NOW()
            WHERE id = $5`
		args = []interface{}{status, driver.Id, driver.Name, driver.CarPlate, tripID}
	} else {
		query = `
            UPDATE trips 
            SET status = $1, updated_at = NOW()
            WHERE id = $2`
		args = []interface{}{status, tripID}
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update trip: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("trip not found: %s", tripID)
	}

	return nil
}

func (r *tripRepository) SaveRideFare(ctx context.Context, fare *domain.RideFareModel) error {
	query := `
        INSERT INTO ride_fares (user_id, package_slug, total_price_in_cents)
        VALUES ($1, $2, $3)
        RETURNING id`

	var idHex string
	err := r.db.QueryRowContext(ctx, query,
		fare.UserID,
		fare.PackageSlug,
		fare.TotalPriceInCents,
	).Scan(&idHex)

	if err != nil {
		return fmt.Errorf("failed to save ride fare: %w", err)
	}

	return nil
}

func (r *tripRepository) GetRideFareByID(ctx context.Context, id string) (*domain.RideFareModel, error) {
	query := `
        SELECT id, user_id, package_slug, total_price_in_cents
        FROM ride_fares 
        WHERE id = $1`

	var fare domain.RideFareModel
	var idHex string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&idHex,
		&fare.UserID,
		&fare.PackageSlug,
		&fare.TotalPriceInCents,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("ride fare not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get ride fare: %w", err)
	}

	return &fare, nil
}
