package log

import (
	"log/slog"
	"os"

	"github.com/nikitaSstepanov/tools/log/handlers"
)

type Config struct {
	Level      string `yaml:"level"       env:"LOGGER_LEVEL"       env-default:"info"`
	AddSource  bool   `yaml:"add_source"  env:"LOGGER_ADD_SOURCE"  env-default:"true"`
	IsJSON     bool   `yaml:"is_json"     env:"LOGGER_IS_JSON"     env-default:"true"`
	Out        string `yaml:"out"         env:"LOGGER_OUT"         env-default:"stdout"`
	OutPath    string `yaml:"out_path"    env:"LOGGER_OUT_PATH"    env-default:""`
	SetDefault bool   `yaml:"set_default" env:"LOGGER_SET_DEFAULT" env-default:"true"`
	Type       string `yaml:"type"        env:"LOGGER_TYPE"        env-default:"default"` 
}

func New(cfg *Config) *Logger {
	handler := setupHandler(cfg)

	logger := slog.New(handler)

	if cfg.SetDefault {
		SetDefault(logger)
	}

	return logger
}

func setupHandler(cfg *Config) Handler {
	level := setLoggerLevel(cfg.Level)

	opts := setHandlerOptions(level, cfg.AddSource)

	out := setOut(cfg)

	var handler Handler
	
	switch cfg.Type {

	case PrettyLogger:
		handler = handlers.NewPretty(out, opts)

	case DiscardLogger:
		handler = handlers.NewDiscard()
		
	default:
		if cfg.IsJSON {
			handler = NewJSONHandler(out, opts)
		} else {
			handler  = NewTextHandler(out, opts)
		}

	}

	return handler
}

func setLoggerLevel(lvl string) Level {
	var level Level

	switch lvl {

	case "debug":
		level = -4
	case "info":
		level = 0
	case "warn":
		level = 4
	case "error":
		level = 8
	default:
		level = 0

	}

	return level
}

func setHandlerOptions(level Level, AddSource bool) *HandlerOptions {
	return &HandlerOptions{AddSource: AddSource, Level: level}
}

func setOut(cfg *Config) *os.File {
	if cfg.Out == FileOut {
		return getLogFile(cfg.OutPath)
	}

	return os.Stdout
}

func getLogFile(path string) *os.File {
	if path == "" {
		path = "logs"
	}

	if err := os.RemoveAll(path); err != nil {
		panic(err)
	}

	if err := os.MkdirAll(path, 0777); err != nil {
		panic(err)
	}

	logFile, err := os.OpenFile(path + "/all.log", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}

	return logFile
}
