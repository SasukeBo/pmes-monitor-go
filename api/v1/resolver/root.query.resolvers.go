package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/SasukeBo/pmes-device-monitor/api/v1/logic"

	"github.com/SasukeBo/pmes-device-monitor/api/v1/generated"
	"github.com/SasukeBo/pmes-device-monitor/api/v1/model"
)

func (r *queryResolver) Dashboards(ctx context.Context, search *string, limit int, page int) (*model.DashboardWrap, error) {
	return logic.Dashboards(ctx, search, limit, page)
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
