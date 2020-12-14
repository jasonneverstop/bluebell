package logger

import (
	"bluebell_renjiexuan/setting"

	"go.uber.org/zap"
)

var lg *zap.Logger

func Init(cfg *setting.LogConfig, mode string) (err error) {
	return
}
