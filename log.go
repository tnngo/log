package log

import (
	"io"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 日志时间格式。
func timeFormat(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// FileCfg 文件内容配置。
type FileCfg struct {
	// Filename 文件名。
	Filename string
	// MaxSize 最大尺寸，默认为100MB。
	MaxSize int
	// MaxBackups 备份数量。
	MaxBackups int
	// MaxAge 最大保存时间。
	MaxAge int
	// Compress 是否压缩打包。
	Compress bool

	// Level 日志级别，默认为info。
	ZapLevel zapcore.Level
}

// ConsoleCfg 控制台内容配置。
type ConsoleCfg struct {
	// Level 日志级别，默认为info。
	ZapLevel zapcore.Level
}

// Options 配置选项。
type Options struct {
	// 文件配置，默认为null。
	File *FileCfg
	// 控制台输出，默认不为null。
	Console *ConsoleCfg
}

// 用于初始化控制台输出。
func newConsole(opt *Options) zapcore.Core {
	consoleWrite := zapcore.AddSync(io.MultiWriter(os.Stdout))
	consoleConfig := zap.NewProductionEncoderConfig()
	consoleConfig.EncodeTime = timeFormat
	// 控制台输出颜色。
	consoleConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	/** 定义日志控制台输出核心。*/
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(consoleConfig),
		consoleWrite,
		opt.Console.ZapLevel,
	)
	return consoleCore

}

var logger *zap.Logger

// NewSimple 新建一个简单的日志对象。
// 仅控制台输出。
func NewSimple() *zap.Logger {
	var cores []zapcore.Core
	opt := &Options{
		Console: &ConsoleCfg{
			ZapLevel: zapcore.DebugLevel,
		},
	}
	cores = append(cores, newConsole(opt))
	core := zapcore.NewTee(cores...)

	logger = zap.New(core, zap.AddCaller())
	return logger
}

// New 新建日志。
func New(opt *Options) *zap.Logger {
	if opt == nil {
		panic("log Options不能nil。")
	}

	var cores []zapcore.Core
	if opt.Console == nil {
		opt.Console = &ConsoleCfg{
			ZapLevel: zapcore.DebugLevel,
		}
	}
	cores = append(cores, newConsole(opt))

	/** 定义日志文件输出核心。 */
	if opt.File != nil {
		hook := &lumberjack.Logger{
			Filename:   opt.File.Filename,
			MaxSize:    opt.File.MaxSize,
			MaxBackups: opt.File.MaxBackups,
			MaxAge:     opt.File.MaxAge,
			Compress:   opt.File.Compress,
		}

		fileWrite := zapcore.AddSync(io.MultiWriter(hook))
		fileConfig := zap.NewProductionEncoderConfig()
		fileConfig.EncodeTime = timeFormat
		fileCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(fileConfig),
			fileWrite,
			opt.File.ZapLevel,
		)
		cores = append(cores, fileCore)
	}
	core := zapcore.NewTee(cores...)

	logger = zap.New(core, zap.AddCaller())
	return logger
}

// L 全局logger对象。
func L() *zap.Logger {
	return logger
}
