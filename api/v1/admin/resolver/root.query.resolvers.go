package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/SasukeBo/pmes-device-monitor/api/v1/admin/generated"
	"github.com/SasukeBo/pmes-device-monitor/api/v1/admin/logic"
	"github.com/SasukeBo/pmes-device-monitor/api/v1/admin/model"
)

func (r *queryResolver) Hello(ctx context.Context, name string) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) AdminDeviceTypes(ctx context.Context, search *string, page int, limit int) (*model.DeviceTypeWrap, error) {
	return logic.AdminDeviceTypes(ctx, search, page, limit)
}

func (r *queryResolver) AdminDeviceType(ctx context.Context, id int) (*model.DeviceType, error) {
	return logic.AdminDeviceType(ctx, id)
}

func (r *queryResolver) AdminDevices(ctx context.Context, search *string, page int, limit int) (*model.DeviceWrap, error) {
	return logic.AdminDevices(ctx, search, page, limit)
}

func (r *queryResolver) AdminDashboards(ctx context.Context, search *string, page int, limit int) (*model.DashboardWrap, error) {
	return logic.AdminDashboards(ctx, search, page, limit)
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
