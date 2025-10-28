package logic

import (
	"context"
	"errors"
	"example.com/shorturl/short-url/zero_remake/common/errmsg"
	"example.com/shorturl/short-url/zero_remake/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"example.com/shorturl/short-url/zero_remake/user-rpc/internal/svc"
	"example.com/shorturl/short-url/zero_remake/user-rpc/types/User"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *User.LoginRequest) (*User.LoginResponse, error) {
	// todo: add your logic here and delete this line
	if in.Password == "" || in.Username == "" {
		return &User.LoginResponse{
			Code:    errmsg.ERROR_FAIL_TO_LOGIN,
			Message: errmsg.GetErrMsg(errmsg.ERROR_FAIL_TO_LOGIN),
		}, nil
	}

	user, err := l.GetFromMysql(in.Username)
	if err != nil {
		return &User.LoginResponse{
			Code:    errmsg.ERROR_FAIL_TO_LOGIN,
			Message: errmsg.GetErrMsg(errmsg.ERROR_FAIL_TO_LOGIN),
		}, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Password))
	if err != nil {
		return &User.LoginResponse{
			Code:    errmsg.ERROR_FAIL_TO_LOGIN,
			Message: errmsg.GetErrMsg(errmsg.ERROR_FAIL_TO_LOGIN),
		}, nil
	}

	return &User.LoginResponse{
		Code:    200,
		Message: "success",
	}, nil
}

func (l *LoginLogic) GetFromMysql(username string) (*models.Usermodel, error) {
	var user models.Usermodel
	err := l.svcCtx.DB.Where("username=?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (l *LoginLogic) checkPassword(hashedPassword, inputPassword string) bool {
	// 比对哈希值和输入密码
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
	return err == nil // 无错误则匹配成功
}
