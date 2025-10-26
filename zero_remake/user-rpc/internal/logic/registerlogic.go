package logic

import (
	"context"
	"errors"
	"example.com/shorturl/short-url/zero_remake/common/errmsg"
	"example.com/shorturl/short-url/zero_remake/models"
	"gorm.io/gorm"

	"example.com/shorturl/short-url/zero_remake/user-rpc/internal/svc"
	"example.com/shorturl/short-url/zero_remake/user-rpc/types/User"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *User.RegisterRequest) (*User.RegisterResponse, error) {
	// todo: add your logic here and delete this line
	if in.Password == "" || in.Username == "" || in.ConfirmPassword == "" {
		return &User.RegisterResponse{
			Code:    errmsg.ERROR_FAIL_TO_REGISTER,
			Message: errmsg.GetErrMsg(errmsg.ERROR_FAIL_TO_REGISTER),
		}, errors.New("invalid params")
	}

	if in.Password != in.ConfirmPassword {
		return &User.RegisterResponse{
			Code:    errmsg.ERROR_FAIL_TO_REGISTER,
			Message: errmsg.GetErrMsg(errmsg.ERROR_FAIL_TO_REGISTER),
		}, nil
	}
	exist, err := l.CheckExist(in.Username)
	if err != nil {
		return &User.RegisterResponse{
			Code:    errmsg.ERROR_FAIL_TO_REGISTER,
			Message: errmsg.GetErrMsg(errmsg.ERROR_FAIL_TO_REGISTER),
		}, nil
	}
	if exist {
		return &User.RegisterResponse{
			Code:    errmsg.ERROR_FAIL_TO_REGISTER,
			Message: errmsg.GetErrMsg(errmsg.ERROR_FAIL_TO_REGISTER),
		}, nil
	}
	userinfo := models.Usermodel{
		Username: in.Username,
		Password: in.Password,
	}

	err = l.SaveToDB(&userinfo)
	if err != nil {
		return &User.RegisterResponse{
			Code:    errmsg.ERROR_FAIL_TO_REGISTER,
			Message: errmsg.GetErrMsg(errmsg.ERROR_FAIL_TO_REGISTER),
		}, nil
	}

	return &User.RegisterResponse{
		Code:    200,
		Message: "success",
	}, nil
}

func (l *RegisterLogic) SaveToDB(user *models.Usermodel) error {
	return l.svcCtx.DB.Create(user).Error
}

func (l *RegisterLogic) CheckExist(username string) (bool, error) {
	var user models.Usermodel
	err := l.svcCtx.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
