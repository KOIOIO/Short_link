package logic

import (
	"context"
	"errors"
	"time"

	"github.com/shorturl/short-url/zero_remake/common/errmsg"
	"github.com/shorturl/short-url/zero_remake/models"
	"github.com/shorturl/short-url/zero_remake/shorturl_rpc/internal/logic/repository"
	"gorm.io/gorm"

	"github.com/shorturl/short-url/zero_remake/shorturl_rpc/internal/svc"
	"github.com/shorturl/short-url/zero_remake/shorturl_rpc/types/shortUrl"

	"github.com/zeromicro/go-zero/core/logx"
)

type FilterByMyBloomFilterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFilterByMyBloomFilterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FilterByMyBloomFilterLogic {
	return &FilterByMyBloomFilterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FilterByMyBloomFilterLogic) FilterByMyBloomFilter(in *shortUrl.GenerateShortUrlRequest) (*shortUrl.GenerateShortUrlResponse, error) {
	if in.Url == "" {
		return &shortUrl.GenerateShortUrlResponse{
			Code:      errmsg.ERROR_URL_IS_NULL,
			Shortcode: "",
		}, errors.New("url is null")
	}

	// 先用布隆过滤器检查URL是否存在
	// 如果存在，则直接从MySQL中读取短链
	// 如果不存在，则生成新的短链ID，保存到MySQL和Redis中，并加入布隆过滤器
	exist := repository.Bloom.MightContain(in.Url)

	if exist {
		shortcode, err := l.ReadFormMysql(in.Url)
		if err == nil && shortcode != nil {
			return &shortUrl.GenerateShortUrlResponse{
				Code:      errmsg.SUCCESS,
				Shortcode: shortcode.Shorturl,
			}, nil
		}
	}

	// 生成短链ID
	id, err := repository.IDGnerator.NextID()
	if err != nil {
		return &shortUrl.GenerateShortUrlResponse{
			Code:      errmsg.ERROR,
			Shortcode: "",
		}, errors.New("ID生成失败")
	}

	shortCode := repository.Base62Encode(id)

	// 解析过期时间
	var expireDuration time.Duration
	if in.Expiration != "" {
		// 使用自定义函数解析过期时间
		expireDuration, err = repository.Duration.ParseCustomDuration(in.Expiration)
		if err != nil {
			return &shortUrl.GenerateShortUrlResponse{
				Code:      errmsg.ERROR_EXPIRATION_ID_WRONG,
				Shortcode: "",
			}, errors.New("failed to parse expiration time")
		}
	}

	shorturl := models.Shorturl{
		ID:       id,
		Shorturl: shortCode,
		Url:      in.Url,
	}

	if err := l.svcCtx.DB.Create(&shorturl).Error; err != nil {
		return &shortUrl.GenerateShortUrlResponse{
			Code:      errmsg.ERROR_FAILED_TO_SAVE_TO_MYSQL,
			Shortcode: "",
		}, errors.New("fail to save to mysql")
	}

	if err := l.svcCtx.Redis.Rdb.Set(l.svcCtx.Redis.Ctx, shortCode, in.Url, expireDuration).Err(); err != nil {
		return &shortUrl.GenerateShortUrlResponse{
			Code:      errmsg.ERROR_FAILED_SAVE_TO_REDIS,
			Shortcode: "",
		}, errors.New("fail to save to redis")
	}

	// 新增后再加入内置布隆过滤器（使用项目内存版 Bloom）
	repository.Bloom.Add(in.Url)

	return &shortUrl.GenerateShortUrlResponse{
		Code:      errmsg.SUCCESS,
		Shortcode: shortCode,
	}, nil
}

// DeleteWithTime 删除数据库中创建时间超过一个月的记录
// 这个函数用于定期清理过期的短URL记录
func (l *FilterByMyBloomFilterLogic) DeleteWithTime() error {
	err := l.svcCtx.DB.Where("created_at < ?", time.Now().Add(-time.Hour*24*30)).Delete(&models.Shorturl{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (l *FilterByMyBloomFilterLogic) ReadFormMysql(url string) (*models.Shorturl, error) {
	var shortURL models.Shorturl
	err := l.svcCtx.DB.Where("url = ?", url).First(&shortURL).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &shortURL, nil
}

func (l *FilterByMyBloomFilterLogic) SaveToMysql(shorturl models.Shorturl) error {
	err := l.svcCtx.DB.Create(&shorturl).Error
	if err != nil {
		return err
	} else {
		return nil
	}
}
