package user

import (
	"context"
	"socialmedia/user/app/user/usecase"
	"socialmedia/user/domain"

	"github.com/gofiber/fiber/v2"
)

type UpdateUserHandler struct {
	usecase usecase.UpdateUserUseCase
}
type UpdateUserRequest struct {
	ID        string  `json:"id"`
	Bio       *string `json:"bio,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	BannerURL *string `json:"banner_url,omitempty"`
	Location  *string `json:"location,omitempty"`
	Website   *string `json:"website,omitempty"`
	IsPrivate *bool   `json:"is_private,omitempty"`
}

type UpdateUserResponse struct {
	Message string
}

func NewUpdateUserHandler(usecase usecase.UpdateUserUseCase) *UpdateUserHandler {
	return &UpdateUserHandler{usecase: usecase}
}

func (h *UpdateUserHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *UpdateUserRequest) (*UpdateUserResponse, error) {

	update := domain.UserUpdate{
		Bio:       req.Bio,
		AvatarURL: req.AvatarURL,
		BannerURL: req.BannerURL,
		Location:  req.Location,
		Website:   req.Website,
		IsPrivate: req.IsPrivate,
	}

	err := h.usecase.Execute(fbrCtx, ctx, update)
	if err != nil {
		return nil, err
	}
	return &UpdateUserResponse{Message: "user updated"}, nil
}
