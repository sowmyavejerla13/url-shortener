package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sowmyavejerla13/url-shortener/internal/dto"
	"github.com/sowmyavejerla13/url-shortener/internal/service"
	"github.com/sowmyavejerla13/url-shortener/internal/utils"
)

type AuthHandler struct {
	userService *service.UserService
}

func NewAuthHandler(service *service.UserService)*AuthHandler{
	return &AuthHandler{
		userService: service,
	}
}

func (h *AuthHandler)Register(c *gin.Context){
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err!=nil{
		if _,ok := err.(validator.ValidationErrors);ok{
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": utils.FormatValidationErrors(err),
			})
			return 
		}
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request body",
			})
			return
	}

	err := h.userService.Register(req.Name,req.Email,req.Password)
	if err !=nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.RegisterResponse{
		Message: "User registered successfully",
	})

}

func (h *AuthHandler)Login(c *gin.Context){

	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err!=nil{
		if _,ok := err.(validator.ValidationErrors);ok{
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": utils.FormatValidationErrors(err),
			})
			return 
		}
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request body",
			})
			return
	}

	token, err:= h.userService.Login(req.Email,req.Password)
	if err!=nil{
		c.JSON(http.StatusUnauthorized,gin.H{
			"error":err.Error(),
		})
		return
	}
		c.JSON(http.StatusCreated,dto.LoginResponse{
			Token: token,
		})

}

func(h *AuthHandler)Me(c * gin.Context){
	userID := c.GetString("userID")

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"message": "Authenticated",
	})
}

func (h *AuthHandler)Urls(c *gin.Context){
	
}