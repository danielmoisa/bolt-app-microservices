package main

import (
	pb "github.com/danielmoisa/bolt-app/pkg/proto/trip"
	"github.com/danielmoisa/bolt-app/pkg/types"
)

// previewTripRequest represents a trip preview request
type previewTripRequest struct {
	UserID      string           `json:"userID" example:"user123" binding:"required"`
	Pickup      types.Coordinate `json:"pickup" binding:"required"`
	Destination types.Coordinate `json:"destination" binding:"required"`
}

func (p *previewTripRequest) toProto() *pb.PreviewTripRequest {
	return &pb.PreviewTripRequest{
		UserID: p.UserID,
		StartLocation: &pb.Coordinate{
			Latitude:  p.Pickup.Latitude,
			Longitude: p.Pickup.Longitude,
		},
		EndLocation: &pb.Coordinate{
			Latitude:  p.Destination.Latitude,
			Longitude: p.Destination.Longitude,
		},
	}
}

// startTripRequest represents a start trip request
type startTripRequest struct {
	RideFareID string `json:"rideFareID" example:"fare123" binding:"required"`
	UserID     string `json:"userID" example:"user123" binding:"required"`
}

func (c *startTripRequest) toProto() *pb.CreateTripRequest {
	return &pb.CreateTripRequest{
		RideFareID: c.RideFareID,
		UserID:     c.UserID,
	}
}
