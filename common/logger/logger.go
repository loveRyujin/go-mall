package logger

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"path"
	"runtime"
)

type logger struct {
	ctx     context.Context
	traceId string
	spanId  string
	pSpanId string
	_logger *zap.Logger
}

func New(ctx context.Context) *logger {
	var traceId, spanId, pSpanId string
	if ctx.Value("traceid") != nil {
		traceId = ctx.Value("traceid").(string)
	}
	if ctx.Value("spanid") != nil {
		spanId = ctx.Value("spanid").(string)
	}
	if ctx.Value("pspanid") != nil {
		pSpanId = ctx.Value("pspanid").(string)
	}
	return &logger{
		ctx:     ctx,
		traceId: traceId,
		spanId:  spanId,
		pSpanId: pSpanId,
		_logger: _logger,
	}
}

func (l *logger) log(level zapcore.Level, msg string, kv ...interface{}) {
	if len(kv)%2 != 0 {
		// 凑齐参数
		kv = append(kv, "unknown")
	}
	// 追加链路追踪参数
	kv = append(kv, "traceid", l.traceId, "spanid", l.spanId, "pspanid", l.pSpanId)
	fields := make([]zap.Field, 0, len(kv)/2)
	// 追加调用者位置信息
	fn, file, line := l.GetLoggerCallerInfo()
	kv = append(kv, "func", fn, "file", file, "line", line)
	for i := 0; i < len(kv); i += 2 {
		k := fmt.Sprintf("%v", kv[i])
		fields = append(fields, zap.Any(k, kv[i+1]))
	}
	ce := l._logger.Check(level, msg)
	ce.Write(fields...)
}

func (l *logger) GetLoggerCallerInfo() (fn, file string, line int) {
	pc, file, line, ok := runtime.Caller(3)
	if !ok {
		return
	}
	file = path.Base(file)
	fn = runtime.FuncForPC(pc).Name()
	return
}

func (l *logger) Info(msg string, kv ...interface{}) {
	l.log(zapcore.InfoLevel, msg, kv...)
}

func (l *logger) Debug(msg string, kv ...interface{}) {
	l.log(zapcore.DebugLevel, msg, kv...)
}

func (l *logger) Warn(msg string, kv ...interface{}) {
	l.log(zapcore.WarnLevel, msg, kv...)
}

func (l *logger) Error(msg string, kv ...interface{}) {
	l.log(zapcore.ErrorLevel, msg, kv...)
}
