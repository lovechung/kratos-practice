package zlog

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"kratos-practice/internal/conf"
	"os"
)

var _ log.Logger = (*ZapLogger)(nil)

// ZapLogger is a logger impl.
type ZapLogger struct {
	log  *zap.Logger
	Sync func() error
}

// Logger 配置zap日志,将zap日志库引入
func Logger(conf *conf.Project) log.Logger {
	//配置zap日志库的编码器
	encoder := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}

	return NewZapLogger(
		conf,
		encoder,
		zap.NewAtomicLevelAt(zapcore.DebugLevel),
		zap.AddStacktrace(
			zap.NewAtomicLevelAt(zapcore.ErrorLevel)),
		zap.AddCaller(),
		zap.AddCallerSkip(2),
		zap.Development(),
	)
}

// NewZapLogger return a zap logger.
func NewZapLogger(conf *conf.Project, encoder zapcore.EncoderConfig, level zap.AtomicLevel, opts ...zap.Option) *ZapLogger {
	//日志切割
	writeSyncer := getLogWriter(conf)
	//设置日志级别
	level.SetLevel(zap.InfoLevel)
	var core zapcore.Core
	//开发模式下打印到标准输出
	// --根据配置文件判断输出到控制台还是日志文件--
	if conf.Profile == "dev" {
		core = zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoder),                      // 编码器配置
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), // 打印到控制台
			level, // 日志级别
		)
	} else {
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoder), // 编码器配置
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(writeSyncer)), // 打印到控制台和文件
			level, // 日志级别
		)
	}
	zapLogger := zap.New(core, opts...)
	return &ZapLogger{log: zapLogger, Sync: zapLogger.Sync}
}

// 日志自动切割，采用 lumberjack 实现的
func getLogWriter(conf *conf.Project) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   conf.Log.File,
		MaxSize:    100,   //日志的最大大小（M）
		MaxBackups: 5,     //日志的最大保存数量
		MaxAge:     30,    //日志文件存储最大天数
		Compress:   false, //是否执行压缩
	}
	return zapcore.AddSync(lumberJackLogger)
}

// Log 实现log接口
func (l *ZapLogger) Log(level log.Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 || len(keyvals)%2 != 0 {
		l.log.Warn(fmt.Sprint("Keyvalues must appear in pairs: ", keyvals))
		return nil
	}
	// Zap.Field is used when keyvals pairs appear
	var data []zap.Field
	for i := 0; i < len(keyvals); i += 2 {
		data = append(data, zap.Any(fmt.Sprint(keyvals[i]), fmt.Sprint(keyvals[i+1])))
	}
	switch level {
	case log.LevelDebug:
		l.log.Debug("", data...)
	case log.LevelInfo:
		l.log.Info("", data...)
	case log.LevelWarn:
		l.log.Warn("", data...)
	case log.LevelError:
		l.log.Error("", data...)
	}
	return nil
}
