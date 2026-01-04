package binding

import (
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
)

func newHTTPBinding() *Binding {
	return &Binding{
		Type: HTTPBindingType,
		IP:   "192.168.1.1",
		Port: 80,
	}
}

func newHTTPSBinding() *Binding {
	return &Binding{
		Type:          HTTPSBindingType,
		IP:            "192.168.1.1",
		Port:          443,
		CertificateID: ptr.Of(uuid.New()),
	}
}

func certCommandsExists(ctrl *gomock.Controller, certID uuid.UUID) certificate.Commands {
	m := certificate.NewMockedCommands(ctrl)
	m.EXPECT().Exists(gomock.Any(), certID).AnyTimes().Return(true, nil)
	m.EXPECT().Exists(gomock.Any(), gomock.Not(certID)).AnyTimes().Return(false, nil)
	return m
}

func certCommandsNotExists(ctrl *gomock.Controller) certificate.Commands {
	m := certificate.NewMockedCommands(ctrl)
	m.EXPECT().Exists(gomock.Any(), gomock.Any()).AnyTimes().Return(false, nil)
	return m
}
