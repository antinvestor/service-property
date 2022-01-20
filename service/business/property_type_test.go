package business

import (
	"context"
	profileV1 "github.com/antinvestor/service-profile-api"
	propertyV1 "github.com/antinvestor/service-property-api"
	"github.com/golang/mock/gomock"
	"github.com/pitabwire/frame"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
)

func Test_propertyTypeBusiness_AddPropertyType(t *testing.T) {

	ctx := context.Background()

	profileCli := getProfileCli(t)

	type fields struct {
		service    *frame.Service
		profileCli *profileV1.ProfileClient
	}
	type args struct {
		ctx     context.Context
		message *propertyV1.PropertyType
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *propertyV1.PropertyType
		wantErr bool
	}{
		{
			name: "AddPropertyTypeSuccess",
			fields: fields{
				service:    getService(ctx, "AddPropertyTypeSuccess"),
				profileCli: profileCli,
			},
			args: args{
				ctx: ctx,
				message: &propertyV1.PropertyType{
					Name:        "Just Building",
					Description: "A simple building that is under test multiply to 50",
					Extra:       map[string]string{"testing": "More test data"},
					CreatedAt:   timestamppb.Now(),
				},
			},
			want: &propertyV1.PropertyType{
				Name:        "Just Building",
				Description: "A simple building that is under test multiply to 50",
				Extra:       map[string]string{"testing": "More test data"},
				CreatedAt:   timestamppb.Now(),
			},
			wantErr: false,
		},
		{
			name: "AddPropertyTypeSuccessWithId",
			fields: fields{
				service:    getService(ctx, "AddPropertyTypeSuccessWithId"),
				profileCli: profileCli,
			},
			args: args{
				ctx: ctx,
				message: &propertyV1.PropertyType{
					ID:          "c2f4j7au6s7f91uqnoxg",
					Name:        "Checking with ID",
					Description: "A simple building that is under test multiply to 50",
					Extra:       map[string]string{"testing": "More test data"},
					CreatedAt:   timestamppb.Now(),
				},
			},
			want: &propertyV1.PropertyType{
				ID:          "c2f4j7au6s7f91uqnoxg",
				Name:        "Checking with ID",
				Description: "A simple building that is under test multiply to 50",
				Extra:       map[string]string{"testing": "More test data"},
				CreatedAt:   timestamppb.Now(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pt := &propertyTypeBusiness{
				service:    tt.fields.service,
				profileCli: tt.fields.profileCli,
			}
			got, err := pt.AddPropertyType(tt.args.ctx, tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddPropertyType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil && got.Name != tt.want.Name {
				t.Errorf("AddPropertyType() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_propertyTypeBusiness_ListPropertyType(t *testing.T) {

	ctx := context.Background()
	profileCli := getProfileCli(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	listTypeStream := propertyV1.NewMockPropertyService_ListTypeServer(ctrl)
	listTypeStream.EXPECT().Send(gomock.Any()).Return(nil).AnyTimes()
	listTypeStream.EXPECT().Context().Return(ctx).AnyTimes()

	type fields struct {
		service    *frame.Service
		profileCli *profileV1.ProfileClient
	}
	type args struct {
		message *propertyV1.SearchRequest
		stream  propertyV1.PropertyService_ListTypeServer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "ListPropertyTypeAll",
			fields: fields{
				service:    getService(ctx, "ListPropertyTypeAll"),
				profileCli: profileCli,
			},
			args: args{
				message: &propertyV1.SearchRequest{Query: ""},
				stream:  listTypeStream,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pt := &propertyTypeBusiness{
				service:    tt.fields.service,
				profileCli: tt.fields.profileCli,
			}
			if err := pt.ListPropertyType(tt.args.message, tt.args.stream); (err != nil) != tt.wantErr {
				t.Errorf("ListPropertyType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
