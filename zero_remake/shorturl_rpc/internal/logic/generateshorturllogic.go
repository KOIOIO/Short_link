package logic

import (
    "context"
    "crypto/sha1"
    "errors"
    "fmt"
    "strconv"
    "strings"
    "time"

    "github.com/shorturl/short-url/zero_remake/common/errmsg"
    "github.com/shorturl/short-url/zero_remake/models"
    "github.com/shorturl/short-url/zero_remake/shorturl_rpc/internal/logic/repository"
    "github.com/shorturl/short-url/zero_remake/shorturl_rpc/internal/svc"
    "github.com/shorturl/short-url/zero_remake/shorturl_rpc/types/shortUrl"
    "gorm.io/gorm"

    "github.com/zeromicro/go-zero/core/logx"
)

type GenerateShortUrlLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGenerateShortUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateShortUrlLogic {
	return &GenerateShortUrlLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GenerateShortUrlLogic) GenerateShortUrl(in *shortUrl.GenerateShortUrlRequest) (*shortUrl.GenerateShortUrlResponse, error) {

	//做一个流量控制器
	// 防止恶意攻击下数据量崩溃
	// ip := limite_processer.GetClientIP(l.ctx)
	// if !limite_processer.AllowIP(ip) {
	// 	return &shortUrl.GenerateShortUrlResponse{
	// 		Code:      errmsg.ERROR_RATE_LIMIT,
	// 		Shortcode: "",
	// 	}, errors.New("rate limit exceeded")
	// }

    if in.Url == "" {
        return &shortUrl.GenerateShortUrlResponse{
            Code:      errmsg.ERROR_URL_IS_NULL,
            Shortcode: "",
        }, errors.New("url is null")
    }

    // 基于Redis的分布式锁，按URL维度防止并发插入造成唯一索引冲突
    lockKey := fmt.Sprintf("shorturl:lock:gen:%x", sha1.Sum([]byte(in.Url)))
    lockVal := fmt.Sprintf("%d", time.Now().UnixNano())
    acquired := false
    for i := 0; i < 10; i++ { // 最长等待约1秒
        ok, err := l.svcCtx.Redis.Rdb.SetNX(l.svcCtx.Redis.Ctx, lockKey, lockVal, 5*time.Second).Result()
        if err != nil {
            return &shortUrl.GenerateShortUrlResponse{Code: errmsg.ERROR, Shortcode: ""}, errors.New("lock acquire failed")
        }
        if ok {
            acquired = true
            break
        }
        time.Sleep(100 * time.Millisecond)
    }
    if !acquired {
        return &shortUrl.GenerateShortUrlResponse{Code: errmsg.ERROR, Shortcode: ""}, errors.New("system busy, please retry")
    }
    defer func() {
        // 仅当值匹配时释放锁，避免误删他人锁
        const unlockScript = `if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("del", KEYS[1]) else return 0 end`
        _ = l.svcCtx.Redis.Rdb.Eval(l.svcCtx.Redis.Ctx, unlockScript, []string{lockKey}, lockVal).Err()
    }()

	// 先用布隆过滤器检查URL是否存在
	// 如果存在，则直接从MySQL中读取短链
	// 如果不存在，则生成新的短链ID，保存到MySQL和Redis中，并加入布隆过滤器
	exist, _ := l.svcCtx.RedisBloom.MightContain(in.Url)

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
	id, err := repository.GetMyFlake().NextID()
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
		expireDuration, err = parseCustomDuration(in.Expiration)
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

	// 新增后再加入布隆过滤器
	_ = l.svcCtx.RedisBloom.Add(in.Url)

	return &shortUrl.GenerateShortUrlResponse{
		Code:      errmsg.SUCCESS,
		Shortcode: shortCode,
	}, nil
}

// parseCustomDuration 自定义函数，支持解析包含 'd' 单位的时间字符串
func parseCustomDuration(s string) (time.Duration, error) {
	if strings.HasSuffix(s, "d") {
		daysStr := strings.TrimSuffix(s, "d")
		days, err := strconv.Atoi(daysStr)
		if err != nil {
			return 0, err
		}
		return time.Duration(days) * 24 * time.Hour, nil
	}
	return time.ParseDuration(s)
}

// DeleteWithTime 删除数据库中创建时间超过一个月的记录
// 这个函数用于定期清理过期的短URL记录
func (l *GenerateShortUrlLogic) DeleteWithTime() error {
	err := l.svcCtx.DB.Where("created_at < ?", time.Now().Add(-time.Hour*24*30)).Delete(&models.Shorturl{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (l *GenerateShortUrlLogic) ReadFormMysql(url string) (*models.Shorturl, error) {
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

func (l *GenerateShortUrlLogic) SaveToMysql(shorturl models.Shorturl) error {
	err := l.svcCtx.DB.Create(&shorturl).Error
	if err != nil {
		return err
	} else {
		return nil
	}
}
