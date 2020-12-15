package logger

import (
	"bluebell_renjiexuan/setting"
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

var lg *zap.Logger

func Init(cfg *setting.LogConfig, mode string) (err error) {
	writeSyncer := getLogWriter(cfg.Filename, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge)
	encoder := getEnocder()
	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		return
	}
	var core zapcore.Core
	if mode == "dev" {
		//进入开发模式，日志输出到终端
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, writeSyncer, l),
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
		)
	} else {
		core = zapcore.NewCore(encoder, writeSyncer, l)
	}

	lg = zap.New(core, zap.AddCaller())

	zap.ReplaceGlobals(lg)
	zap.L().Info("init logger sruccess")
	return
}

//编码器(如何写入日志)
func getEnocder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder         //时间格式ISO8601 UTC
	encoderConfig.TimeKey = "time"                                //输出时间的key名
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder  //将日志中日志级别的字符串转换为大写
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder //执行消耗的时间转化成浮点型的秒
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder       //以包/文件:行号 格式化调用堆栈
	return zapcore.NewJSONEncoder(encoderConfig)
}

//指定日志写到哪里去
func getLogWriter(filename string, maxSize, maxBackupS, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackupS,
		MaxAge:     maxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}
