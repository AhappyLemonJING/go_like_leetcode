package middlewares

import (
	"gin_gorm_oj/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthAdminCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// TODO: check if user is admin
		auth := ctx.GetHeader("Authorization")
		userClaim, err := helper.AnalyseToken(auth)
		if err != nil {
			ctx.Abort()
			ctx.JSON(http.StatusOK, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Unauthorized",
			})
			return
		}
		if userClaim == nil || userClaim.IsAdmin != 1 {
			ctx.Abort()
			ctx.JSON(http.StatusOK, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Unauthorized",
			})
			return
		}
		ctx.Next()
	}
}
