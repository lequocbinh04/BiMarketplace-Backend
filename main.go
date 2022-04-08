package main

import (
	"BiMarketplace/blockchain"
	"BiMarketplace/component"
	"BiMarketplace/middleware"
	"BiMarketplace/modules/config/configstorage"
	"BiMarketplace/modules/egg/eggstorage"
	"BiMarketplace/modules/pet/petstorage"
	"BiMarketplace/modules/transaction/transactionstorage"
	"BiMarketplace/modules/user/userstorage"
	"BiMarketplace/modules/user/usertransport/ginnonce"
	"BiMarketplace/modules/user/usertransport/ginuser"
	"context"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	dsn := os.Getenv("DBConnectionStr")
	secretKey := os.Getenv("SYSTEM_SECRET")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	rand.Seed(time.Now().UnixNano())
	if err != nil {
		log.Fatalln(err)
	}
	db = db.Debug()
	if err := runService(db, secretKey); err != nil {
		log.Fatalln(err)
	}
}

func runService(db *gorm.DB, secretKey string) error {
	appCtx := component.NewAppContext(db, secretKey)
	if err := blockchainListener(appCtx); err != nil {
		return err
	}
	r := gin.Default()
	r.Use(middleware.Recover(appCtx))
	r.Static("/assets", "./assets")
	return mainRoute(appCtx, r)
}

func mainRoute(appCtx *component.AppCtx, r *gin.Engine) error {
	v1 := r.Group("/v1")

	auth := v1.Group("/auth")
	{
		auth.POST("/login", ginuser.Login(appCtx))
		auth.GET("/nonce", ginnonce.GetNone(appCtx))
	}

	return r.Run()
}

func blockchainListener(appCtx *component.AppCtx) error {
	eggStore := eggstorage.NewSQLStore(appCtx.GetMainDBConnection())
	petStore := petstorage.NewSQLStore(appCtx.GetMainDBConnection())
	txStore := transactionstorage.NewSQLStore(appCtx.GetMainDBConnection())
	userStore := userstorage.NewSQLStore(appCtx.GetMainDBConnection())
	configStore := configstorage.NewSQLStore(appCtx.GetMainDBConnection())
	logCrawler := blockchain.NewLogCrawler("", 0, "", configStore)
	logHandler := blockchain.NewMkpHdl(eggStore, petStore, txStore, userStore, "./eggAbi.abi")
	if err := logCrawler.Start(); err != nil {
		return err
	}
	go logHandler.Run(context.Background(), logCrawler.GetLogChan())
	return nil
}
