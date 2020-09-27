package logger

import (
	"github.com/ZYallers/zgin/app"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
	"sync/atomic"
	"time"
)

const (
	fileMaxSize  = 100 // Unit: MB
	maxBackups   = 20
	fileSuffix   = `.log`
	cacheMaxSize = 100
)

var (
	logDir                   = app.LogDir
	loggerCache, lumberCache sync.Map
	loggerCounter            int32
	levelEnabler             zap.LevelEnablerFunc = func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	}
	jsonEncoder = zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:       "ts",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006/01/02 15:04:05.000"))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
)

func AppLogger() *zap.Logger {
	file := logDir + "/" + app.Name + fileSuffix
	logger, _ := newLogger(file)
	return logger
}

func Use(filename string) *zap.Logger {
	if filename == "" {
		return nil
	}

	file := logDir + "/" + app.Name + "-" + filename + fileSuffix

	// 判断是否已存在缓存中
	if logger, ok := loggerCache.Load(file); ok {
		return logger.(*zap.Logger)
	}

	// 判断容器是否达到最大数量
	if atomic.LoadInt32(&loggerCounter) >= cacheMaxSize {
		// 超出最大容量，随机删除一半
		var counter, clean int32 = 0, cacheMaxSize / 2
		loggerCache.Range(func(key, value interface{}) bool {
			if counter++; counter > clean {
				return false
			}
			if lumber, ok := lumberCache.Load(key); ok {
				_ = lumber.(*lumberjack.Logger).Close()
			}
			lumberCache.Delete(key)
			loggerCache.Delete(key)
			return true
		})
		atomic.AddInt32(&loggerCounter, -clean)
	}

	logger, hook := newLogger(file)
	loggerCache.Store(file, logger)
	lumberCache.Store(file, hook)
	atomic.AddInt32(&loggerCounter, 1)

	return logger
}

func newLogger(filename string) (*zap.Logger, *lumberjack.Logger) {
	hook := &lumberjack.Logger{MaxSize: fileMaxSize, MaxBackups: maxBackups, LocalTime: true, Compress: false, Filename: filename}
	logger := zap.New(zapcore.NewCore(jsonEncoder, zapcore.AddSync(hook), levelEnabler))
	logger.Info("new logger succeed", zap.String("filename", hook.Filename))
	return logger, hook
}
