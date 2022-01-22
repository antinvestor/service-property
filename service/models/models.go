package models

import (
	"github.com/antinvestor/apis/common"
	propertyV1 "github.com/antinvestor/service-property-api"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"

	"github.com/pitabwire/frame"
	"gorm.io/datatypes"
)

// Locality Table holds the location data relating to our properties
type Locality struct {
	frame.BaseModel

	ParentID    string `gorm:"type:varchar(50)"`
	Name        string `gorm:"type:varchar(50)"`
	Description string `gorm:"type:text"`
	Point       datatypes.JSON
	Boundary    datatypes.JSON
	Extra       datatypes.JSONMap
}

func (l *Locality) ToApi() *propertyV1.Locality {

	locality := &propertyV1.Locality{
		ID:          l.GetID(),
		ParentID:    l.ParentID,
		Name:        l.Name,
		Description: l.Description,
		Extras:      frame.DBPropertiesToMap(l.Extra),
		CreatedAt:   timestamppb.New(l.CreatedAt),
	}

	if l.Boundary.String() != "{}" {
		locality.Feature = &propertyV1.Locality_Boundary{Boundary: l.Boundary.String()}
	} else {

		locality.Feature = &propertyV1.Locality_Point{Point: l.Point.String()}
	}

	return locality
}

type PropertyType struct {
	frame.BaseModel

	Name        string `gorm:"type:varchar(250)"`
	Description string `gorm:"type:text"`
	Extra       datatypes.JSONMap
}

func (pt *PropertyType) ToApi() *propertyV1.PropertyType {
	return &propertyV1.PropertyType{
		ID:          pt.GetID(),
		Name:        pt.Name,
		Description: pt.Description,
		Extra:       frame.DBPropertiesToMap(pt.Extra),
		CreatedAt:   timestamppb.New(pt.CreatedAt),
	}
}

type Subscription struct {
	frame.BaseModel

	PropertyID  string `gorm:"type:varchar(50)"`
	ProfileID   string `gorm:"type:varchar(50)"`
	Role        string `gorm:"type:varchar(250)"`
	Description string `gorm:"type:text"`
	Extra       datatypes.JSONMap
	ExpiresAt   time.Time
}

func (s *Subscription) ToApi() *propertyV1.Subscription {

	return &propertyV1.Subscription{
		ID:         s.GetID(),
		ProfileID:  s.ProfileID,
		PropertyID: s.PropertyID,
		Role:       s.Role,
		Extra:      frame.DBPropertiesToMap(s.Extra),
		ExpiresAt:  timestamppb.New(s.ExpiresAt),
		CreatedAt:  timestamppb.New(s.CreatedAt),
	}
}

type Property struct {
	frame.BaseModel

	ParentID string `gorm:"type:varchar(50)"`

	PropertyTypeID string `gorm:"type:varchar(50)"`
	LocalityID     string `gorm:"type:varchar(50)"`

	Name        string `gorm:"type:varchar(250)"`
	Description string `gorm:"type:text"`
	Extra       datatypes.JSONMap

	StartedAt time.Time
	StateID   string `gorm:"type:varchar(50)"`
}

type PropertyState struct {
	frame.BaseModel
	PropertyID string `gorm:"type:varchar(50)"`

	Name        string `gorm:"type:varchar(250)"`
	Description string `gorm:"type:text"`
	Extra       datatypes.JSONMap
	State       int32
	Status      int32
}

func (ps *PropertyState) ToApi() *propertyV1.PropertyState {
	return &propertyV1.PropertyState{
		ID:          ps.GetID(),
		PropertyID:  ps.PropertyID,
		Name:        ps.Name,
		Description: ps.Description,
		Extras:      frame.DBPropertiesToMap(ps.Extra),
		Status:      common.STATUS(ps.Status),
		State:       common.STATE(ps.State),
		CreatedAt:   timestamppb.New(ps.CreatedAt),
	}
}
