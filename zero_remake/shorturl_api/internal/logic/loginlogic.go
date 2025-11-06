package logic

import (
    "context"
    "errors"
    jwts "github.com/shorturl/short-url/zero_remake/common/Auth"
    "github.com/shorturl/short-url/zero_remake/shorturl_api/internal/svc"
    "github.com/shorturl/short-url/zero_remake/shorturl_api/internal/types"
    "github.com/shorturl/short-url/zero_remake/user-rpc/types/User"
    "strings"

    "github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.UserLoginRequest) (resp *types.UserLoginResponse, err error) {
	// todo: add your logic here and delete this line
	username := strings.TrimSpace(req.Username)
	password := strings.TrimSpace(req.Password)
	auth := l.svcCtx.Config.Auth
	loginResp, err := l.svcCtx.UserRpc.Login(l.ctx, &User.LoginRequest{
		Username: username,
		Password: password,
	})
	if loginResp.Code != 200 {
		return nil, errors.New(loginResp.Message)
	}

	if err != nil {
		return nil, err
	}

	token, err := jwts.GenToken(jwts.JwtPayLoad{
		Code:     200,
		Username: req.Username,
	}, auth.AccessSecret, auth.AccessExpire)
	if err != nil {
		return nil, err
	}
	return &types.UserLoginResponse{
		Code:        200,
		AccessToken: token,
		ExpiresIn:   auth.AccessExpire,
		Message:     "login success",
	}, nil
}
