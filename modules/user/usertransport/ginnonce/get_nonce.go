package ginnonce

import (
	"BiMarketplace/appCommon"
	"BiMarketplace/component"
	"BiMarketplace/modules/user/userbiz"
	"BiMarketplace/modules/user/userstorage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetNone(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		address := c.Query("address")
		store := userstorage.NewSQLStore(appCtx.GetMainDBConnection())
		biz := userbiz.NewGetNonceBiz(store)
		nonce, err := biz.GetNonce(address)
		if err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, appCommon.SimpleSuccessResponse(nonce))
	}
}
