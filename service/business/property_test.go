package business

import (
	"context"
	"github.com/antinvestor/apis/common"
	partitionV1 "github.com/antinvestor/service-partition-api"
	profileV1 "github.com/antinvestor/service-profile-api"
	propertyV1 "github.com/antinvestor/service-property-api"
	"github.com/antinvestor/service-property/service/events"
	"github.com/golang/mock/gomock"
	"github.com/pitabwire/frame"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
)

func getService(ctx context.Context, serviceName string) *frame.Service {
	dbUrl := frame.GetEnv("TEST_DATABASE_URL", "postgres://ant:secret@localhost:5437/service_property?sslmode=disable")
	testDb := frame.Datastore(ctx, dbUrl, false)

	service := frame.NewService(serviceName, testDb, frame.NoopHttpOptions())

	eventList := frame.RegisterEvents(&events.PropertyStateSave{Service: service})
	service.Init(eventList)
	_ = service.Run(ctx, "")
	return service
}

func getProfileCli(t *testing.T) *profileV1.ProfileClient {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockProfileService := profileV1.NewMockProfileServiceClient(ctrl)
	profileCli := profileV1.InstantiateProfileClient(nil, mockProfileService)
	return profileCli
}

func getPartitionCli(t *testing.T) *partitionV1.PartitionClient {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPartitionService := partitionV1.NewMockPartitionServiceClient(ctrl)

	mockPartitionService.EXPECT().
		GetAccess(gomock.Any(), gomock.Any()).
		Return(&partitionV1.AccessObject{
			AccessId: "test_access-id",
			Partition: &partitionV1.PartitionObject{
				PartitionId: "test_partition-id",
				TenantId:    "test_tenant-id",
			},
		}, nil).AnyTimes()

	profileCli := partitionV1.InstantiatePartitionsClient(nil, mockPartitionService)
	return profileCli
}

func TestNewPartitionBusiness(t *testing.T) {

	ctx := context.Background()

	profileCli := getProfileCli(t)
	partitionCli := getPartitionCli(t)

	type args struct {
		ctx          context.Context
		service      *frame.Service
		profileCli   *profileV1.ProfileClient
		partitionCli *partitionV1.PartitionClient
	}
	tests := []struct {
		name      string
		args      args
		want      PropertyBusiness
		expectErr bool
	}{

		{name: "NewPropertyBusiness",
			args: args{
				ctx:          ctx,
				service:      getService(ctx, "NewPropertyBusinessTest"),
				profileCli:   profileCli,
				partitionCli: partitionCli},
			expectErr: false},

		{name: "NewPropertyBusinessWithNils",
			args: args{
				ctx:        ctx,
				service:    nil,
				profileCli: nil,
			},
			expectErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := NewPropertyBusiness(tt.args.ctx, tt.args.service, tt.args.profileCli, tt.args.partitionCli); !tt.expectErr && (err != nil || got == nil) {
				t.Errorf("NewPropertyBusiness() = could not get a valid propertyBusiness at %s", tt.name)
			}
		})
	}
}

func Test_propertyBusiness_CreateProperty(t *testing.T) {

	ctx := context.Background()
	profileCli := getProfileCli(t)

	pt := &propertyTypeBusiness{
		service:    getService(ctx, "profileTypeService"),
		profileCli: profileCli,
	}
	propertyType, err := pt.AddPropertyType(ctx, &propertyV1.PropertyType{
		Name: "Residential",
	})
	if err != nil {
		t.Errorf("CreateProperty() we couldn't create a new property %v", err)
	}

	type fields struct {
		service    *frame.Service
		profileCli *profileV1.ProfileClient
	}
	type args struct {
		ctx     context.Context
		message *propertyV1.Property
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *propertyV1.PropertyState
		wantErr bool
	}{
		{name: "NormalPassingCreateProperty",
			fields: fields{
				service:    getService(ctx, "NormalCreatePropertyTest"),
				profileCli: profileCli,
			},
			args: args{
				ctx: ctx,
				message: &propertyV1.Property{
					ID:           "justtestingId",
					Name:         "epochTesting",
					Description:  "Hello we are just testing things out Multiply so we over 50",
					StartedAt:    timestamppb.Now(),
					PropertyType: propertyType,
				},
			},
			wantErr: false,
			want: &propertyV1.PropertyState{
				ID:     "123456",
				State:  common.STATE_CREATED,
				Status: common.STATUS_QUEUED,
			},
		},
		{name: "NormalPassingCreatePropertyWithID",
			fields: fields{
				service:    getService(ctx, "NormalPassingCreatePropertyWithID"),
				profileCli: getProfileCli(t),
			},
			args: args{
				ctx: ctx,
				message: &propertyV1.Property{
					ID:           "c7k7ad5498t587vbh1k0",
					PropertyType: propertyType,
					Name:         "epochTesting",
					Description:  "Hello we are just testing things out, increase line to over 50",
					StartedAt:    timestamppb.Now(),
				},
			},
			wantErr: false,
			want: &propertyV1.PropertyState{
				PropertyID: "c7k7ad5498t587vbh1k0",
				State:      common.STATE_CREATED,
				Status:     common.STATUS_QUEUED,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			propBusiness := &propertyBusiness{
				service:    tt.fields.service,
				profileCli: tt.fields.profileCli,
			}
			got, err := propBusiness.CreateProperty(tt.args.ctx, tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateProperty() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.GetStatus() != tt.want.GetStatus() || got.GetState() != tt.want.GetState() {
				t.Errorf("CreateProperty() got = %v, want %v", got, tt.want)
			}

			if tt.name == "NormalPassingCreatePropertyWithID" && got.GetPropertyID() != tt.want.GetPropertyID() {
				t.Errorf("CreateProperty() expecting id %s to be reused, got : %s", tt.want.GetID(), got.GetID())
			}
		})
	}
}


