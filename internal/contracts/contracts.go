package contracts

import (
	"github.com/agkountis/go-listr-backend/internal/model"
	"github.com/google/uuid"
)

type GetListsResponse struct {
	Lists []model.List `json:"lists"`
}

type CreateListRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateListRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreateListItemRequest struct {
	Data string `json:"data" binding:"required"`
}

type CreateListItemResponse struct {
	ListItemId uuid.UUID `json:"list_item_id"`
	Data       string    `json:"data"`
	ListId     uuid.UUID `json:"list_id"`
}

type GetListItemsResponse struct {
	Items []model.ListItem `json:"items"`
}

type DeleteListItemRequest struct {
	ListItemId uuid.UUID `json:"list_item_id" binding:"required"`
}

type DeleteListRequest struct {
	ID uuid.UUID `json:"id" binding:"required"`
}
