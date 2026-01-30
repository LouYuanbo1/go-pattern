package controller

import (
	"go-pattern/internal/model"
	service "go-pattern/internal/service/order"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	orderService service.OrderService
}

func NewOrderController(orderService service.OrderService) *OrderController {
	return &OrderController{
		orderService: orderService,
	}
}

func (oc *OrderController) RegisterRoutes(router *gin.Engine) {
	group := router.Group("/api/orders")
	{
		group.POST("", oc.CreateOrder)
		group.GET("/:orderId", oc.GetOrderByID)
		group.GET("", oc.GetOrdersByPage)
		group.PATCH("/:orderId", oc.UpdateOrder)
		group.DELETE("/:orderId", oc.DeleteOrder)
	}
}

type CreateUserReq struct {
	UserID    uint64 `json:"user_id"`
	ProductID uint64 `json:"product_id"`
}

func (oc *OrderController) CreateOrder(c *gin.Context) {
	var req CreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := oc.orderService.CreateOrder(c.Request.Context(), &model.Order{
		UserID:    req.UserID,
		ProductID: req.ProductID,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"code": http.StatusCreated, "msg": "success", "data": nil})
}

func (oc *OrderController) GetOrderByID(c *gin.Context) {
	orderId := c.Param("orderId")
	oid, err := strconv.ParseUint(orderId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	order, err := oc.orderService.GetOrder(c.Request.Context(), oid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "msg": "success", "data": order})
}

type GetUsersByPageReq struct {
	Page uint64 `form:"page" binding:"required"`
	Size uint64 `form:"size" binding:"required"`
}

func (oc *OrderController) GetOrdersByPage(c *gin.Context) {
	var req GetUsersByPageReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	orders, err := oc.orderService.GetOrdersByPage(c.Request.Context(), req.Page, req.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "msg": "success", "data": orders})
}

type UpdateUserReq struct {
	OrderID   uint64 `json:"order_id"`
	UserID    uint64 `json:"user_id"`
	ProductID uint64 `json:"product_id"`
}

func (oc *OrderController) UpdateOrder(c *gin.Context) {
	var req UpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := oc.orderService.UpdateOrder(c.Request.Context(), &model.Order{
		ID:        req.OrderID,
		UserID:    req.UserID,
		ProductID: req.ProductID,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "msg": "success", "data": nil})
}

func (oc *OrderController) DeleteOrder(c *gin.Context) {
	orderId := c.Param("orderId")
	oid, err := strconv.ParseUint(orderId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := oc.orderService.DeleteOrder(c.Request.Context(), oid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "msg": "success", "data": nil})
}
