package zap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

func NewConsoleCore(level zapcore.Level) zapcore.Core {
	// 自定义控制台输出编码器（带颜色）
	consoleEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // 彩色编码
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
	return zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level)
}

func DefaultFileConfig() *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    100,  // 单个文件最大尺寸（MB）
		MaxBackups: 5,    // 最多保留旧文件个数
		MaxAge:     30,   // 保留旧文件最大天数
		Compress:   true, // 是否压缩旧文件
	}
}

func NewFileCore(level zapcore.Level, config *lumberjack.Logger) zapcore.Core {
	// JSON 文件输出编码器
	fileEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
	var fileWriter zapcore.WriteSyncer
	if config != nil {
		fileWriter = zapcore.AddSync(config)
	} else {
		fileWriter = zapcore.AddSync(DefaultFileConfig())
	}

	return zapcore.NewCore(fileEncoder, fileWriter, level)
}

func NewZapLogger() *zap.Logger {
	c1 := NewConsoleCore(zap.DebugLevel)
	c2 := NewFileCore(zap.DebugLevel, nil)

	// 设置日志级别
	core := zapcore.NewTee(
		c1, c2,
	)
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel), zap.AddCallerSkip(3))
}
