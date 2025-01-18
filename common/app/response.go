package app

import (
	"github.com/gin-gonic/gin"
	"github.com/loveRyujin/go-mall/common/errcode"
	"github.com/loveRyujin/go-mall/common/logger"
)

type Response struct {
	ctx        *gin.Context
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	RequestId  string      `json:"request_id"`
	Data       interface{} `json:"data,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

func NewResponse(c *gin.Context) *Response {
	return &Response{ctx: c}
}

func (r *Response) SetPagination(pagination *Pagination) *Response {
	r.Pagination = pagination
	return r
}

// Success 返回成功响应（有响应数据）
func (r *Response) Success(data interface{}) {
	r.Code = errcode.Success.Code()
	r.Message = errcode.Success.Message()
	requestId := ""
	if val, exists := r.ctx.Get("traceid"); exists {
		requestId = val.(string)
	}
	r.RequestId = requestId
	r.Data = data
	r.ctx.JSON(errcode.Success.HttpStatusCode(), r)
}

// SuccessOk 返回成功响应
func (r *Response) SuccessOk() {
	r.Success(nil)
}

// Error 返回错误响应
func (r *Response) Error(err *errcode.AppError) {
	r.Code = err.Code()
	r.Message = err.Message()
	requestId := ""
	if val, exists := r.ctx.Get("traceid"); exists {
		id, ok := val.(string)
		if !ok {
			logger.New(r.ctx).Warn("traceid type error", "traceid", val)
		} else {
			requestId = id
		}
	}
	r.RequestId = requestId
	// 兜底记录错误日志，项目自定义的AppError中有错误链条，方便出错后排查问题
	logger.New(r.ctx).Error("api_response_error", "err", err)
	r.ctx.JSON(err.HttpStatusCode(), r)
}
