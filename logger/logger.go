package logger

import (
	"bluebell_renjiexuan/setting"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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

//GinLogger 接受gin框架默认日志（作为gin的中间件）
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()
		cost := time.Since(start)
		lg.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

//GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
// TODO 研究这个函数具体写的些什么
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection,as it is not really a
				//condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					lg.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					lg.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					lg.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
