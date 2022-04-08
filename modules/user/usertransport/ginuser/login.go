package ginuser

import (
	"BiMarketplace/appCommon"
	"BiMarketplace/component"
	"BiMarketplace/component/tokenprovider/jwt"
	"BiMarketplace/modules/user/userbiz"
	"BiMarketplace/modules/user/usermodel"
	"BiMarketplace/modules/user/userstorage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginUserData usermodel.UserLogin
		if err := c.ShouldBind(&loginUserData); err != nil {
			panic(appCommon.ErrInvalidRequest(err))
		}

		db := appCtx.GetMainDBConnection()
		tokenProvider := jwt.NewTokenJWTProvider(appCtx.SecretKey())
		store := userstorage.NewSQLStore(db)
		biz := userbiz.NewLoginBiz(store, tokenProvider, 60*60*24*3000)
		accessToken, err := biz.Login(c.Request.Context(), loginUserData)
		if err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, appCommon.SimpleSuccessResponse(accessToken))
	}
}
