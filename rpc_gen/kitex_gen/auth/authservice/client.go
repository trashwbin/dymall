// Code generated by Kitex v0.9.1. DO NOT EDIT.

package authservice

import (
	"context"
	client "github.com/cloudwego/kitex/client"
	callopt "github.com/cloudwego/kitex/client/callopt"
	auth "github.com/trashwbin/dymall/rpc_gen/kitex_gen/auth"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	DeliverTokenByRPC(ctx context.Context, Req *auth.DeliverTokenReq, callOptions ...callopt.Option) (r *auth.DeliveryResp, err error)
	VerifyTokenByRPC(ctx context.Context, Req *auth.VerifyTokenReq, callOptions ...callopt.Option) (r *auth.VerifyResp, err error)
	AddPolicy(ctx context.Context, Req *auth.PolicyReq, callOptions ...callopt.Option) (r *auth.PolicyResp, err error)
	RemovePolicy(ctx context.Context, Req *auth.PolicyReq, callOptions ...callopt.Option) (r *auth.PolicyResp, err error)
	AddRoleForUser(ctx context.Context, Req *auth.RoleBindingReq, callOptions ...callopt.Option) (r *auth.PolicyResp, err error)
	RemoveRoleForUser(ctx context.Context, Req *auth.RoleBindingReq, callOptions ...callopt.Option) (r *auth.PolicyResp, err error)
	GetRolesForUser(ctx context.Context, Req *auth.RoleQueryReq, callOptions ...callopt.Option) (r *auth.RoleQueryResp, err error)
}

// NewClient creates a client for the service defined in IDL.
func NewClient(destService string, opts ...client.Option) (Client, error) {
	var options []client.Option
	options = append(options, client.WithDestService(destService))

	options = append(options, opts...)

	kc, err := client.NewClient(serviceInfo(), options...)
	if err != nil {
		return nil, err
	}
	return &kAuthServiceClient{
		kClient: newServiceClient(kc),
	}, nil
}

// MustNewClient creates a client for the service defined in IDL. It panics if any error occurs.
func MustNewClient(destService string, opts ...client.Option) Client {
	kc, err := NewClient(destService, opts...)
	if err != nil {
		panic(err)
	}
	return kc
}

type kAuthServiceClient struct {
	*kClient
}

func (p *kAuthServiceClient) DeliverTokenByRPC(ctx context.Context, Req *auth.DeliverTokenReq, callOptions ...callopt.Option) (r *auth.DeliveryResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.DeliverTokenByRPC(ctx, Req)
}

func (p *kAuthServiceClient) VerifyTokenByRPC(ctx context.Context, Req *auth.VerifyTokenReq, callOptions ...callopt.Option) (r *auth.VerifyResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.VerifyTokenByRPC(ctx, Req)
}

func (p *kAuthServiceClient) AddPolicy(ctx context.Context, Req *auth.PolicyReq, callOptions ...callopt.Option) (r *auth.PolicyResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.AddPolicy(ctx, Req)
}

func (p *kAuthServiceClient) RemovePolicy(ctx context.Context, Req *auth.PolicyReq, callOptions ...callopt.Option) (r *auth.PolicyResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.RemovePolicy(ctx, Req)
}

func (p *kAuthServiceClient) AddRoleForUser(ctx context.Context, Req *auth.RoleBindingReq, callOptions ...callopt.Option) (r *auth.PolicyResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.AddRoleForUser(ctx, Req)
}

func (p *kAuthServiceClient) RemoveRoleForUser(ctx context.Context, Req *auth.RoleBindingReq, callOptions ...callopt.Option) (r *auth.PolicyResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.RemoveRoleForUser(ctx, Req)
}

func (p *kAuthServiceClient) GetRolesForUser(ctx context.Context, Req *auth.RoleQueryReq, callOptions ...callopt.Option) (r *auth.RoleQueryResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetRolesForUser(ctx, Req)
}
