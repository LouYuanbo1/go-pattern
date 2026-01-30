package controller

import (
	"go-pattern/internal/model"
	service "go-pattern/internal/service/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (uc *UserController) RegisterRoutes(router *gin.Engine) {
	group := router.Group("/api/users")
	{
		group.POST("", uc.CreateUser)
		group.GET("/:userId", uc.GetUserByID)
		group.GET("", uc.GetUsersByPage)
		group.PATCH("/:userId", uc.UpdateUser)
		group.DELETE("/:userId", uc.DeleteUser)
	}
}

type CreateUserReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (uc *UserController) CreateUser(c *gin.Context) {
	var req CreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := uc.userService.CreateUser(c.Request.Context(), &model.User{
		Name:  req.Username,
		Email: req.Email,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"code": http.StatusCreated, "msg": "success", "data": nil})
}

func (uc *UserController) GetUserByID(c *gin.Context) {
	userId := c.Param("userId")
	uid, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := uc.userService.GetUser(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "msg": "success", "data": user})
}

type GetUsersByPageReq struct {
	Page uint64 `form:"page" binding:"required"`
	Size uint64 `form:"size" binding:"required"`
}

func (uc *UserController) GetUsersByPage(c *gin.Context) {
	var req GetUsersByPageReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	users, err := uc.userService.GetUsersByPage(c.Request.Context(), req.Page, req.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "msg": "success", "data": users})
}

type UpdateUserReq struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (uc *UserController) UpdateUser(c *gin.Context) {
	var req UpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := uc.userService.UpdateUser(c.Request.Context(), &model.User{
		ID:    req.UserID,
		Name:  req.Username,
		Email: req.Email,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "msg": "success", "data": nil})
}

func (uc *UserController) DeleteUser(c *gin.Context) {
	userId := c.Param("userId")
	uid, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := uc.userService.DeleteUser(c.Request.Context(), uid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "msg": "success", "data": nil})
}
