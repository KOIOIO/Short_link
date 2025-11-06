package logic

import (
    "context"
    "errors"
    "github.com/shorturl/short-url/zero_remake/user-rpc/types/User"
    "strings"

    "github.com/shorturl/short-url/zero_remake/shorturl_api/internal/svc"
    "github.com/shorturl/short-url/zero_remake/shorturl_api/internal/types"

    "github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.UserRegisterRequset) (resp *types.UserRegisterResponse, err error) {
	// todo: add your logic here and delete this line
	username := strings.TrimSpace(req.Username)
	password := strings.TrimSpace(req.Password)
	confirmPassword := strings.TrimSpace(req.ConfirmPassword)

	RegisterResponse, err := l.svcCtx.UserRpc.Register(l.ctx, &User.RegisterRequest{
		Username:        username,
		Password:        password,
		ConfirmPassword: confirmPassword,
	})

	if RegisterResponse.Code != 200 {
		return nil, errors.New("register failed")
	}

	resp = &types.UserRegisterResponse{
		Code: int(RegisterResponse.Code),
	}

	return resp, nil
}
