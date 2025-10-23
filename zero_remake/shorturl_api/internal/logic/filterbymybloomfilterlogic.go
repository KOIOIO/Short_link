package logic

import (
	"context"
	"errors"
	"net/http"

	"example.com/shorturl/short-url/zero_remake/common/errmsg"
	"example.com/shorturl/short-url/zero_remake/shorturl_rpc/types/shortUrl"

	"strings"

	"example.com/shorturl/short-url/zero_remake/shorturl_api/internal/svc"
	"example.com/shorturl/short-url/zero_remake/shorturl_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FilterByMyBloomFilterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFilterByMyBloomFilterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FilterByMyBloomFilterLogic {
	return &FilterByMyBloomFilterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FilterByMyBloomFilterLogic) FilterByMyBloomFilter(req *types.GenerateRequest) (resp *types.GenerateResponse, err error) {
	// todo: add your logic here and delete this line
	url := strings.TrimSpace(req.Url)
	expiration := strings.TrimSpace(req.Expiration)

	GenerateShortUrlResponse, _ := l.svcCtx.ShortUrlRpc.FilterByMyBloomFilter(l.ctx, &shortUrl.GenerateShortUrlRequest{
		Url:        url,
		Expiration: expiration,
	})
	if GenerateShortUrlResponse.Code == errmsg.SUCCESS {
		return &types.GenerateResponse{
			Code:     errmsg.SUCCESS,
			ShortUrl: GenerateShortUrlResponse.Shortcode,
			Message:  "短链接生成成功",
		}, nil
	} else {
		return &types.GenerateResponse{
			Code:     http.StatusBadRequest,
			ShortUrl: "",
			Message:  "生成短链接失败",
		}, errors.New("生成短链接失败")
	}
}
