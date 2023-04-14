package controllers

import (
	"go-jwt-challenge/database"
	"go-jwt-challenge/helpers"
	"go-jwt-challenge/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	appJSON = "application/json"
)

func UserRegister(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	_, _ = db, contentType
	User := models.User{}

	if contentType == appJSON {
		c.ShouldBindJSON(&User)
	} else {
		c.ShouldBind(&User)
	}

	err := db.Debug().Create(&User).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":        User.ID,
		"role":      User.Role,
		"full_name": User.FullName,
		"email":     User.Email,
	})
}

func UserLogin(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)

	_, _ = db, contentType

	User := models.User{}
	password := ""

	if contentType == appJSON {
		c.ShouldBindJSON(&User)
	} else {
		c.ShouldBind(&User)
	}

	password = User.Password

	err := db.Debug().Where("email = ?", User.Email).Take(&User).Error

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "invalid email/password",
		})

		return
	}

	comparePass := helpers.ComparePass([]byte(User.Password), []byte(password))

	if !comparePass {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "invalid email/password",
		})

		return
	}

	token := helpers.GenerateToken(User.ID, User.Email, User.Role)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func GetAllUsers(c *gin.Context) {
    db := database.GetDB()
    var users []models.User

    err := db.Debug().Select("id, role, full_name, email").Find(&users).Error

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   "Internal Server Error",
            "message": err.Error(),
        })

        return
    }

    c.JSON(http.StatusOK, gin.H{
        "users": users,
    })
}


func GetUserByID(c *gin.Context) {
	db := database.GetDB()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid user ID",
		})
	
		return
	}
	
	var user models.User
	
	err = db.Debug().Where("id = ?", id).Take(&user).Error
	
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "User not found",
		})
	
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"id":        user.ID,
		"role":      user.Role,
		"full_name": user.FullName,
		"email":     user.Email,
	})
}	
