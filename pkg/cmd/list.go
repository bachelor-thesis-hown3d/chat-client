package cmd

import (
	"context"
	"fmt"

	"github.com/bachelor-thesis-hown3d/chat-api-server/proto/rocket/v1"
	"github.com/bachelor-thesis-hown3d/chat-client/pkg/oauth"
	"github.com/bachelor-thesis-hown3d/chat-client/pkg/types"
	"github.com/bachelor-thesis-hown3d/chat-client/pkg/util"
)

func List(ctx context.Context, apiServer types.HttpURL) error {
	ctx, err := oauth.LoadTokenIntoContext(ctx)
	if err != nil {
		return err
	}

	client, err := util.NewRocketClient(ctx, apiServer)
	req := &rocket.GetAllRequest{}
	resp, err := client.GetAll(ctx, req)
	if err != nil {
		return err
	}

	for _, rocket := range resp.Rockets {
		fmt.Println(rocket)
	}
	return nil
}
