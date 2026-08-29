package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/lucasdillmann/nginx-ignition/internal/api/common/authorization"
	"github.com/lucasdillmann/nginx-ignition/internal/api/common/converter"
	"github.com/lucasdillmann/nginx-ignition/internal/core/user"
)

type createHandler struct {
	commands user.Commands
}

func (h createHandler) handle(ctx *gin.Context) {
	payload := &userRequestDTO{}
	if err := ctx.BindJSON(payload); err != nil {
		panic(err)
	}

	domainModel := converter.Wrap(ctx.Request.Context(), toDomain, payload)
	domainModel.ID = uuid.New()
	currentUserID := authorization.CurrentSubject(ctx).User.ID

	if err := h.commands.Save(ctx.Request.Context(), domainModel, &currentUserID); err != nil {
		panic(err)
	}

	ctx.JSON(
		http.StatusCreated,
		map[string]any{
			"id": domainModel.ID,
		},
	)
}
