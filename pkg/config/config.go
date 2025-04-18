package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v2"
)

type YMLConfig struct {
	LogLevel  string `yaml:"logLevel"`
	LogFile   string `yaml:"logFile"`
	Port      string `yaml:"port"`
	CacheSize string `yaml:"cache-size"`
}

type appLogger struct {
	Log *zap.Logger
}

var (
	Logger    *appLogger
	Port      string
	CacheSize = 10
	err       error
)

func (yml *YMLConfig) readConfigYMLFile() bool {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("unable to fetch current directory, error:", err)
		return false
	}

	data, err := os.ReadFile(dir + "/config.yml")
	if err != nil {
		fmt.Println("unable to read config.yml file, error:", err)
		return false
	}

	// Unmarshal YAML data into Config struct
	if err := yaml.Unmarshal(data, yml); err != nil {
		fmt.Println("unable to unmarshal config file data into struct, error:", err)
		return false
	}

	return true
}

func initLogger(logLevel, logFile string) *appLogger {

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:  "time",
		LevelKey: "level",
		// NameKey:       "logger",
		CallerKey:  "caller",
		MessageKey: "message",
		// StacktraceKey: "stacktrace",
		// LineEnding:    zapcore.DefaultLineEnding,
		EncodeTime: zapcore.ISO8601TimeEncoder,
		// EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename: logFile,
		MaxSize:  20,
	})

	// Create core with file writer and encoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), fileWriter, zap.NewAtomicLevelAt(translateLogLevel(logLevel)),
	)

	return &appLogger{Log: zap.New(core, zap.AddCaller())}

}

func translateLogLevel(logLevel string) zapcore.Level {
	if logLevel == "info" {
		return zapcore.InfoLevel
	} else if logLevel == "error" {
		return zapcore.ErrorLevel
	} else if logLevel == "debug" {
		return zapcore.DebugLevel
	} else if logLevel == "warning" {
		return zapcore.WarnLevel
	} else {
		return zapcore.PanicLevel
	}
}

func LoadCommon() bool {
	y := YMLConfig{}

	if !y.readConfigYMLFile() {
		return false
	}

	y.LogFile = strings.Replace(y.LogFile, "<<instanceId>>", y.Port, 1)

	Logger = initLogger(y.LogLevel, y.LogFile)

	Port = y.Port

	CacheSize, err = strconv.Atoi(y.CacheSize)
	if err != nil {
		Logger.Log.Sugar().Errorf("cache size is not numeric in config file, hence taking default size %d", CacheSize)
	} else {
		Logger.Log.Sugar().Infof("cache size is %d", CacheSize)
	}

	err := godotenv.Load()
	if err != nil {
		Logger.Log.Error("unable to load .env file")
	}

	Port = os.Getenv("PORT")
	if Port == "" {
		Port = "10000"
	}

	return true
}
