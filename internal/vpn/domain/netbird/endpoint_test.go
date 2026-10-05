package netbird

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/vpn"
)

func Test_endpoint(t *testing.T) {
	t.Run("start", func(t *testing.T) {
		t.Run("rejects invalid source names", func(t *testing.T) {
			for _, name := range []string{
				"../../../../../opt/nginx-ignition/frontend",
				"..",
				"1vpn",
				"vpn name",
				"vpn/name",
				"vpn.name",
			} {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				endpoint := vpn.NewMockedEndpoint(ctrl)
				endpoint.EXPECT().SourceName().Return(name).AnyTimes()

				e := &netbirdEndpoint{endpoint: endpoint}
				err := e.start(t.Context())

				assert.Error(t, err)
			}
		})
	})
}
