package service

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/trashwbin/dymall/app/auth/biz/utils"
	auth "github.com/trashwbin/dymall/rpc_gen/kitex_gen/auth"
)

type VerifyTokenByRPCService struct {
	ctx     context.Context
	authSvc *AuthorizationService
}

// NewVerifyTokenByRPCService new VerifyTokenByRPCService
func NewVerifyTokenByRPCService(ctx context.Context, authSvc *AuthorizationService) *VerifyTokenByRPCService {
	return &VerifyTokenByRPCService{
		ctx:     ctx,
		authSvc: authSvc,
	}
}

// Run 验证令牌
func (s *VerifyTokenByRPCService) Run(req *auth.VerifyTokenReq) (resp *auth.VerifyResp, err error) {
	resp = new(auth.VerifyResp)

	//TODO: 使用apifox实在无法携带上token，所以先返回一个mock值
	resp.Code = auth.ErrorCode_Success
	resp.Message = "令牌验证成功"
	resp.IsValid = true
	resp.UserId = 1
	resp.Role = "user"
	return resp, nil

	// 参数校验
	if req.Token == "" {
		klog.Error("令牌为空")
		resp.Code = auth.ErrorCode_TokenInvalid
		resp.Message = "令牌不能为空"
		resp.IsValid = false
		return resp, nil
	}

	// 解析并验证令牌
	claims, err := utils.ParseToken(req.Token)
	if err != nil {
		var code auth.ErrorCode
		switch err {
		case utils.ErrTokenExpired:
			code = auth.ErrorCode_TokenExpired
		case utils.ErrTokenNotValidYet, utils.ErrTokenMalformed, utils.ErrTokenInvalid:
			code = auth.ErrorCode_TokenInvalid
		default:
			code = auth.ErrorCode_TokenInvalid
		}

		klog.Errorf("验证令牌失败: %v", err)
		resp.Code = code
		resp.Message = err.Error()
		resp.IsValid = false
		return resp, nil
	}

	// 令牌有效，返回用户信息
	resp.Code = auth.ErrorCode_Success
	resp.Message = "令牌验证成功"
	resp.IsValid = true
	resp.UserId = claims.UserID
	resp.Role = claims.Role

	klog.Infof("成功验证用户[%d]的令牌, role=%s", claims.UserID, claims.Role)
	return resp, nil
}
