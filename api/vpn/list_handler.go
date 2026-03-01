package vpn

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/api/common/pagination"
	"dillmann.com.br/nginx-ignition/core/vpn"
)

type listHandler struct {
	commands vpn.Commands
}

func (h listHandler) handle(ctx *gin.Context) {
	pageSize, pageNumber, searchTerms, err := pagination.ExtractPaginationParameters(ctx)
	if err != nil {
		panic(err)
	}

	enabledOnly := ctx.Query("enabledOnly") == "true"

	page, err := h.commands.List(
		ctx.Request.Context(),
		pageSize,
		pageNumber,
		searchTerms,
		enabledOnly,
	)
	if err != nil {
		panic(err)
	}

	pageData := pagination.Convert(page, func(vpn *vpn.VPN) *vpnResponse {
		driver, err := h.commands.GetAvailableDriverByID(ctx, vpn.Driver)
		if err != nil {
			panic(err)
		}

		return toDTO(vpn, driver)
	})

	ctx.JSON(http.StatusOK, pageData)
}
