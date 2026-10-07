package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lucasdillmann/nginx-ignition/internal/api/core/authorization"
	"github.com/lucasdillmann/nginx-ignition/internal/api/core/converter"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/user"
)

type apiTokenCreateHandler struct {
	commands   user.Commands
	authorizer *authorization.ABAC
}

func (h apiTokenCreateHandler) handle(ctx *gin.Context) {
	currentSubject := authorization.CurrentSubject(ctx)
	if currentSubject == nil || currentSubject.User == nil {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	requestPayload := &apiTokenCreateRequestDTO{}
	if err := ctx.BindJSON(requestPayload); err != nil {
		panic(err)
	}

	request := converter.Wrap(ctx.Request.Context(), toAPITokenDomain, requestPayload)

	domainModel, err := h.commands.CreateAPIToken(
		ctx.Request.Context(),
		currentSubject.User.ID,
		request,
	)
	if err != nil {
		panic(err)
	}

	accessToken, err := h.authorizer.Jwt().GenerateToken(
		currentSubject.User,
		authorization.APIKind,
		&domainModel.ID,
		domainModel.Expiration,
	)
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusCreated, toAPITokenCreatedDTO(domainModel, *accessToken))
}
