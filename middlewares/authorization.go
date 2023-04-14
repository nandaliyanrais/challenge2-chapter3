package middlewares

import (
	"go-jwt-challenge/database"
	"go-jwt-challenge/models"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func ProductAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := database.GetDB()
		productId, err := strconv.Atoi(c.Param("productId"))
		userData := c.MustGet("userData").(jwt.MapClaims)

		userRole := userData["role"]

		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": "invalid parameter",
			})

			return
		}

		if userRole != "admin" {
			userID := uint(userData["id"].(float64))
			Product := models.Product{}

			err = db.Select("user_id").First(&Product, uint(productId)).Error

			if err != nil {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"error":   "Data Not Found",
					"message": "Data doesn't exist",
				})

				return
			}

			if Product.UserID != userID {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error":   "Unauthorized",
					"message": "You don't have permission to access this data",
				})

				return
			}
		}

		c.Next()
	}
}
