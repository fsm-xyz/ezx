package log

import (
	"context"
	"os"

	"github.com/fsm-xyz/ezx/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger ...
type Logger struct {
	logger *zap.Logger
}

// New 返回一个Logger
func New() *Logger {
	writer := getLogWriter()
	encoder := getEncoder()

	core := zapcore.NewCore(encoder, writer, getLevel(config.C.Log.Level))

	if config.C.Log.Dev {
		return &Logger{
			logger: zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.Development()),
		}
	}
	return &Logger{logger: zap.New(core)}
}

var levelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

// 返回指定的level， 默认info
func getLevel(level string) zapcore.Level {
	if x, ok := levelMap[level]; ok {
		return x
	}
	return zap.InfoLevel
}

func getEncoder() zapcore.Encoder {
	c := zap.NewProductionEncoderConfig()
	c.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	c.EncodeCaller = zapcore.FullCallerEncoder
	return zapcore.NewJSONEncoder(c)
}

// 日志输出位置
const (
	stdoutOutput = "stdout"
	stderrOutput = "stderr"
	fileOutput   = "file"
)

func getLogWriter() zapcore.WriteSyncer {
	c := config.C.Log
	if c.Output == stdoutOutput {
		return os.Stdout
	}
	if c.Output == stderrOutput {
		return os.Stderr
	}

	lumberJackLogger := &lumberjack.Logger{
		Filename:   c.Rotate.Filename,
		MaxSize:    int(c.Rotate.MaxSize),
		MaxAge:     int(c.Rotate.MaxAge),
		MaxBackups: int(c.Rotate.MaxBackups),
		LocalTime:  c.Rotate.LocalTime,
		Compress:   c.Rotate.Compress,
	}
	return zapcore.AddSync(lumberJackLogger)
}

// Debug ...
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

// Info ...
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

// Warn ...
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

// Error ...
func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

// DPanic ...
func (l *Logger) DPanic(msg string, fields ...zap.Field) {
	l.logger.DPanic(msg, fields...)
}

// Panic ...
func (l *Logger) Panic(msg string, fields ...zap.Field) {
	l.logger.Panic(msg, fields...)
}

// Fatal ...
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

// With ...
func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{l.logger.With(fields...)}
}

// Sync ...
func (l *Logger) Sync() error {
	return l.logger.Sync()
}

// 包级别的logger
var (
	defaultL *Logger
)

// SetDefault 设置默认
func SetDefault(logger *Logger) {
	defaultL = logger
}

// GetLogger 提供默认logger, 实现自定义修改
func GetLogger() *Logger {
	return defaultL
}

// NewWith ...
func NewWith(ctx context.Context, fields ...zap.Field) *Logger {
	fields = append(fields, getTrace(ctx))
	return &Logger{logger: defaultL.logger.With(fields...)}
}

// Debug ...
func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, getTrace(ctx))
	defaultL.logger.Debug(msg, fields...)
}

// Info ...
func Info(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, getTrace(ctx))
	defaultL.logger.Info(msg, fields...)
}

// Warn ...
func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, getTrace(ctx))
	defaultL.logger.Warn(msg, fields...)
}

// Error ...
func Error(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, getTrace(ctx))
	defaultL.logger.Error(msg, fields...)
}

// DPanic ...
func DPanic(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, getTrace(ctx))
	defaultL.logger.DPanic(msg, fields...)
}

// Panic ...
func Panic(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, getTrace(ctx))
	defaultL.logger.Panic(msg, fields...)
}

// Fatal ...
func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, getTrace(ctx))
	defaultL.logger.Fatal(msg, fields...)
}

// Sync ...
func Sync() error {
	return defaultL.logger.Sync()
}
