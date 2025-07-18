package domain

import (
	"context"
	"time"

	pbd "github.com/danielmoisa/bolt-app/pkg/proto/driver"
	pb "github.com/danielmoisa/bolt-app/pkg/proto/trip"
	"github.com/danielmoisa/bolt-app/pkg/types"
	tripTypes "github.com/danielmoisa/bolt-app/services/trip-service/pkg/types"
)

type TripModel struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    string
	Status    string
	RideFare  *RideFareModel
	Driver    *pb.TripDriver
}

func (t *TripModel) ToProto() *pb.Trip {
	return &pb.Trip{
		Id:           t.ID,
		UserID:       t.UserID,
		SelectedFare: t.RideFare.ToProto(),
		Status:       t.Status,
		Driver:       t.Driver,
		Route:        t.RideFare.Route.ToProto(),
	}
}

type TripRepository interface {
	CreateTrip(ctx context.Context, trip *TripModel) (*TripModel, error)
	SaveRideFare(ctx context.Context, f *RideFareModel) error
	GetRideFareByID(ctx context.Context, id string) (*RideFareModel, error)
	GetTripByID(ctx context.Context, id string) (*TripModel, error)
	UpdateTrip(ctx context.Context, tripID string, status string, driver *pbd.Driver) error
}

type TripService interface {
	CreateTrip(ctx context.Context, fare *RideFareModel) (*TripModel, error)
	GetRoute(ctx context.Context, pickup, destination *types.Coordinate, useOsrmApi bool) (*tripTypes.OsrmApiResponse, error)
	EstimatePackagesPriceWithRoute(route *tripTypes.OsrmApiResponse) []*RideFareModel
	GenerateTripFares(
		ctx context.Context,
		fares []*RideFareModel,
		userID string,
		Route *tripTypes.OsrmApiResponse,
	) ([]*RideFareModel, error)
	GetAndValidateFare(ctx context.Context, fareID, userID string) (*RideFareModel, error)
	GetTripByID(ctx context.Context, id string) (*TripModel, error)
	UpdateTrip(ctx context.Context, tripID string, status string, driver *pbd.Driver) error
}
