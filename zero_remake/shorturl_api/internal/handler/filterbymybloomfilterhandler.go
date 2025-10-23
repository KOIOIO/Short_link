package handler

import (
	"net/http"

	"example.com/shorturl/short-url/zero_remake/shorturl_api/internal/logic"
	"example.com/shorturl/short-url/zero_remake/shorturl_api/internal/svc"
	"example.com/shorturl/short-url/zero_remake/shorturl_api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func FilterByMyBloomFilterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GenerateRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewFilterByMyBloomFilterLogic(r.Context(), svcCtx)
		resp, err := l.FilterByMyBloomFilter(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
