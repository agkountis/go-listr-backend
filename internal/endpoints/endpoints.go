package endpoints

import (
	"encoding/json"
	"net/http"

	"fmt"

	"github.com/agkountis/go-listr-backend/internal/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type createListRequest struct {
	Name string
}

func CreateList(c *gin.Context) {
	db, ok := c.MustGet("db").(*gorm.DB)

	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	decoder := json.NewDecoder(c.Request.Body)

	var bodyJson createListRequest
	err := decoder.Decode(&bodyJson)

	if err != nil {
		// Failed JSON decodig might not always be the users fault.
		c.AbortWithError(http.StatusBadRequest, err)
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

type getListsResponse struct {
	Lists []model.List `json:"lists"`
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

	encoder := json.NewEncoder(c.Writer)
	err = encoder.Encode(&getListsResponse{Lists: lists})

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}
