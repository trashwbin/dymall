package auth

import (
	"context"
	auth "github.com/trashwbin/dymall/rpc_gen/kitex_gen/auth"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
	
)

type RPCClient interface {
	KitexClient() authservice.Client
	Service() string
	DeliverTokenByRPC(ctx context.Context, Req *auth.DeliverTokenReq, callOptions ...callopt.Option) (r *auth.DeliveryResp, err error)
	VerifyTokenByRPC(ctx context.Context, Req *auth.VerifyTokenReq, callOptions ...callopt.Option) (r *auth.VerifyResp, err error)
	AddPolicy(ctx context.Context, Req *auth.PolicyReq, callOptions ...callopt.Option) (r *auth.PolicyResp, err error)
	RemovePolicy(ctx context.Context, Req *auth.PolicyReq, callOptions ...callopt.Option) (r *auth.PolicyResp, err error)
	AddRoleForUser(ctx context.Context, Req *auth.RoleBindingReq, callOptions ...callopt.Option) (r *auth.PolicyResp, err error)
	RemoveRoleForUser(ctx context.Context, Req *auth.RoleBindingReq, callOptions ...callopt.Option) (r *auth.PolicyResp, err error)
	GetRolesForUser(ctx context.Context, Req *auth.RoleQueryReq, callOptions ...callopt.Option) (r *auth.RoleQueryResp, err error)
}

func NewRPCClient(dstService string, opts ...client.Option) (RPCClient, error) {
	kitexClient, err := authservice.NewClient(dstService, opts...)
	if err != nil {
		return nil, err
	}
	cli := &clientImpl{
		service:     dstService,
		kitexClient: kitexClient,
	}

	return cli, nil
}

type clientImpl struct {
	service     string
	kitexClient authservice.Client
}

func (c *clientImpl) Service() string {
	return c.service
}

func (c *clientImpl) KitexClient() authservice.Client {
	return c.kitexClient
}

func (c *clientImpl) DeliverTokenByRPC(ctx context.Context, Req *auth.DeliverTokenReq, callOptions ...callopt.Option) (r *auth.DeliveryResp, err error) {
	return c.kitexClient.DeliverTokenByRPC(ctx, Req, callOptions...)
}

func (c *clientImpl) VerifyTokenByRPC(ctx context.Context, Req *auth.VerifyTokenReq, callOptions ...callopt.Option) (r *auth.VerifyResp, err error) {
	return c.kitexClient.VerifyTokenByRPC(ctx, Req, callOptions...)
}

func (c *clientImpl) AddPolicy(ctx context.Context, Req *auth.PolicyReq, callOptions ...callopt.Option) (r *auth.PolicyResp, err error) {
	return c.kitexClient.AddPolicy(ctx, Req, callOptions...)
}

func (c *clientImpl) RemovePolicy(ctx context.Context, Req *auth.PolicyReq, callOptions ...callopt.Option) (r *auth.PolicyResp, err error) {
	return c.kitexClient.RemovePolicy(ctx, Req, callOptions...)
}

func (c *clientImpl) AddRoleForUser(ctx context.Context, Req *auth.RoleBindingReq, callOptions ...callopt.Option) (r *auth.PolicyResp, err error) {
	return c.kitexClient.AddRoleForUser(ctx, Req, callOptions...)
}

func (c *clientImpl) RemoveRoleForUser(ctx context.Context, Req *auth.RoleBindingReq, callOptions ...callopt.Option) (r *auth.PolicyResp, err error) {
	return c.kitexClient.RemoveRoleForUser(ctx, Req, callOptions...)
}

func (c *clientImpl) GetRolesForUser(ctx context.Context, Req *auth.RoleQueryReq, callOptions ...callopt.Option) (r *auth.RoleQueryResp, err error) {
	return c.kitexClient.GetRolesForUser(ctx, Req, callOptions...)
}
