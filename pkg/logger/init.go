package logger

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	Debug      bool
	Level      string // debug/info/warn/error/dpanic/panic/fatal
	CallerSkip int
	File       FileConfig
	Hooks      []HookConfig
}

type FileConfig struct {
	Enable     bool
	Path       string
	MaxSize    int
	MaxBackups int
}

type HookConfig struct {
	Enable    bool
	Level     string // debug/info/warn/error/dpanic/panic/fatal
	Type      string // gorm
	MaxBuffer int
	MaxThread int
	Extra     map[string]string
	Options   *toml.Primitive
}

type HookHandlerFunc func(md *toml.MetaData, cfg *HookConfig) (*Hook, error)

func InitWithConfig(name string, hookHandle HookHandlerFunc) (func(), error) {
	var tomlCfg struct {
		Logger Config
	}

	tomlMD, err := toml.DecodeFile(name, &tomlCfg)
	if err != nil {
		return nil, err
	}

	cfg := tomlCfg.Logger
	var zconfig zap.Config
	if cfg.Debug {
		cfg.Level = "debug"
		zconfig = zap.NewDevelopmentConfig()
	} else {
		zconfig = zap.NewProductionConfig()
	}

	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}
	zconfig.Level.SetLevel(level)

	var (
		logger   *zap.Logger
		cleanFns []func()
	)

	if cfg.File.Enable {
		filename := cfg.File.Path
		_ = os.MkdirAll(filepath.Dir(filename), 0777)
		fileWriter := &lumberjack.Logger{
			Filename:   filename,
			MaxSize:    cfg.File.MaxSize,
			MaxBackups: cfg.File.MaxBackups,
			Compress:   false,
			LocalTime:  true,
		}

		cleanFns = append(cleanFns, func() {
			_ = fileWriter.Close()
		})

		zc := zapcore.NewCore(
			zapcore.NewJSONEncoder(zconfig.EncoderConfig),
			zapcore.AddSync(fileWriter),
			zconfig.Level,
		)
		logger = zap.New(zc)
	} else {
		ilogger, err := zconfig.Build()
		if err != nil {
			return nil, err
		}
		logger = ilogger
	}

	skip := cfg.CallerSkip
	if skip <= 0 {
		skip = 2
	}

	logger = logger.WithOptions(
		zap.WithCaller(true),
		zap.AddStacktrace(zap.ErrorLevel),
		zap.AddCallerSkip(skip),
	)

	for _, h := range cfg.Hooks {
		if !h.Enable || hookHandle == nil {
			continue
		}

		writer, err := hookHandle(&tomlMD, &h)
		if err != nil {
			return nil, err
		} else if writer == nil {
			continue
		}

		cleanFns = append(cleanFns, func() {
			writer.Flush()
		})

		hookLevel := zap.NewAtomicLevel()
		if level, err := zapcore.ParseLevel(h.Level); err == nil {
			hookLevel.SetLevel(level)
		} else {
			hookLevel.SetLevel(zap.InfoLevel)
		}

		hookEncoderConfig := zap.NewProductionEncoderConfig()
		hookEncoderConfig.EncodeTime = zapcore.EpochMillisTimeEncoder
		hookEncoderConfig.EncodeDuration = zapcore.MillisDurationEncoder
		hookCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(hookEncoderConfig),
			zapcore.AddSync(writer),
			hookLevel,
		)

		logger = logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewTee(core, hookCore)
		}))
	}

	zap.ReplaceGlobals(logger)
	return func() {
		for _, fn := range cleanFns {
			fn()
		}
	}, nil
}
