package business

import (
	"context"
	"errors"
	profileV1 "github.com/antinvestor/service-profile-api"
	propertyV1 "github.com/antinvestor/service-property-api"
	"github.com/antinvestor/service-property/service/models"
	"github.com/antinvestor/service-property/service/repository"
	"github.com/pitabwire/frame"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
)

type localityBusiness struct {
	service    *frame.Service
	profileCli *profileV1.ProfileClient
}

func (l *localityBusiness) AddLocality(ctx context.Context, message *propertyV1.Locality) (*propertyV1.Locality, error) {

	err := message.Validate()
	if err != nil {
		return nil, err
	}

	locality := models.Locality{
		Name:        message.GetName(),
		Description: message.GetDescription(),
		ParentID:    message.GetParentID(),
		Extra:       frame.DBPropertiesFromMap(message.Extras),
	}

	switch v := message.GetFeature().(type) {
	case *propertyV1.Locality_Boundary:

		var geoJsonFeature geom.T
		err = geojson.Unmarshal([]byte(v.Boundary), &geoJsonFeature)
		if err != nil {
			return nil, err
		}

		_, ok := geoJsonFeature.(*geom.Polygon)
		if !ok {
			return nil, errors.New("supplied geometry is not a polygon")
		}

		locality.Boundary = []byte(v.Boundary)
		locality.Point = []byte(`{}`)

	case *propertyV1.Locality_Point:
		var geoJsonFeature geom.T
		err = geojson.Unmarshal([]byte(v.Point), &geoJsonFeature)
		if err != nil {
			return nil, err
		}

		_, ok := geoJsonFeature.(*geom.Point)
		if !ok {
			return nil, errors.New("supplied geometry is not a point")
		}

		locality.Boundary = []byte(`{}`)
		locality.Point = []byte(v.Point)
	}

	if locality.ValidXID(message.GetID()) {
		locality.ID = message.GetID()
	} else {
		locality.GenID(ctx)
	}

	localityRepository := repository.NewLocalityRepository(ctx, l.service)
	err = localityRepository.Save(&locality)
	if err != nil {
		return nil, err
	}

	return locality.ToApi(), nil
}

func (l *localityBusiness) DeleteLocality(ctx context.Context, message *propertyV1.RequestID) error {

	err := message.Validate()
	if err != nil {
		return err
	}

	localityRepository := repository.NewLocalityRepository(ctx, l.service)

	var locality models.Locality
	err = localityRepository.GetByID(message.GetID(), &locality)
	if err != nil {
		return err
	}

	return localityRepository.Delete(locality.GetID())
}
