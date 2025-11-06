package handler

import (
    "net/http"

    "github.com/shorturl/short-url/zero_remake/shorturl_api/internal/logic"
    "github.com/shorturl/short-url/zero_remake/shorturl_api/internal/svc"
    "github.com/shorturl/short-url/zero_remake/shorturl_api/internal/types"
    "github.com/zeromicro/go-zero/rest/httpx"
)

func GenerateIDbMySnowflakeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GenerateRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGenerateIDbMySnowflakeLogic(r.Context(), svcCtx)
		resp, err := l.GenerateIDbMySnowflake(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
