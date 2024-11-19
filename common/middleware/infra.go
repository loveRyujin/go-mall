package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/loveRyujin/go-mall/common/utils"
)

func StartTrace() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := c.Request.Header.Get("traceid")
		pSpanId := c.Request.Header.Get("spanid")
		spanId := utils.GenSpanID(c.Request.RemoteAddr)
		if traceId == "" {
			traceId = spanId
		}
		c.Set("traceid", traceId)
		c.Set("pspanid", pSpanId)
		c.Set("spanid", spanId)
		c.Next()
	}
}
