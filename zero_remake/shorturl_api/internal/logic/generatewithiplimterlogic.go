package logic

import (
    "context"
    "errors"
    "net/http"

    "github.com/shorturl/short-url/zero_remake/common/errmsg"
    "github.com/shorturl/short-url/zero_remake/shorturl_rpc/types/shortUrl"

    "strings"

    "github.com/shorturl/short-url/zero_remake/shorturl_api/internal/svc"
    "github.com/shorturl/short-url/zero_remake/shorturl_api/internal/types"

    "github.com/zeromicro/go-zero/core/logx"
)

type GenerateWithIPLimterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenerateWithIPLimterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateWithIPLimterLogic {
	return &GenerateWithIPLimterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenerateWithIPLimterLogic) GenerateWithIPLimter(req *types.GenerateRequest) (resp *types.GenerateResponse, err error) {
	// todo: add your logic here and delete this line
	url := strings.TrimSpace(req.Url)
	expiration := strings.TrimSpace(req.Expiration)

	GenerateShortUrlResponse, _ := l.svcCtx.ShortUrlRpc.GenerateWithIPLimter(l.ctx, &shortUrl.GenerateShortUrlRequest{
		Url:        url,
		Expiration: expiration,
	})
	if GenerateShortUrlResponse.Code == errmsg.SUCCESS {
		return &types.GenerateResponse{
			Code:     errmsg.SUCCESS,
			ShortUrl: GenerateShortUrlResponse.Shortcode,
			Message:  "短链接生成成功",
		}, nil
	} else if GenerateShortUrlResponse.Code == errmsg.ERROR_RATE_LIMIT {
		return &types.GenerateResponse{
			Code:     errmsg.ERROR_RATE_LIMIT,
			ShortUrl: "",
			Message:  "请求过于频繁，请稍后再试",
		}, errors.New("请求过于频繁，请稍后再试")
	} else {
		return &types.GenerateResponse{
			Code:     http.StatusBadRequest,
			ShortUrl: "",
			Message:  "生成短链接失败",
		}, errors.New("生成短链接失败")
	}
}
