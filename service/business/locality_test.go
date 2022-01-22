package business

import (
	"context"
	profileV1 "github.com/antinvestor/service-profile-api"
	propertyV1 "github.com/antinvestor/service-property-api"
	"github.com/pitabwire/frame"
	"testing"
)

func Test_localityBusiness_AddLocality(t *testing.T) {

	ctx := context.Background()
	profileCli := getProfileCli(t)

	type fields struct {
		service    *frame.Service
		profileCli *profileV1.ProfileClient
	}
	type args struct {
		ctx     context.Context
		message *propertyV1.Locality
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *propertyV1.Locality
		wantErr bool
	}{
		{
			name: "AddLocalitySuccessBoundary",
			fields: fields{
				service:    getService(ctx, "AddLocalitySuccessBoundary"),
				profileCli: profileCli,
			},
			args: args{
				ctx: ctx,
				message: &propertyV1.Locality{
					Feature: &propertyV1.Locality_Boundary{Boundary: "{\"type\": \"Polygon\",\"coordinates\": [[[2442627.9025405287, -3705499.954308534],[2425506.008204649,-3886502.837287831],[2425506.008204649,-3886502.837287831],[2555143.2081763083,-3910962.686339088],[2442627.9025405287,-3705499.954308534]]]}"},
					Name:    "TestLocalityBoundary",
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "AddLocalitySuccessPoint",
			fields: fields{
				service:    getService(ctx, "AddLocalitySuccessPoint"),
				profileCli: profileCli,
			},
			args: args{
				ctx: ctx,
				message: &propertyV1.Locality{
					Feature: &propertyV1.Locality_Point{Point: "{\"type\": \"Point\",\"coordinates\": [2442627.9025405287, -3705499.954308534]}"},
					Name:    "TestLocalityPoint",
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "AddLocalityFailValidation",
			fields: fields{
				service:    getService(ctx, "AddLocalityFailValidation"),
				profileCli: profileCli,
			},
			args: args{
				ctx: ctx,
				message: &propertyV1.Locality{
					Name: "TestLocality",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "AddLocalityFailedPointSupplied",
			fields: fields{
				service:    getService(ctx, "AddLocalityFailedPointSupplied"),
				profileCli: profileCli,
			},
			args: args{
				ctx: ctx,
				message: &propertyV1.Locality{
					Feature: &propertyV1.Locality_Point{Point: "{\"type\": \"Point\",\"coordinates\": [2442627.9025405287]}"},
					Name:    "TestLocalityPoint",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "AddLocalityFailedBoundarySupplied",
			fields: fields{
				service:    getService(ctx, "AddLocalityFailedBoundarySupplied"),
				profileCli: profileCli,
			},
			args: args{
				ctx: ctx,
				message: &propertyV1.Locality{
					Feature: &propertyV1.Locality_Boundary{Boundary: "{\"type\": \"Polygon\",\"coordinates\": [[[2442627.9025405287, -3705499.954308534, 43],[2425506.008204649,-3886502.837287831],[2425506.008204649,-3886502.837287831],[2555143.2081763083,-3910962.686339088],[2442627.9025405287,-3705499.954308534]]]}"},
					Name:    "TestLocalityBoundary",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, _ := NewLocalityBusiness(ctx, tt.fields.service, tt.fields.profileCli)

			got, err := l.AddLocality(tt.args.ctx, tt.args.message)

			if (err != nil) == tt.wantErr {
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("AddLocality() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got.GetName() != tt.args.message.GetName() {
				t.Errorf("AddLocality() got = %v, want %v", got.GetName(), tt.args.message.GetName())
				return
			}

			switch got.GetFeature().(type) {
			case *propertyV1.Locality_Point:
				if got.GetPoint() != tt.args.message.GetPoint() {
					t.Errorf("AddLocality() - Point got = %v, want %v", got.GetPoint(), tt.args.message.GetPoint())
				}

			case *propertyV1.Locality_Boundary:

				if got.GetBoundary() != tt.args.message.GetBoundary() {
					t.Errorf("AddLocality() - Boundary got = %v, want %v", got.GetBoundary(), tt.args.message.GetBoundary())
				}

			}

		})
	}
}

func Test_localityBusiness_DeleteLocality(t *testing.T) {

	ctx := context.Background()
	profileCli := getProfileCli(t)

	testPropertyID, err := getTestPropertyID(ctx, profileCli)
	if err != nil {
		t.Errorf("AddSubscription() we couldn't create a new property for test : %v", err)
	}

	lb := &localityBusiness{
		service:    getService(ctx, "LBTest"),
		profileCli: profileCli,
	}

	locality, err := lb.AddLocality(ctx, &propertyV1.Locality{
		ParentID: testPropertyID,
		//Feature:  &propertyV1.Locality_Boundary{Boundary: "{\"type\": \"Polygon\",\"coordinates\": [[[2442627.9025405287, -3705499.954308534],[2425506.008204649,-3886502.837287831],[2425506.008204649,-3886502.837287831],[2555143.2081763083,-3910962.686339088],[2442627.9025405287,-3705499.954308534]]]}"},
		Feature: &propertyV1.Locality_Point{Point: "{\"type\": \"Point\",\"coordinates\": [2442627.9025405287, -3705499.954308534]}"},
		Name:    "TestDeleteLocality",
	})
	if err != nil {
		t.Errorf("DeleteLocality() could not create delete locality because : %v \n\n\n ", err)
		return
	}

	type fields struct {
		service    *frame.Service
		profileCli *profileV1.ProfileClient
	}
	type args struct {
		ctx     context.Context
		message *propertyV1.RequestID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "DeleteLocalitySuccess",
			fields: fields{
				service:    getService(ctx, "DeleteLocalitySuccess"),
				profileCli: profileCli,
			},
			args: args{
				ctx: ctx,
				message: &propertyV1.RequestID{
					ID: locality.GetID(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &localityBusiness{
				service:    tt.fields.service,
				profileCli: tt.fields.profileCli,
			}
			if err := l.DeleteLocality(tt.args.ctx, tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("DeleteLocality() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
