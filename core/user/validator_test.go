package user

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/common/ptr"
)

func Test_validator(t *testing.T) {
	t.Run("validate", func(t *testing.T) {
		t.Run("valid user passes", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usr := newUser()
			request := newSaveRequest()

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(t.Context(), "testuser").Return(nil, nil)
			userValidator := newValidator(repo)

			err := userValidator.validate(t.Context(), usr, nil, request, nil)

			assert.NoError(t, err)
		})

		t.Run("cannot disable own user fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usr := newUser()
			usr.Enabled = false
			request := newSaveRequest()
			request.Enabled = false
			currentUser := &User{ID: usr.ID}
			currentUserID := usr.ID

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(t.Context(), "testuser").Return(nil, nil)
			userValidator := newValidator(repo)

			err := userValidator.validate(t.Context(), usr, currentUser, request, &currentUserID)

			assert.Error(t, err)
		})

		t.Run("password missing for new user fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usr := newUser()
			request := newSaveRequest()
			request.Password = nil

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(t.Context(), "testuser").Return(nil, nil)
			userValidator := newValidator(repo)

			err := userValidator.validate(t.Context(), usr, nil, request, nil)

			assert.Error(t, err)
		})

		t.Run("password can be nil for existing user", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usr := newUser()
			request := newSaveRequest()
			request.Password = nil
			currentUser := &User{ID: usr.ID}

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(t.Context(), "testuser").Return(currentUser, nil)
			userValidator := newValidator(repo)

			err := userValidator.validate(t.Context(), usr, currentUser, request, nil)

			assert.NoError(t, err)
		})

		t.Run("duplicate username fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usr := newUser()
			request := newSaveRequest()
			otherID := uuid.New()

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(t.Context(), "testuser").Return(&User{ID: otherID}, nil)
			userValidator := newValidator(repo)

			err := userValidator.validate(t.Context(), usr, nil, request, nil)

			assert.Error(t, err)
		})

		t.Run("username too short fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usr := newUser()
			usr.Username = "ab"
			request := newSaveRequest()
			request.Username = "ab"

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(t.Context(), "ab").Return(nil, nil)
			userValidator := newValidator(repo)

			err := userValidator.validate(t.Context(), usr, nil, request, nil)

			assert.Error(t, err)
		})

		t.Run("name too short fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usr := newUser()
			usr.Name = "ab"
			request := newSaveRequest()
			request.Name = "ab"

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(t.Context(), "testuser").Return(nil, nil)
			userValidator := newValidator(repo)

			err := userValidator.validate(t.Context(), usr, nil, request, nil)

			assert.Error(t, err)
		})

		t.Run("password too short fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usr := newUser()
			request := newSaveRequest()
			request.Password = ptr.Of("short")

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(t.Context(), "testuser").Return(nil, nil)
			userValidator := newValidator(repo)

			err := userValidator.validate(t.Context(), usr, nil, request, nil)

			assert.Error(t, err)
		})

		t.Run("invalid hosts access level fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usr := newUser()
			usr.Permissions.Hosts = "INVALID"
			request := newSaveRequest()
			request.Permissions.Hosts = "INVALID"

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(t.Context(), "testuser").Return(nil, nil)
			userValidator := newValidator(repo)

			err := userValidator.validate(t.Context(), usr, nil, request, nil)

			assert.Error(t, err)
		})

		t.Run("invalid streams access level fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usr := newUser()
			usr.Permissions.Streams = "INVALID"
			request := newSaveRequest()
			request.Permissions.Streams = "INVALID"

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(t.Context(), "testuser").Return(nil, nil)
			userValidator := newValidator(repo)

			err := userValidator.validate(t.Context(), usr, nil, request, nil)

			assert.Error(t, err)
		})

		t.Run("invalid certificates access level fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usr := newUser()
			usr.Permissions.Certificates = "INVALID"
			request := newSaveRequest()
			request.Permissions.Certificates = "INVALID"

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(t.Context(), "testuser").Return(nil, nil)
			userValidator := newValidator(repo)

			err := userValidator.validate(t.Context(), usr, nil, request, nil)

			assert.Error(t, err)
		})

		t.Run("invalid logs access level fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usr := newUser()
			usr.Permissions.Logs = "INVALID"
			request := newSaveRequest()
			request.Permissions.Logs = "INVALID"

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(t.Context(), "testuser").Return(nil, nil)
			userValidator := newValidator(repo)

			err := userValidator.validate(t.Context(), usr, nil, request, nil)

			assert.Error(t, err)
		})

		t.Run("invalid integrations access level fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usr := newUser()
			usr.Permissions.Integrations = "INVALID"
			request := newSaveRequest()
			request.Permissions.Integrations = "INVALID"

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(t.Context(), "testuser").Return(nil, nil)
			userValidator := newValidator(repo)

			err := userValidator.validate(t.Context(), usr, nil, request, nil)

			assert.Error(t, err)
		})

		t.Run("invalid access lists access level fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usr := newUser()
			usr.Permissions.AccessLists = "INVALID"
			request := newSaveRequest()
			request.Permissions.AccessLists = "INVALID"

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(t.Context(), "testuser").Return(nil, nil)
			userValidator := newValidator(repo)

			err := userValidator.validate(t.Context(), usr, nil, request, nil)

			assert.Error(t, err)
		})

		t.Run("invalid settings access level fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usr := newUser()
			usr.Permissions.Settings = "INVALID"
			request := newSaveRequest()
			request.Permissions.Settings = "INVALID"

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(t.Context(), "testuser").Return(nil, nil)
			userValidator := newValidator(repo)

			err := userValidator.validate(t.Context(), usr, nil, request, nil)

			assert.Error(t, err)
		})

		t.Run("invalid users access level fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usr := newUser()
			usr.Permissions.Users = "INVALID"
			request := newSaveRequest()
			request.Permissions.Users = "INVALID"

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(t.Context(), "testuser").Return(nil, nil)
			userValidator := newValidator(repo)

			err := userValidator.validate(t.Context(), usr, nil, request, nil)

			assert.Error(t, err)
		})

		t.Run("invalid nginx server access level fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usr := newUser()
			usr.Permissions.NginxServer = "INVALID"
			request := newSaveRequest()
			request.Permissions.NginxServer = "INVALID"

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(t.Context(), "testuser").Return(nil, nil)
			userValidator := newValidator(repo)

			err := userValidator.validate(t.Context(), usr, nil, request, nil)

			assert.Error(t, err)
		})

		t.Run("invalid export data access level fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usr := newUser()
			usr.Permissions.ExportData = "INVALID"
			request := newSaveRequest()
			request.Permissions.ExportData = "INVALID"

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(t.Context(), "testuser").Return(nil, nil)
			userValidator := newValidator(repo)

			err := userValidator.validate(t.Context(), usr, nil, request, nil)

			assert.Error(t, err)
		})

		t.Run("invalid VPNs access level fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usr := newUser()
			usr.Permissions.VPNs = "INVALID"
			request := newSaveRequest()
			request.Permissions.VPNs = "INVALID"

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(t.Context(), "testuser").Return(nil, nil)
			userValidator := newValidator(repo)

			err := userValidator.validate(t.Context(), usr, nil, request, nil)

			assert.Error(t, err)
		})

		t.Run("invalid caches access level fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usr := newUser()
			usr.Permissions.Caches = "INVALID"
			request := newSaveRequest()
			request.Permissions.Caches = "INVALID"

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(t.Context(), "testuser").Return(nil, nil)
			userValidator := newValidator(repo)

			err := userValidator.validate(t.Context(), usr, nil, request, nil)

			assert.Error(t, err)
		})

		t.Run("nginxServer no access fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usr := newUser()
			usr.Permissions.NginxServer = NoAccessAccessLevel
			request := newSaveRequest()
			request.Permissions.NginxServer = NoAccessAccessLevel

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(t.Context(), "testuser").Return(nil, nil)
			userValidator := newValidator(repo)

			err := userValidator.validate(t.Context(), usr, nil, request, nil)

			assert.Error(t, err)
		})

		t.Run("logs read-write access fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usr := newUser()
			usr.Permissions.Logs = ReadWriteAccessLevel
			request := newSaveRequest()
			request.Permissions.Logs = ReadWriteAccessLevel

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(t.Context(), "testuser").Return(nil, nil)
			userValidator := newValidator(repo)

			err := userValidator.validate(t.Context(), usr, nil, request, nil)

			assert.Error(t, err)
		})

		t.Run("export data read-write access fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usr := newUser()
			usr.Permissions.ExportData = ReadWriteAccessLevel
			request := newSaveRequest()
			request.Permissions.ExportData = ReadWriteAccessLevel

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(t.Context(), "testuser").Return(nil, nil)
			userValidator := newValidator(repo)

			err := userValidator.validate(t.Context(), usr, nil, request, nil)

			assert.Error(t, err)
		})
	})
}
