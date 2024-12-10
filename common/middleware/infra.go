package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/loveRyujin/go-mall/common/logger"
	"github.com/loveRyujin/go-mall/common/utils"
	"io"
	"time"
)

const (
	AccessStart = "access_start"
	AccessEnd   = "access_end"
)

// StartTrace 开启链路追踪
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

// LogAccess 记录进出日志
func LogAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		body, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewReader(body))
		startTime := time.Now()
		blw := &bodyLogWriter{
			ResponseWriter: c.Writer,
			buf:            bytes.NewBufferString(""),
		}
		c.Writer = blw
		accessLog(c, AccessStart, time.Since(startTime), body, nil)
		defer func() {
			accessLog(c, AccessEnd, time.Since(startTime), body, blw.buf.String())
		}()
		c.Next()
		return
	}
}

func accessLog(c *gin.Context, accessType string, dur time.Duration, b []byte, dataOut interface{}) {
	req := c.Request
	bodyStr := string(b)
	method := req.Method
	cip := c.ClientIP()
	query := req.URL.RawQuery
	path := req.URL.Path
	logger.New(c).Info("AccessLog",
		"type", accessType,
		"ip", cip,
		"method", method,
		"path", path,
		"query", query,
		"body", bodyStr,
		"output", dataOut,
		"time(ms)", int64(dur/time.Millisecond),
	)
}

// 包装gin.ResponseWriter,暂存响应体
type bodyLogWriter struct {
	gin.ResponseWriter
	buf *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.buf.Write(b)
	return w.ResponseWriter.Write(b)
}
