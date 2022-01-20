package business

import (
	"context"
	profileV1 "github.com/antinvestor/service-profile-api"
	propertyV1 "github.com/antinvestor/service-property-api"
	"github.com/pitabwire/frame"
	"reflect"
	"testing"
)

func Test_subscriptionBusiness_AddSubscription(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &subscriptionBusiness{
				service:    tt.fields.service,
				profileCli: tt.fields.profileCli,
			}
			got, err := s.AddSubscription(tt.args.ctx, tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddSubscription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddSubscription() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_subscriptionBusiness_DeleteSubscription(t *testing.T) {
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
		// TODO: Add test cases.
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

func Test_subscriptionBusiness_ListSubscriptions(t *testing.T) {
	type fields struct {
		service    *frame.Service
		profileCli *profileV1.ProfileClient
	}
	type args struct {
		message *propertyV1.SubscriptionListRequest
		stream  propertyV1.PropertyService_ListSubscriptionsServer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &subscriptionBusiness{
				service:    tt.fields.service,
				profileCli: tt.fields.profileCli,
			}
			if err := s.ListSubscriptions(tt.args.message, tt.args.stream); (err != nil) != tt.wantErr {
				t.Errorf("ListSubscriptions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
