package util

import (
	"context"

	rocketpb "github.com/bachelor-thesis-hown3d/chat-api-server/proto/rocket/v1"
	tenantpb "github.com/bachelor-thesis-hown3d/chat-api-server/proto/tenant/v1"
	"github.com/bachelor-thesis-hown3d/chat-client/pkg/types"
	"google.golang.org/grpc"
)

func NewTenantClient(ctx context.Context, apiServer types.HttpURL) (tenantpb.TenantServiceClient, error) {
	conn, err := newConn(ctx, apiServer.Host)
	if err != nil {
		return nil, err
	}
	return tenantpb.NewTenantServiceClient(conn), nil
}

func NewRocketClient(ctx context.Context, apiServer types.HttpURL) (rocketpb.RocketServiceClient, error) {
	conn, err := newConn(ctx, apiServer.Host)
	if err != nil {
		return nil, err
	}
	return rocketpb.NewRocketServiceClient(conn), nil

}

func newConn(ctx context.Context, host string) (*grpc.ClientConn, error) {
	return grpc.DialContext(ctx, host, grpc.WithInsecure())
}
