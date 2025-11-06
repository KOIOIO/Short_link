package logic

import (
	"context"
	"errors"
	"github.com/shorturl/short-url/zero_remake/common/errmsg"
	"github.com/shorturl/short-url/zero_remake/models"
	"github.com/shorturl/short-url/zero_remake/user-rpc/internal/svc"
	"github.com/shorturl/short-url/zero_remake/user-rpc/types/User"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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
	//检查用户名是否存在
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

	//密码加密
	bycrtypassword, err := l.hashPassword(in.Password)
	if err != nil {
		return &User.RegisterResponse{
			Code:    errmsg.ERROR_FAIL_TO_REGISTER,
			Message: errmsg.GetErrMsg(errmsg.ERROR_FAIL_TO_REGISTER),
		}, nil
	}
	userinfo := models.Usermodel{
		Username: in.Username,
		Password: bycrtypassword,
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
	// 查询数据库中是否存在该用户名
	err := l.svcCtx.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil // 用户名不存在
		}
		return false, err // 其他数据库错误
	}
	return true, nil // 用户名存在
}

func (l *RegisterLogic) hashPassword(password string) (string, error) {
	hashbytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashbytes), nil
}
