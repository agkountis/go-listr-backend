package contracts

import "github.com/agkountis/go-listr-backend/internal/model"

type GetListsResponse struct {
	Lists []model.List `json:"lists"`
}

type CreateListRequest struct {
	Name string
}

type GetListItemsResponse struct {
	Items []model.Item `json:"items"`
}
