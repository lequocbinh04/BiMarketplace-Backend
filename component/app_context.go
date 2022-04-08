package component

import "gorm.io/gorm"

type AppContext interface {
	GetMainDBConnection() *gorm.DB
	SecretKey() string
}

type AppCtx struct {
	db        *gorm.DB
	secretKey string
}

func NewAppContext(db *gorm.DB, secretKey string) *AppCtx {
	return &AppCtx{db: db, secretKey: secretKey}
}

func (ctx *AppCtx) GetMainDBConnection() *gorm.DB {
	return ctx.db
}

func (ctx *AppCtx) SecretKey() string { return ctx.secretKey }
