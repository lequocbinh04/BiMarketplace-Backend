package appCommon

import (
	"errors"
	log "github.com/sirupsen/logrus"
)

var (
	RecordNotFound = errors.New("record not found")
)

func AppRecover() {
	if err := recover(); err != nil {
		log.Println("Recovery error:", err)
	}
}
