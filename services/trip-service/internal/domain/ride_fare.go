package domain

import (
	pb "github.com/danielmoisa/bolt-app/pkg/proto/trip"
	"github.com/danielmoisa/bolt-app/services/trip-service/pkg/types"
)

type RideFareModel struct {
	ID                string                 `json:"id"`
	UserID            string                 `json:"userID"`
	PackageSlug       string                 `json:"packageSlug"` // ex: van, luxury, sedan
	TotalPriceInCents float64                `json:"totalPriceInCents"`
	Route             *types.OsrmApiResponse `json:"route"`
}

func (r *RideFareModel) ToProto() *pb.RideFare {
	return &pb.RideFare{
		Id:                r.ID,
		UserID:            r.UserID,
		PackageSlug:       r.PackageSlug,
		TotalPriceInCents: r.TotalPriceInCents,
	}
}

func ToRideFaresProto(fares []*RideFareModel) []*pb.RideFare {
	var protoFares []*pb.RideFare
	for _, f := range fares {
		protoFares = append(protoFares, f.ToProto())
	}
	return protoFares
}
