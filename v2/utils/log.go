package utils

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Options struct {
	LogFileDir    string
	AppName       string
	ErrorFileName string
	WarnFileName  string
	InfoFileName  string
	DebugFileName string
	Level         zapcore.Level
	MaxSize       int
	MaxBackups    int
	MaxAge        int
	Development   bool
	zap.Config
}

const ContextReqUUid = "req_uuid"

type ModOptions func(options *Options)

var (
	l                              *Logger
	sp                             = string(filepath.Separator)
	errWs, warnWs, infoWs, debugWs zapcore.WriteSyncer
	debugConsoleWS                 = zapcore.Lock(os.Stdout)
	errorConsoleWS                 = zapcore.Lock(os.Stderr)
)

type Logger struct {
	*zap.Logger
	sync.Mutex
	Opts      *Options `json:"opts"`
	zapConfig zap.Config
	inited    bool
}

var logger *zap.Logger

func NewLoggerServer() {
	logger = NewLogger(
		SetAppName("go-kit"),
		SetDevelopment(true),
		SetLevel(zap.DebugLevel),
	)
}

func GetLogger() *zap.Logger {
	return logger
}

func NewLogger(mod ...ModOptions) *zap.Logger {
	l = &Logger{}
	l.Lock()
	defer l.Unlock()
	if l.inited {
		l.Info("[NewLogger] logger Inited")
		return nil
	}
	l.Opts = &Options{
		LogFileDir:    "",
		AppName:       "app_log",
		ErrorFileName: "error.log",
		WarnFileName:  "warn.log",
		InfoFileName:  "info.log",
		DebugFileName: "debug.log",
		Level:         zapcore.DebugLevel,
		MaxSize:       100,
		MaxBackups:    60,
		MaxAge:        30,
	}
	if l.Opts.LogFileDir == "" {
		l.Opts.LogFileDir, _ = filepath.Abs(filepath.Dir(filepath.Join(".")))
		l.Opts.LogFileDir += sp + "logs" + sp
	}
	if l.Opts.Development {
		l.zapConfig = zap.NewDevelopmentConfig()
		l.zapConfig.EncoderConfig.EncodeTime = timeEncoder
	} else {
		l.zapConfig = zap.NewProductionConfig()
		l.zapConfig.EncoderConfig.EncodeTime = timeUnixNano
	}

	if l.Opts.OutputPaths == nil || len(l.Opts.OutputPaths) == 0 {
		l.zapConfig.OutputPaths = []string{"stdout"}
	}
	if l.Opts.ErrorOutputPaths == nil || len(l.Opts.ErrorOutputPaths) == 0 {
		l.zapConfig.OutputPaths = []string{"stderr"}
	}
	for _, fn := range mod {
		fn(l.Opts)
	}
	l.zapConfig.Level.SetLevel(l.Opts.Level)
	l.init()
	l.inited = true
	l.Info("[NewLogger] success")
	return l.Logger
}

func (l *Logger) init() {
	l.setSyncers()
	var err error
	l.Logger, err = l.zapConfig.Build(l.cores())
	if err != nil {
		panic(err)
	}
	defer l.Logger.Sync()
}

func (l *Logger) setSyncers() {
	f := func(fN string) zapcore.WriteSyncer {
		return zapcore.AddSync(&lumberjack.Logger{
			Filename:   l.Opts.LogFileDir + sp + l.Opts.AppName + "-" + fN,
			MaxSize:    l.Opts.MaxSize,
			MaxAge:     l.Opts.MaxAge,
			MaxBackups: l.Opts.MaxBackups,
			LocalTime:  true,
			Compress:   true,
		})
	}
	errWs = f(l.Opts.ErrorFileName)
	warnWs = f(l.Opts.WarnFileName)
	infoWs = f(l.Opts.InfoFileName)
	debugWs = f(l.Opts.DebugFileName)
	return
}

func SetMaxSize(MaxSize int) ModOptions {
	return func(option *Options) {
		option.MaxSize = MaxSize
	}
}

func SetMaxBackups(MaxBackups int) ModOptions {
	return func(option *Options) {
		option.MaxBackups = MaxBackups
	}
}

func SetMaxAge(MaxAge int) ModOptions {
	return func(option *Options) {
		option.MaxAge = MaxAge
	}
}

func SetLogFileDir(logFileDir string) ModOptions {
	return func(option *Options) {
		option.LogFileDir = logFileDir
	}
}

func SetAppName(AppName string) ModOptions {
	return func(option *Options) {
		option.AppName = AppName
	}
}

func SetLevel(Level zapcore.Level) ModOptions {
	return func(option *Options) {
		option.Level = Level
	}
}

func SetErrorFileName(ErrorFileName string) ModOptions {
	return func(option *Options) {
		option.ErrorFileName = ErrorFileName
	}
}

func SetWarnFileName(WarnFileName string) ModOptions {
	return func(option *Options) {
		option.WarnFileName = WarnFileName
	}
}

func SetInfoFileName(InfoFileName string) ModOptions {
	return func(option *Options) {
		option.InfoFileName = InfoFileName
	}
}

func SetDebugFileName(DebugFileName string) ModOptions {
	return func(option *Options) {
		option.DebugFileName = DebugFileName
	}
}

func SetDevelopment(Development bool) ModOptions {
	return func(option *Options) {
		option.Development = Development
	}
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2020-10-01 15:12:02"))
}

func timeUnixNano(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendInt64(t.UnixNano() / 1e6)
}

func (l *Logger) cores() zap.Option {
	fileEncoder := zapcore.NewJSONEncoder(l.zapConfig.EncoderConfig)
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = timeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	errPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.ErrorLevel && zapcore.ErrorLevel-l.zapConfig.Level.Level() > -1
	})
	warnPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.WarnLevel && zapcore.WarnLevel-l.zapConfig.Level.Level() > -1
	})
	infoPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel && zapcore.InfoLevel-l.zapConfig.Level.Level() > -1
	})
	debugPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.DebugLevel && zapcore.DebugLevel-l.zapConfig.Level.Level() > -1
	})

	cores := []zapcore.Core{
		zapcore.NewCore(fileEncoder, errWs, errPriority),
		zapcore.NewCore(fileEncoder, warnWs, warnPriority),
		zapcore.NewCore(fileEncoder, infoWs, infoPriority),
		zapcore.NewCore(fileEncoder, debugWs, debugPriority),
	}
	if l.Opts.Development {
		cores = append(cores, []zapcore.Core{
			zapcore.NewCore(consoleEncoder, errorConsoleWS, errPriority),
			zapcore.NewCore(consoleEncoder, debugConsoleWS, warnPriority),
			zapcore.NewCore(consoleEncoder, debugConsoleWS, infoPriority),
			zapcore.NewCore(consoleEncoder, debugConsoleWS, debugPriority),
		}...)
	}
	return zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(cores...)
	})
}
