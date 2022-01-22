package business

import (
	"context"
	profileV1 "github.com/antinvestor/service-profile-api"
	propertyV1 "github.com/antinvestor/service-property-api"
	"github.com/golang/mock/gomock"
	"github.com/pitabwire/frame"
	"reflect"
	"testing"
)

func Test_subscriptionBusiness_AddSubscription(t *testing.T) {

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
		message *propertyV1.Subscription
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *propertyV1.Subscription
		wantErr bool
	}{
		{
			name: "AddSubscriptionSuccess",
			fields: fields{
				service:    getService(ctx, "AddSubscriptionSuccess"),
				profileCli: profileCli,
			},
			args: args{
				ctx: ctx,
				message: &propertyV1.Subscription{
					PropertyID: testPropertyID,
					ProfileID:  "test_profile",
					Role:       "tester",
				},
			},
			want: &propertyV1.Subscription{
				PropertyID: testPropertyID,
				ProfileID:  "test_profile",
				Role:       "tester",
			},
			wantErr: false,
		},

		{
			name: "AddSubscriptionFailure",
			fields: fields{
				service:    getService(ctx, "AddSubscriptionFailure"),
				profileCli: profileCli,
			},
			args: args{
				ctx: ctx,
				message: &propertyV1.Subscription{
					PropertyID: "randomMissingProperty",
					ProfileID:  "test_profile",
					Role:       "tester",
				},
			},
			want: &propertyV1.Subscription{
				PropertyID: testPropertyID,
				ProfileID:  "test_profile",
				Role:       "tester",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			s, _ := NewSubscriptionBusiness(ctx, tt.fields.service, tt.fields.profileCli)
			got, err := s.AddSubscription(tt.args.ctx, tt.args.message)

			if tt.wantErr == (err != nil) {
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("AddSubscription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Role, tt.want.Role) {
				t.Errorf("AddSubscription() got = %v, want %v", got.Role, tt.want.Role)
			}

			if !reflect.DeepEqual(got.PropertyID, tt.want.PropertyID) {
				t.Errorf("AddSubscription() got = %v, want %v", got.PropertyID, tt.want.PropertyID)
			}
		})
	}
}

func Test_subscriptionBusiness_DeleteSubscription(t *testing.T) {
	ctx := context.Background()
	profileCli := getProfileCli(t)

	testPropertyID, err := getTestPropertyID(ctx, profileCli)
	if err != nil {
		t.Errorf("DeleteSubscription() we couldn't create a new property for test : %v", err)
	}

	testSubBiz := &subscriptionBusiness{
		service:    getService(ctx, "DeleteSubscription"),
		profileCli: profileCli,
	}
	subscription, err := testSubBiz.AddSubscription(ctx, &propertyV1.Subscription{
		PropertyID: testPropertyID,
		ProfileID:  "test_profile_delete",
		Role:       "Testing",
	})
	if err != nil {
		t.Errorf("DeleteSubscription() we couldn't create a new property for test : %v", err)
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
			name: "DeleteSubscriptionSuccess",
			fields: fields{
				service:    getService(ctx, "DeleteSubscriptionSuccess"),
				profileCli: profileCli,
			},
			args: args{
				ctx: ctx,
				message: &propertyV1.RequestID{
					ID: subscription.GetID(),
				},
			},
			wantErr: false,
		},
		{
			name: "DeleteSubscriptionFail",
			fields: fields{
				service:    getService(ctx, "DeleteSubscriptionFail"),
				profileCli: profileCli,
			},
			args: args{
				ctx: ctx,
				message: &propertyV1.RequestID{
					ID: "some_random_id",
				},
			},
			wantErr: true,
		},

		{
			name: "DeleteSubscriptionFailValidation",
			fields: fields{
				service:    getService(ctx, "DeleteSubscriptionFailValidation"),
				profileCli: profileCli,
			},
			args: args{
				ctx: ctx,
				message: &propertyV1.RequestID{
					ID: "",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &subscriptionBusiness{
				service:    tt.fields.service,
				profileCli: tt.fields.profileCli,
			}
			if err := s.DeleteSubscription(tt.args.ctx, tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("DeleteSubscription() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func getListSubscriptionStream(ctx context.Context, controller *gomock.Controller, expectedSendCount int) *propertyV1.MockPropertyService_ListSubscriptionsServer {
	listSubscriptionStream := propertyV1.NewMockPropertyService_ListSubscriptionsServer(controller)
	listSubscriptionStream.EXPECT().Context().Return(ctx).AnyTimes()

	listSubscriptionStream.EXPECT().Send(gomock.Any()).Return(nil).Times(expectedSendCount)

	return listSubscriptionStream
}

func Test_subscriptionBusiness_ListSubscriptions(t *testing.T) {
	ctx := context.Background()
	profileCli := getProfileCli(t)

	testPropertyID, err := getTestPropertyID(ctx, profileCli)
	if err != nil {
		t.Errorf("ListSubscriptions() we couldn't create a new property for test : %v", err)
	}

	testSubBiz := &subscriptionBusiness{
		service:    getService(ctx, "SubscriptionTest"),
		profileCli: profileCli,
	}

	for _, subsc := range []*propertyV1.Subscription{{
		PropertyID: testPropertyID,
		ProfileID:  "test_profile_list_1",
		Role:       "Testing_list_1",
	}, {
		PropertyID: testPropertyID,
		ProfileID:  "test_profile_list_2",
		Role:       "Testing_list_2",
	}, {
		PropertyID: testPropertyID,
		ProfileID:  "test_profile_list_3",
		Role:       "Testing_list_3",
	}} {

		_, err := testSubBiz.AddSubscription(ctx, subsc)
		if err != nil {
			t.Errorf("ListSubscriptions() we couldn't create a new property for test : %v", err)
		}
	}
	type fields struct {
		service    *frame.Service
		profileCli *profileV1.ProfileClient
	}

	type args struct {
		message           *propertyV1.SubscriptionListRequest
		expectedSendCount int
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		resultCount int
		wantErr     bool
	}{
		{
			name: "ListSubscriptionsSuccess",
			fields: fields{
				service:    getService(ctx, "ListSubscriptionsSuccess"),
				profileCli: profileCli,
			},
			args: args{
				message: &propertyV1.SubscriptionListRequest{
					PropertyID: testPropertyID},
				expectedSendCount: 3,
			},
			wantErr: false,
		},
		{
			name: "ListSubscriptionsSuccessQuery",
			fields: fields{
				service:    getService(ctx, "ListSubscriptionsSuccessQuery"),
				profileCli: profileCli,
			},
			args: args{
				message: &propertyV1.SubscriptionListRequest{
					PropertyID: testPropertyID, Query: "_list_2"},
				expectedSendCount: 1,
			},
			wantErr: false,
		},
		{
			name: "ListSubscriptionsFailure",
			fields: fields{
				service:    getService(ctx, "ListSubscriptionsFailure"),
				profileCli: profileCli,
			},
			args: args{
				message: &propertyV1.SubscriptionListRequest{
					PropertyID: "randomNonId"},
				expectedSendCount: 0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			controller := gomock.NewController(t)

			s := &subscriptionBusiness{
				service:    tt.fields.service,
				profileCli: tt.fields.profileCli,
			}
			if err := s.ListSubscriptions(tt.args.message, getListSubscriptionStream(ctx, controller, tt.args.expectedSendCount)); (err != nil) != tt.wantErr {
				t.Errorf("ListSubscriptions() error = %v, wantErr %v", err, tt.wantErr)
			}

			controller.Finish()

		})
	}
}
