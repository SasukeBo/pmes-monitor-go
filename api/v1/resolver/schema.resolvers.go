package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/SasukeBo/pmes-device-monitor/api/v1/generated"
	"github.com/SasukeBo/pmes-device-monitor/api/v1/logic"
	"github.com/SasukeBo/pmes-device-monitor/api/v1/model"
)

func (r *dashboardResolver) Devices(ctx context.Context, obj *model.Dashboard) ([]*model.Device, error) {
	return logic.LoadDevices(ctx, obj.DeviceIDs), nil
}

func (r *deviceResolver) DeviceType(ctx context.Context, obj *model.Device) (*model.DeviceType, error) {
	return logic.LoadDeviceType(ctx, obj.DeviceTypeID), nil
}

func (r *deviceTypeResolver) ErrorCode(ctx context.Context, obj *model.DeviceType) (*model.ErrorCode, error) {
	return logic.LoadErrorCode(ctx, obj.ErrorCodeID), nil
}

// Dashboard returns generated.DashboardResolver implementation.
func (r *Resolver) Dashboard() generated.DashboardResolver { return &dashboardResolver{r} }

// Device returns generated.DeviceResolver implementation.
func (r *Resolver) Device() generated.DeviceResolver { return &deviceResolver{r} }

// DeviceType returns generated.DeviceTypeResolver implementation.
func (r *Resolver) DeviceType() generated.DeviceTypeResolver { return &deviceTypeResolver{r} }

type dashboardResolver struct{ *Resolver }
type deviceResolver struct{ *Resolver }
type deviceTypeResolver struct{ *Resolver }
