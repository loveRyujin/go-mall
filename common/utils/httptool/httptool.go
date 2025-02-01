package httptool

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/loveRyujin/go-mall/common/errcode"
	"github.com/loveRyujin/go-mall/common/logger"
	"github.com/loveRyujin/go-mall/common/utils"
	"io"
	"net/http"
	"time"
)

func Request(method, url string, options ...Option) (httpStatusCode int, respBody []byte, err error) {
	start := time.Now()
	requestOption := defaultRequestOption()
	for _, opt := range options {
		if err = opt.apply(requestOption); err != nil {
			return
		}
	}
	log := logger.New(requestOption.ctx)
	defer func() {
		if err != nil {
			log.Error("HTTP_REQUEST_ERROR_LOG", "method", method, "url", url, "body", requestOption.data, "reply", respBody, "err", err)
		}
	}()

	// 创建请求对象
	req, err := http.NewRequest(method, url, bytes.NewReader(requestOption.data))
	if err != nil {
		return
	}
	req = req.WithContext(requestOption.ctx)
	defer func() {
		_ = req.Body.Close()
	}()

	// 在Header中添加追踪信息 把内部服务串起来
	traceId, spanId, _ := utils.GetTraceInfoFromCtx(requestOption.ctx)
	requestOption.headers["traceid"] = traceId
	requestOption.headers["spanid"] = spanId
	if len(requestOption.headers) != 0 { // 设置请求头
		for key, value := range requestOption.headers {
			req.Header.Add(key, value)
		}
	}

	// 发起请求
	client := &http.Client{
		Timeout: requestOption.timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	dur := time.Since(start).Milliseconds()
	defer func() {
		if dur >= 3000 {
			// 记录慢请求日志
			log.Warn("HTTP_REQUEST_SLOW_LOG", "method", method, "url", url, "body", requestOption.data, "reply", respBody, "err", err, "dur/ms", dur)
		} else {
			log.Debug("HTTP_REQUEST_DEBUG_LOG", "method", method, "url", url, "body", string(requestOption.data), "reply", string(respBody), "err", err, "dur/ms", dur)
		}
	}()

	httpStatusCode = resp.StatusCode
	if httpStatusCode != http.StatusOK {
		// 返回非 200 时Go的 http 库不回返回error, 这里处理成error 调用方好判断
		err = errcode.Wrap("request api error", errors.New(fmt.Sprintf("non 200 response, response code: %d", httpStatusCode)))
		return
	}

	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return
}

// Get 发起GET请求
func Get(ctx context.Context, url string, options ...Option) (httpStatusCode int, respBody []byte, err error) {
	options = append(options, WithContext(ctx))
	return Request("GET", url, options...)
}

// Post 发起POST请求
func Post(ctx context.Context, url string, data []byte, options ...Option) (httpStatusCode int, respBody []byte, err error) {
	// 默认自带Header Content-Type: application/json 可通过 传递 WithHeaders 增加或者覆盖Header信息
	defaultHeader := map[string]string{"Content-Type": "application/json"}
	var newOptions []Option
	newOptions = append(newOptions, WithHeaders(defaultHeader), WithData(data), WithContext(ctx))
	newOptions = append(newOptions, options...)

	httpStatusCode, respBody, err = Request("POST", url, newOptions...)
	return
}

type Option interface {
	apply(option *requestOption) error
}

type optionFunc func(option *requestOption) error

func (f optionFunc) apply(opts *requestOption) error {
	return f(opts)
}

type requestOption struct {
	ctx     context.Context
	timeout time.Duration
	data    []byte
	headers map[string]string
}

func defaultRequestOption() *requestOption {
	return &requestOption{
		ctx:     context.Background(),
		timeout: 5 * time.Second,
		data:    nil,
		headers: make(map[string]string),
	}
}

func WithContext(ctx context.Context) Option {
	return optionFunc(func(opts *requestOption) error {
		opts.ctx = ctx
		return nil
	})
}

func WithTimeout(timeout time.Duration) Option {
	return optionFunc(func(opts *requestOption) error {
		opts.timeout = timeout
		return nil
	})
}

func WithData(data []byte) Option {
	return optionFunc(func(opts *requestOption) error {
		opts.data = data
		return nil
	})
}

func WithHeaders(headers map[string]string) Option {
	return optionFunc(func(opts *requestOption) error {
		for k, v := range headers {
			opts.headers[k] = v
		}
		return nil
	})
}
