package endpoints

import (
	"errors"
	"net/http"

	"fmt"

	"github.com/agkountis/go-listr-backend/internal/app/listr-server/contracts"
	"github.com/agkountis/go-listr-backend/internal/app/listr-server/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateList(c *gin.Context) {
	db, ok := c.MustGet("db").(*gorm.DB)

	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var bodyJson contracts.CreateListRequest
	if err := c.BindJSON(&bodyJson); err != nil {
		// Failed JSON decodig might not always be the users fault.
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Failed to deserialize JSON body.",
			"error":   err.Error(),
		})
		return
	}

	if bodyJson.Name == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "List name cannot be emtpy"})
		return
	}

	record := model.List{Name: bodyJson.Name}
	result := db.Create(&record)

	if err := result.Error; err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": fmt.Sprintf("%v", record.ID),
	})
}

func UpdateList(c *gin.Context) {
	db, ok := c.MustGet("db").(*gorm.DB)

	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	listId, err := uuid.Parse(c.Param("id"))

	if !listExists(db, listId) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("List with id `%v` does not exist.", listId),
		})
		return
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Malformed list UUID.",
			"error":   err.Error(),
		})
		return
	}

	var updateListRequest contracts.UpdateListRequest
	if err := c.BindJSON(&updateListRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Malformed request body JSON",
			"error":   err.Error(),
		})
		return
	}

	err = db.Model(&model.List{ID: listId}).Update("name", updateListRequest.Name).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"list_id":  listId,
		"new_name": updateListRequest.Name,
	})
}

func DeleteList(c *gin.Context) {
	db, ok := c.MustGet("db").(*gorm.DB)

	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var deleteListRequest contracts.DeleteListRequest
	if err := c.BindJSON(&deleteListRequest); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err := db.Where(&model.ListItem{ListID: deleteListRequest.ID}).Delete(&model.ListItem{}).Error

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = db.Delete(&model.List{ID: deleteListRequest.ID}).Error

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": deleteListRequest,
	})
}

func GetLists(c *gin.Context) {
	db, ok := c.MustGet("db").(*gorm.DB)

	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var lists []model.List
	err := db.Find(&lists).Error

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, &contracts.GetListsResponse{Lists: lists})
}

func CreateListItem(c *gin.Context) {
	db, ok := c.MustGet("db").(*gorm.DB)

	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	listIdStr := c.Params.ByName("id")

	listId, err := uuid.Parse(listIdStr)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Malformed list ID string.",
			"reason":  err.Error(),
		})
		return
	}

	if !listExists(db, listId) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("List with id '%v' does not exist", listId),
		})
		return
	}

	var createListItemRequest contracts.CreateListItemRequest
	if err = c.BindJSON(&createListItemRequest); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	record := model.ListItem{Data: createListItemRequest.Data, ListID: listId}
	err = db.Create(&record).Error

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, &record)
}

func GetListItems(c *gin.Context) {
	db, ok := c.MustGet("db").(*gorm.DB)

	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	listId, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if !listExists(db, listId) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("List with id '%v' does not exist", listId),
		})
		return
	}

	var items []model.ListItem
	err = db.Where(&model.ListItem{ListID: listId}).Find(&items).Error

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, &contracts.GetListItemsResponse{Items: items})
}

func DeleteListItem(c *gin.Context) {
	db, ok := c.MustGet("db").(*gorm.DB)

	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var reqBody contracts.DeleteListItemRequest
	if err := c.BindJSON(&reqBody); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	listId, uuidParseError := uuid.Parse(c.Params.ByName("id"))
	if uuidParseError != nil {
		c.AbortWithError(http.StatusBadRequest, uuidParseError)
		return
	}

	err := db.Delete(&model.ListItem{ID: reqBody.ListItemId, ListID: listId}).Error

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, &reqBody)
}

func listExists(db *gorm.DB, listId uuid.UUID) bool {
	err := db.Where(&model.List{ID: listId}).First(&model.List{}).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}
