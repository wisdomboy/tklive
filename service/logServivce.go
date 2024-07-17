package service

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
)

type LogService struct {
	logger *zap.Logger
}

// GetConfYaml 获取全局配置文件
func GetConfYaml() {
	viper.SetConfigFile("config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to read config file: %v", err))
	}
}

func createFile() string {
	GetConfYaml()
	filePath := viper.GetString("log.log_file_path")
	// 创建日志目录
	err := os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		panic(fmt.Errorf("failed to create log directory: %s", err))
	}
	return filepath.Join(filePath, "app.log")
}

func NewLogService() (*LogService, error) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	filename := createFile()
	logFile := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    viper.GetInt("log.max_file_size"), // 每个日志文件的最大尺寸（MB）
		MaxBackups: viper.GetInt("log.max_backups"),   // 最多保留的旧日志文件数
		MaxAge:     viper.GetInt("log.max_age"),       // 保留的旧日志文件的最大天数
		Compress:   true,                              // 是否压缩旧日志文件
	}

	fileBackend := zapcore.AddSync(logFile)
	consoleBackend := zapcore.AddSync(os.Stdout)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(fileBackend, consoleBackend),
		zap.NewAtomicLevelAt(zap.InfoLevel),
	)

	logger := zap.New(core)

	return &LogService{
		logger: logger,
	}, nil
}

func (ls *LogService) Info(message string) {
	ls.logger.Info(message)
	ls.logger.Info("")
}

func (ls *LogService) Error(message string) {
	ls.logger.Error(message)
	ls.logger.Error("")
}

func (ls *LogService) Close() error {
	err := ls.logger.Sync()
	if err != nil {
		return err
	}

	return nil
}
