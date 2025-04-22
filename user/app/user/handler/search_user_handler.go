package user

import (
	"context"
	"socialmedia/user/app/user/usecase"
	"socialmedia/user/domain"

	"github.com/gofiber/fiber/v2"
)

type SearchUserHandler struct {
	usecase usecase.SearchUserUseCase
}
type SearchUserRequest struct {
	Identifier string `json:"identifier"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
}

type SearchUserResponse struct {
	Users []*domain.UserSearchResult
}

func NewSearchUserHandler(usecase usecase.SearchUserUseCase) *SearchUserHandler {
	return &SearchUserHandler{usecase: usecase}
}

func (h *SearchUserHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *SearchUserRequest) (*SearchUserResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 || req.Limit > 20 {
		req.Limit = 10
	}
	users, err := h.usecase.Execute(fbrCtx, ctx, req.Identifier, req.Page, req.Limit)
	if err != nil {
		return nil, err
	}
	return &SearchUserResponse{Users: users}, nil
}
