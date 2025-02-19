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

var cacheKey_cart = "cart"

type UpdatecartRequest struct {
	Userid    uint `json:"user_id" binding:"required"`
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  uint `json:"quantity" binding:"required"`
}

func Setcart(ctx *gin.Context) {
	var cart models.Cart

	if err := ctx.ShouldBindJSON(&cart); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := global.DB.AutoMigrate(&cart); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := global.DB.Create(&cart).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	var product models.Product
	if err := global.DB.Where("product_id = ?", cart.ProductID).Find(&product).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	//对应商品数量减少
	product.Number -= cart.Quantity

	if err := global.RedisDB.Del(cacheKey_cart).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// 重新获取所有购物车数据，并更新缓存
	var carts []models.Cart
	if err := global.DB.Find(&carts).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// 将购物车数据序列化为JSON格式
	cartjson, err := json.Marshal(carts)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// 设置新的缓存
	if err := global.RedisDB.Set(cacheKey_cart, cartjson, 0).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusOK, cart)
}

func Getcart(ctx *gin.Context) {
	cacheddata, err := global.RedisDB.Get(cacheKey_cart).Bytes()
	if err == redis.Nil {
		var carts []models.Cart
		if err := global.DB.Find(&carts).Error; err != nil {
			if errors.Is(err, redis.Nil) {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "cart not found"})
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			}
			return
		}
		cartjson, err := json.Marshal(carts)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		if err := global.RedisDB.Set(cacheKey_cart, cartjson, 0).Err(); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		ctx.JSON(http.StatusOK, carts)
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	} else {
		var carts []models.Cart
		if err := json.Unmarshal(cacheddata, &carts); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		ctx.JSON(http.StatusOK, carts)
	}
}
func Addcart(ctx *gin.Context) {
	cartid := ctx.Param("id")
	var req UpdatecartRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := global.DB.Preload("User").Where("id = ?", req.Userid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "cart not found"})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
	}

	var product models.Product
	if err := global.DB.Where("product_id = ?", req.ProductID).First(&product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	product.Number -= req.Quantity

	var cart models.Cart
	if err := global.DB.Where("cart_id = ?", cartid).First(&cart).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	cart.Quantity += req.Quantity

	if err := global.DB.Save(&cart).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if err := global.DB.Save(&product).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if err := global.RedisDB.Del(cacheKey_cart).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// 重新获取所有购物车数据，并更新缓存
	var carts []models.Cart
	if err := global.DB.Where("user_id = ?", req.Userid).Find(&carts).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// 将购物车数据序列化为JSON格式
	cartjson, err := json.Marshal(carts)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// 设置新的缓存
	if err := global.RedisDB.Set(cacheKey_cart, cartjson, 0).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusOK, cart)
}
func Clearcart(ctx *gin.Context) {
	cartid := ctx.Param("id")

	// 删除 Redis 中的缓存
	if err := global.RedisDB.Del(cacheKey_cart, cartid).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// 批量更新数据库中的购物车项
	if err := global.DB.Model(&models.Cart{}).Where("cart_id = ?", cartid).Update("quantity", 0).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Cart cleared successfully"})
}
func Deletecart(ctx *gin.Context) {
	cartid := ctx.Param("id")

	// 从数据库中删除购物车项
	if err := global.DB.Delete(&models.Cart{}, "cart_id = ?", cartid).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete cart from database"})
		return
	}

	// 从缓存中删除购物车项
	if err := global.RedisDB.Del(cacheKey_cart + cartid).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete cart from cache"})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{"message": "Cart item deleted successfully"})
}
