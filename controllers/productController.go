package controllers

import (
	"go-jwt-challenge/database"
	"go-jwt-challenge/helpers"
	"go-jwt-challenge/models"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func CreateProduct(c *gin.Context) {

	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)

	Product := models.Product{}
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&Product)
	} else {
		c.ShouldBind(&Product)
	}

	Product.UserID = userID

	err := db.Debug().Create(&Product).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to create product",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
		"data":    Product,
	})
}


func GetAllProducts(c *gin.Context) {

	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	userRole := userData["role"]

	var products []models.Product

	if userRole == "admin" {
		err := db.Find(&products).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to retrieve products",
				"error":   err.Error(),
			})
			return
		}
	} else if userRole == "user" {
		userID := uint(userData["id"].(float64))
		err := db.Where("user_id = ?", userID).Find(&products).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to retrieve products",
				"error":   err.Error(),
			})
			return
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
			"error":   "You don't have permission to access this data",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"product": products,
	})
}

func GetProductById(c *gin.Context) {

	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	userRole := userData["role"]
	productID, err := strconv.Atoi(c.Param("productId"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid product ID",
			"error":   err.Error(),
		})
		return
	}
	
	product := models.Product{}
	
	if userRole == "admin" {
		err = db.First(&product, productID).Error
	
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Product not found",
				"error":   err.Error(),
			})
			return
		}
	} else if userRole == "user" {
		userID := uint(userData["id"].(float64))
		err = db.Where("id = ? AND user_id = ?", productID, userID).First(&product).Error
	
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Product not found",
				"error":   err.Error(),
			})
			return
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
			"error":   "You don't have permission to access this data",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"product": product,
	})
	
}

func UpdateProduct(c *gin.Context) {

	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	Product := models.Product{}

	productId, _ := strconv.Atoi(c.Param("productId"))
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&Product)
	} else {
		c.ShouldBind(&Product)
	}

	Product.UserID = userID
	Product.ID = uint(productId)

	err := db.Model(&Product).Where("id = ?", productId).Updates(models.Product{Title: Product.Title, Description: Product.Description}).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to update product",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated successfully",
		"product":    Product,
	})
}


func DeleteProduct(c *gin.Context) {

	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	productId, err := strconv.Atoi(c.Param("productId"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid product ID",
			"error":   err.Error(),
		})
		return
	}

	Product := models.Product{}
	err = db.First(&Product, uint(productId)).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Product not found",
			"error":   err.Error(),
		})
		return
	}

	userRole := userData["role"]

	if userRole != "admin" {
		userID := uint(userData["id"].(float64))

		if Product.UserID != userID {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "You don't have permission to delete this product",
			})
			return
		}
	}

	err = db.Delete(&Product).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete product",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
	})

}
