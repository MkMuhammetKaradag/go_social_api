package user

import (
	"context"
	"socialmedia/user/app/user/usecase"

	"github.com/gofiber/fiber/v2"
)

type UpdateAvatarHandler struct {
	usecase usecase.UpdateAvatarUseCase
}
type UpdateAvatarRequest struct {
	AvatarURL string `json:"avatar_url"`
	// BannerURL *string `json:"banner_url,omitempty"`

}

type UpdateAvatarResponse struct {
	Message string
}

func NewUpdateAvatarHandler(usecase usecase.UpdateAvatarUseCase) *UpdateAvatarHandler {
	return &UpdateAvatarHandler{usecase: usecase}
}

func (h *UpdateAvatarHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *UpdateAvatarRequest) (*UpdateAvatarResponse, error) {

	err := h.usecase.Execute(fbrCtx, ctx, req.AvatarURL)
	if err != nil {
		return nil, err
	}
	return &UpdateAvatarResponse{Message: "user updated"}, nil
}
