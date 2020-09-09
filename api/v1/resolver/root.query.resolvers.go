package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/SasukeBo/pmes-device-monitor/api/v1/generated"
	"github.com/SasukeBo/pmes-device-monitor/api/v1/logic"
	"github.com/SasukeBo/pmes-device-monitor/api/v1/model"
)

func (r *queryResolver) Dashboards(ctx context.Context, search *string, limit int, page int) (*model.DashboardWrap, error) {
	return logic.Dashboards(ctx, search, limit, page)
}

func (r *queryResolver) DashboardDevices(ctx context.Context, id int) ([]*model.DashboardDevice, error) {
	return logic.DashboardDevices(ctx, id)
}

func (r *queryResolver) DashboardDeviceFresh(ctx context.Context, id int, pid int, sid int) (*model.DashboardDeviceFreshResponse, error) {
	return logic.DashboardDeviceFresh(ctx, id, pid, sid)
}

func (r *queryResolver) DashboardOverviewAnalyze(ctx context.Context, id int) (*model.DashboardOverviewAnalyzeResponse, error) {
	return logic.DashboardOverviewAnalyze(ctx, id)
}

func (r *queryResolver) DashboardDeviceStatus(ctx context.Context, id int) (*model.DashboardDeviceStatusResponse, error) {
	return logic.DashboardDeviceStatus(ctx, id)
}

func (r *queryResolver) DashboardDeviceErrors(ctx context.Context, id int) (*model.DashboardDeviceErrorsResponse, error) {
	return logic.DashboardDeviceErrors(ctx, id)
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
