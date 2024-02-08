package controllers

import (
	"fmt"
	"net/http"

	"assignment-2/database"
	"assignment-2/models"

	"github.com/gin-gonic/gin"
)

func CreateOrder(ctx *gin.Context) {
	db := database.GetDB()
	newOrder := models.Order{}

	if err := ctx.ShouldBindJSON(&newOrder); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err := db.Create(&newOrder).Error
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error_status":  "Error insert data",
			"error_message": "Something wrong when try to insert data",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"order": newOrder,
	})
}

func GetOrder(ctx *gin.Context) {
	orderID := ctx.Param("orderID")
	db := database.GetDB()

	order := models.Order{}
	err := db.Preload("Items").First(&order, "ID=?", orderID).Error
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error_status":  "Failed get data order with items",
			"error_message": fmt.Sprintf("Failed get data order with items with id %v", orderID),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"order": order,
	})
}

func GetAllOrder(ctx *gin.Context) {
	var db = database.GetDB()

	var orders []models.Order
	err := db.Preload("Items").Find(&orders).Error

	if err != nil {
		fmt.Println("Error getting data:", err.Error())
	}

	ctx.JSON(http.StatusOK, gin.H{"data": orders})
}

func DeleteOrder(c *gin.Context) {
	orderID := c.Param("orderID")
	var db = database.GetDB()

	var order models.Order
	if err := db.Preload("Items").First(&order, "id = ?", orderID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Delete associated items
	for _, item := range order.Items {
		if err := db.Delete(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete associated items"})
			return
		}
	}

	// Delete the order
	if err := db.Delete(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Order with id %v has been successfully deleted", orderID)})
}

func UpdateOrder(c *gin.Context) {
	orderID := c.Param("orderID")
	var db = database.GetDB()

	// Parse request body
	var updateData struct {
		CustomerName string `json:"customerName"`
		Items        []struct {
			ItemCode    string `json:"itemCode"`
			Description string `json:"description"`
			Quantity    int    `json:"quantity"`
		} `json:"items"`
	}

	if err := c.BindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	// Fetch the order
	var order models.Order
	if err := db.Preload("Items").First(&order, "id = ?", orderID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Update order details
	order.CustomerName = updateData.CustomerName

	// Update items
	for _, updatedItem := range updateData.Items {
		// Check if item exists in order
		var itemExists bool
		var existingItem models.Item
		for _, item := range order.Items {
			if item.ItemCode == updatedItem.ItemCode {
				itemExists = true
				existingItem = item
				break
			}
		}

		if itemExists {
			// Update existing item
			existingItem.Description = updatedItem.Description
			existingItem.Quantity = updatedItem.Quantity
			if err := db.Save(&existingItem).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item"})
				return
			}
		} else {
			// Create new item
			newItem := models.Item{
				ItemCode:    updatedItem.ItemCode,
				Description: updatedItem.Description,
				Quantity:    updatedItem.Quantity,
				OrderID:     order.ID,
			}
			if err := db.Create(&newItem).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item"})
				return
			}
		}
	}

	// Save the updated order
	if err := db.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Order with id %v has been successfully updated", orderID)})
}
