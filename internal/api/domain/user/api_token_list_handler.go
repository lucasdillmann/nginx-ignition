package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lucasdillmann/nginx-ignition/internal/api/core/authorization"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/user"
)

type apiTokenListHandler struct {
	commands user.Commands
}

func (h apiTokenListHandler) handle(ctx *gin.Context) {
	currentSubject := authorization.CurrentSubject(ctx)
	if currentSubject == nil || currentSubject.User == nil {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	tokens, err := h.commands.ListAPITokens(ctx.Request.Context(), currentSubject.User.ID)
	if err != nil {
		panic(err)
	}

	responsePayload := make([]apiTokenResponseDTO, 0, len(tokens))
	for _, token := range tokens {
		responsePayload = append(responsePayload, *toAPITokenDTO(&token))
	}

	ctx.JSON(http.StatusOK, responsePayload)
}
