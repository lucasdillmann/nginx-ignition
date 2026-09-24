package vpn

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lucasdillmann/nginx-ignition/internal/api/core/pagination"

	vpn2 "github.com/lucasdillmann/nginx-ignition/internal/business/domain/vpn"
)

type listHandler struct {
	commands vpn2.Commands
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

	pageData := pagination.Convert(page, func(vpn *vpn2.VPN) *vpnResponse {
		driver, err := h.commands.GetAvailableDriverByID(ctx, vpn.Driver)
		if err != nil {
			panic(err)
		}

		return toDTO(vpn, driver)
	})

	ctx.JSON(http.StatusOK, pageData)
}
