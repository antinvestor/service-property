package business

import (
	"context"
	profileV1 "github.com/antinvestor/service-profile-api"
	propertyV1 "github.com/antinvestor/service-property-api"
	"github.com/pitabwire/frame"
	"reflect"
	"testing"
)

func Test_localityBusiness_AddLocality(t *testing.T) {

	ctx := context.Background()
	profileCli := getProfileCli(t)

	testPropertyID, err := getTestPropertyID(ctx, profileCli)
	if err != nil {
		t.Errorf("AddSubscription() we couldn't create a new property for test : %v", err)
	}

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
			name: "AddLocalitySuccess",
			fields: fields{
				service:    getService(ctx, "AddLocalitySuccess"),
				profileCli: profileCli,
			},
			args: args{
				ctx: ctx,
				message: &propertyV1.Locality{
					ParentID: testPropertyID,
					Boundary: "[]",
					Name:     "TestLocality",
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &localityBusiness{
				service:    tt.fields.service,
				profileCli: tt.fields.profileCli,
			}
			got, err := l.AddLocality(tt.args.ctx, tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddLocality() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddLocality() got = %v, want %v", got, tt.want)
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
		Boundary: "",
		Name:     "TestDeleteLocality",
	})
	if err != nil {
		t.Errorf("DeleteLocality() could not create delete locality because : %v", err)
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
