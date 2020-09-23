package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/SasukeBo/pmes-device-monitor/api/v1/admin/generated"
	"github.com/SasukeBo/pmes-device-monitor/api/v1/admin/logic"
	"github.com/SasukeBo/pmes-device-monitor/api/v1/admin/model"
)

func (r *mutationResolver) ImportErrors(ctx context.Context, deviceID int, fileToken string) (string, error) {
	return logic.ImportErrors(ctx, deviceID, fileToken)
}

func (r *mutationResolver) AdminDeviceTypeCreate(ctx context.Context, name string) (string, error) {
	return logic.AdminDeviceTypeCreate(ctx, name)
}

func (r *mutationResolver) AdminDeviceTypeDelete(ctx context.Context, id int) (string, error) {
	return logic.AdminDeviceTypeDelete(ctx, id)
}

func (r *mutationResolver) AdminDeviceTypeAddErrorCode(ctx context.Context, deviceTypeID int, errors []string) (string, error) {
	return logic.AdminDeviceTypeAddErrorCode(ctx, deviceTypeID, errors)
}

func (r *mutationResolver) AdminSaveErrorCode(ctx context.Context, id int, errors []string) (string, error) {
	return logic.AdminSaveErrorCode(ctx, id, errors)
}

func (r *mutationResolver) AdminCreateDevices(ctx context.Context, input model.CreateDeviceInput) (string, error) {
	return logic.AdminCreateDevices(ctx, input)
}

func (r *mutationResolver) AdminDeleteDevice(ctx context.Context, id int) (string, error) {
	return logic.AdminDeleteDevice(ctx, id)
}

func (r *mutationResolver) AdminEditDevice(ctx context.Context, id int, update model.DeviceInput) (string, error) {
	return logic.AdminEditDevice(ctx, id, update)
}

func (r *mutationResolver) AdminCreateDashboard(ctx context.Context, name string, deviceIDs []int) (string, error) {
	return logic.AdminCreateDashboard(ctx, name, deviceIDs)
}

func (r *mutationResolver) AdminDashboardDelete(ctx context.Context, id int) (string, error) {
	return logic.AdminDashboardDelete(ctx, id)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
