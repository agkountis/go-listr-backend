package server

import (
	"errors"
	"net/http"

	"fmt"

	"github.com/agkountis/go-listr-backend/internal/app/listr-server/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type getListsResponse struct {
	Lists []database.List `json:"lists"`
}

type createListRequest struct {
	Name string `json:"name" binding:"required"`
}

type updateListRequest struct {
	Name string `json:"name" binding:"required"`
}

type createListItemRequest struct {
	Data string `json:"data" binding:"required"`
}

type createListItemResponse struct {
	ListItemId uuid.UUID `json:"list_item_id"`
	ListId     uuid.UUID `json:"list_id"`
}

type getListItemsResponse struct {
	Items []database.ListItem `json:"items"`
}

type deleteListItemRequest struct {
	ListItemId uuid.UUID `json:"list_item_id" binding:"required"`
}

type deleteListRequest struct {
	ID uuid.UUID `json:"id" binding:"required"`
}

type deleteListResponse struct {
	ID uuid.UUID `json:"id" binding:"required"`
}

type Server struct {
	db *gorm.DB
}

func New() (*Server, error) {
	db, err := database.OpenConnection()

	if err != nil {
		return nil, err
	}

	return &Server{db: db}, nil
}

func (server *Server) Start(addr string) {
	r := server.createAndInitGinEngine()
	r.Run(addr)
}

func (server *Server) StartTLS(addr, certFilePath, keyFilePath string) {
	r := server.createAndInitGinEngine()
	r.RunTLS(addr, certFilePath, keyFilePath)
}

func (server *Server) createAndInitGinEngine() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/v1")

	v1.POST("/lists", server.createList)
	v1.POST("/lists/:id", server.createListItem)

	v1.GET("/lists", server.getLists)
	v1.GET("/lists/:id", server.getListItems)

	v1.DELETE("lists", server.deleteList)
	v1.DELETE("lists/:id", server.deleteListItem)

	v1.PATCH("lists/:id", server.updateList)

	return r
}

func (server *Server) createList(c *gin.Context) {
	db := server.db

	var createListRequest createListRequest
	if err := c.BindJSON(&createListRequest); err != nil {
		// Failed JSON decodig might not always be the users fault.
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Failed to deserialize JSON body.",
			"error":   err.Error(),
		})
		return
	}

	if createListRequest.Name == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "List name cannot be emtpy"})
		return
	}

	record := database.List{Name: createListRequest.Name}
	if err := db.Create(&record).Error; err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": fmt.Sprintf("%v", record.ID),
	})
}

func (server *Server) updateList(c *gin.Context) {
	db := server.db

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

	var updateListRequest updateListRequest
	if err := c.BindJSON(&updateListRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Malformed request body JSON",
			"error":   err.Error(),
		})
		return
	}

	err = db.Model(&database.List{ID: listId}).Update("name", updateListRequest.Name).Error

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

func (server *Server) deleteList(c *gin.Context) {
	db := server.db

	var deleteListRequest deleteListRequest
	if err := c.BindJSON(&deleteListRequest); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err := db.Where(&database.ListItem{ListID: deleteListRequest.ID}).Delete(&database.ListItem{}).Error

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = db.Delete(&database.List{ID: deleteListRequest.ID}).Error

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, &deleteListResponse{ID: deleteListRequest.ID})
}

func (server *Server) getLists(c *gin.Context) {
	db := server.db

	var lists []database.List
	err := db.Find(&lists).Error

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, &getListsResponse{Lists: lists})
}

func (server *Server) createListItem(c *gin.Context) {
	db := server.db

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

	var createListItemRequest createListItemRequest
	if err = c.BindJSON(&createListItemRequest); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	record := database.ListItem{Data: createListItemRequest.Data, ListID: listId}
	err = db.Create(&record).Error

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, &createListItemResponse{ListItemId: record.ID, ListId: record.ListID})
}

func (server *Server) getListItems(c *gin.Context) {
	db := server.db

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

	var items []database.ListItem
	err = db.Where(&database.ListItem{ListID: listId}).Find(&items).Error

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, &getListItemsResponse{Items: items})
}

func (server *Server) deleteListItem(c *gin.Context) {
	db := server.db

	var reqBody deleteListItemRequest
	if err := c.BindJSON(&reqBody); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	listId, uuidParseError := uuid.Parse(c.Params.ByName("id"))
	if uuidParseError != nil {
		c.AbortWithError(http.StatusBadRequest, uuidParseError)
		return
	}

	err := db.Delete(&database.ListItem{ID: reqBody.ListItemId, ListID: listId}).Error

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, &reqBody)
}

func listExists(db *gorm.DB, listId uuid.UUID) bool {
	err := db.Where(&database.List{ID: listId}).First(&database.List{}).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}
