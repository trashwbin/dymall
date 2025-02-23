package auth

import (
	"context"
	auth "github.com/trashwbin/dymall/rpc_gen/kitex_gen/auth"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/pkg/klog"
)

func DeliverTokenByRPC(ctx context.Context, req *auth.DeliverTokenReq, callOptions ...callopt.Option) (resp *auth.DeliveryResp, err error) {
	resp, err = defaultClient.DeliverTokenByRPC(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "DeliverTokenByRPC call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func VerifyTokenByRPC(ctx context.Context, req *auth.VerifyTokenReq, callOptions ...callopt.Option) (resp *auth.VerifyResp, err error) {
	resp, err = defaultClient.VerifyTokenByRPC(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "VerifyTokenByRPC call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func AddPolicy(ctx context.Context, req *auth.PolicyReq, callOptions ...callopt.Option) (resp *auth.PolicyResp, err error) {
	resp, err = defaultClient.AddPolicy(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "AddPolicy call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func RemovePolicy(ctx context.Context, req *auth.PolicyReq, callOptions ...callopt.Option) (resp *auth.PolicyResp, err error) {
	resp, err = defaultClient.RemovePolicy(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "RemovePolicy call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func AddRoleForUser(ctx context.Context, req *auth.RoleBindingReq, callOptions ...callopt.Option) (resp *auth.PolicyResp, err error) {
	resp, err = defaultClient.AddRoleForUser(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "AddRoleForUser call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func RemoveRoleForUser(ctx context.Context, req *auth.RoleBindingReq, callOptions ...callopt.Option) (resp *auth.PolicyResp, err error) {
	resp, err = defaultClient.RemoveRoleForUser(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "RemoveRoleForUser call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func GetRolesForUser(ctx context.Context, req *auth.RoleQueryReq, callOptions ...callopt.Option) (resp *auth.RoleQueryResp, err error) {
	resp, err = defaultClient.GetRolesForUser(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "GetRolesForUser call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}
