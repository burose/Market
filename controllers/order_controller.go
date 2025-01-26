package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"market/global"
	"market/models"
	"net/http"
)

var cacheKey_order = "order"

type UpdateorderRequest struct {
	userid    uint `json:"userid"`
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  uint `json:"quantity" binding:"required"`
}

func Createorder(ctx *gin.Context) {
	var order models.Order
	if err := ctx.ShouldBindJSON(&order); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := global.DB.AutoMigrate(&order); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := global.DB.Create(&order).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := global.RedisDB.Del(cacheKey_order).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 重新获取所有订单数据，并更新缓存
	var orders []models.Order
	if err := global.DB.Find(&orders).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// 将订单数据序列化为JSON格式
	orderjson, err := json.Marshal(orders)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// 设置新的缓存
	if err := global.RedisDB.Set(cacheKey_cart, orderjson, 0).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusOK, order)
}
func Getorder(ctx *gin.Context) {
	cachedata, err := global.RedisDB.Get(cacheKey_order).Bytes()
	if err == redis.Nil {
		var orders []models.Order
		if err := global.DB.Find(&orders).Error; err != nil {
			if errors.Is(err, redis.Nil) {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		orderjson, err := json.Marshal(orders)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := global.RedisDB.Set(cacheKey_order, orderjson, 0).Err(); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, orders)
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		var orders []models.Order
		if err := json.Unmarshal(cachedata, &orders); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, orders)
	}
}
func Updateorder(ctx *gin.Context) {
	orderid := ctx.Param("id")
	req := UpdateorderRequest{}
	var order models.Order
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := global.DB.Preload("Products").First(&order, orderid).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var product models.Product
	if err := global.DB.Preload("product").Where("product_id = ?", req.ProductID).First(&product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	order.Total += req.Quantity * product.Price
	product.Number -= req.Quantity
	if err := global.DB.Save(&product).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 重新获取所有订单数据，并更新缓存
	var orders []models.Order
	if err := global.DB.Find(&orders).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// 将订单数据序列化为JSON格式
	orderjson, err := json.Marshal(orders)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// 设置新的缓存
	if err := global.RedisDB.Set(cacheKey_cart, orderjson, 0).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, order)
}
func Cancelorder(ctx *gin.Context) {
	orderid := ctx.Param("id")
	if err := global.RedisDB.Del(cacheKey_order).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := global.DB.Where("id = ?", orderid).Delete(&models.Order{}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 重新获取所有订单数据，并更新缓存
	var orders []models.Order
	if err := global.DB.Find(&orders).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// 将订单数据序列化为JSON格式
	orderjson, err := json.Marshal(orders)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// 设置新的缓存
	if err := global.RedisDB.Set(cacheKey_cart, orderjson, 0).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Order cleared successfully"})
}
func Payorder(ctx *gin.Context) {
	id := ctx.Param("id")
	var order models.Order
	if err := global.DB.Where("order_id = ?", id).First(&order).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	order.Status = "paid"
	if err := global.DB.Save(&order).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := global.RedisDB.Del(cacheKey_order).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 重新获取所有订单数据，并更新缓存
	var orders []models.Order
	if err := global.DB.Find(&orders).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// 将订单数据序列化为JSON格式
	orderjson, err := json.Marshal(orders)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// 设置新的缓存
	if err := global.RedisDB.Set(cacheKey_cart, orderjson, 0).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Order paid successfully", "data": order})
}
