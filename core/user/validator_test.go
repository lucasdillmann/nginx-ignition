package user

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/common/ptr"
)

func validUser() (*User, *SaveRequest) {
	id := uuid.New()
	user := &User{
		ID:       id,
		Username: "testuser",
		Name:     "Test User",
		Enabled:  true,
		Permissions: Permissions{
			Hosts:        NoAccessAccessLevel,
			Streams:      NoAccessAccessLevel,
			Certificates: NoAccessAccessLevel,
			Logs:         NoAccessAccessLevel,
			Integrations: NoAccessAccessLevel,
			AccessLists:  NoAccessAccessLevel,
			Settings:     NoAccessAccessLevel,
			Users:        NoAccessAccessLevel,
			NginxServer:  ReadOnlyAccessLevel,
			ExportData:   NoAccessAccessLevel,
			VPNs:         NoAccessAccessLevel,
			Caches:       NoAccessAccessLevel,
		},
	}
	request := &SaveRequest{
		ID:       id,
		Username: "testuser",
		Name:     "Test User",
		Enabled:  true,
		Password: ptr.Of("password123"),
		Permissions: Permissions{
			Hosts:        NoAccessAccessLevel,
			Streams:      NoAccessAccessLevel,
			Certificates: NoAccessAccessLevel,
			Logs:         NoAccessAccessLevel,
			Integrations: NoAccessAccessLevel,
			AccessLists:  NoAccessAccessLevel,
			Settings:     NoAccessAccessLevel,
			Users:        NoAccessAccessLevel,
			NginxServer:  ReadOnlyAccessLevel,
			ExportData:   NoAccessAccessLevel,
			VPNs:         NoAccessAccessLevel,
			Caches:       NoAccessAccessLevel,
		},
	}
	return user, request
}

func Test_Validator(t *testing.T) {
	ctx := context.Background()

	t.Run("validate", func(t *testing.T) {
		t.Run("valid user passes", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			user, request := validUser()
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(ctx, "testuser").Return(nil, nil)
			val := newValidator(repo)

			err := val.validate(ctx, user, nil, request, nil)

			assert.NoError(t, err)
		})

		t.Run("cannot disable own user fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			user, request := validUser()
			user.Enabled = false
			request.Enabled = false
			currentUser := &User{ID: user.ID}
			currentUserID := user.ID
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(ctx, "testuser").Return(nil, nil)
			val := newValidator(repo)

			err := val.validate(ctx, user, currentUser, request, &currentUserID)

			assert.Error(t, err)
		})

		t.Run("password missing for new user fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			user, request := validUser()
			request.Password = nil
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(ctx, "testuser").Return(nil, nil)
			val := newValidator(repo)

			err := val.validate(ctx, user, nil, request, nil)

			assert.Error(t, err)
		})

		t.Run("password can be nil for existing user", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			user, request := validUser()
			request.Password = nil
			currentUser := &User{ID: user.ID}
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(ctx, "testuser").Return(currentUser, nil)
			val := newValidator(repo)

			err := val.validate(ctx, user, currentUser, request, nil)

			assert.NoError(t, err)
		})

		t.Run("duplicate username fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			user, request := validUser()
			otherID := uuid.New()
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(ctx, "testuser").Return(&User{ID: otherID}, nil)
			val := newValidator(repo)

			err := val.validate(ctx, user, nil, request, nil)

			assert.Error(t, err)
		})

		t.Run("username too short fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			user, request := validUser()
			user.Username = "ab"
			request.Username = "ab"
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(ctx, "ab").Return(nil, nil)
			val := newValidator(repo)

			err := val.validate(ctx, user, nil, request, nil)

			assert.Error(t, err)
		})

		t.Run("name too short fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			user, request := validUser()
			user.Name = "ab"
			request.Name = "ab"
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(ctx, "testuser").Return(nil, nil)
			val := newValidator(repo)

			err := val.validate(ctx, user, nil, request, nil)

			assert.Error(t, err)
		})

		t.Run("password too short fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			user, request := validUser()
			request.Password = ptr.Of("short")
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(ctx, "testuser").Return(nil, nil)
			val := newValidator(repo)

			err := val.validate(ctx, user, nil, request, nil)

			assert.Error(t, err)
		})

		t.Run("invalid hosts access level fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			user, request := validUser()
			user.Permissions.Hosts = AccessLevel("INVALID")
			request.Permissions.Hosts = AccessLevel("INVALID")
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(ctx, "testuser").Return(nil, nil)
			val := newValidator(repo)

			err := val.validate(ctx, user, nil, request, nil)

			assert.Error(t, err)
		})

		t.Run("invalid streams access level fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			user, request := validUser()
			user.Permissions.Streams = AccessLevel("INVALID")
			request.Permissions.Streams = AccessLevel("INVALID")
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(ctx, "testuser").Return(nil, nil)
			val := newValidator(repo)

			err := val.validate(ctx, user, nil, request, nil)

			assert.Error(t, err)
		})

		t.Run("invalid certificates access level fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			user, request := validUser()
			user.Permissions.Certificates = AccessLevel("INVALID")
			request.Permissions.Certificates = AccessLevel("INVALID")
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(ctx, "testuser").Return(nil, nil)
			val := newValidator(repo)

			err := val.validate(ctx, user, nil, request, nil)

			assert.Error(t, err)
		})

		t.Run("invalid logs access level fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			user, request := validUser()
			user.Permissions.Logs = AccessLevel("INVALID")
			request.Permissions.Logs = AccessLevel("INVALID")
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(ctx, "testuser").Return(nil, nil)
			val := newValidator(repo)

			err := val.validate(ctx, user, nil, request, nil)

			assert.Error(t, err)
		})

		t.Run("invalid integrations access level fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			user, request := validUser()
			user.Permissions.Integrations = AccessLevel("INVALID")
			request.Permissions.Integrations = AccessLevel("INVALID")
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(ctx, "testuser").Return(nil, nil)
			val := newValidator(repo)

			err := val.validate(ctx, user, nil, request, nil)

			assert.Error(t, err)
		})

		t.Run("invalid access lists access level fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			user, request := validUser()
			user.Permissions.AccessLists = AccessLevel("INVALID")
			request.Permissions.AccessLists = AccessLevel("INVALID")
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(ctx, "testuser").Return(nil, nil)
			val := newValidator(repo)

			err := val.validate(ctx, user, nil, request, nil)

			assert.Error(t, err)
		})

		t.Run("invalid settings access level fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			user, request := validUser()
			user.Permissions.Settings = AccessLevel("INVALID")
			request.Permissions.Settings = AccessLevel("INVALID")
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(ctx, "testuser").Return(nil, nil)
			val := newValidator(repo)

			err := val.validate(ctx, user, nil, request, nil)

			assert.Error(t, err)
		})

		t.Run("invalid users access level fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			user, request := validUser()
			user.Permissions.Users = AccessLevel("INVALID")
			request.Permissions.Users = AccessLevel("INVALID")
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(ctx, "testuser").Return(nil, nil)
			val := newValidator(repo)

			err := val.validate(ctx, user, nil, request, nil)

			assert.Error(t, err)
		})

		t.Run("invalid nginx server access level fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			user, request := validUser()
			user.Permissions.NginxServer = AccessLevel("INVALID")
			request.Permissions.NginxServer = AccessLevel("INVALID")
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(ctx, "testuser").Return(nil, nil)
			val := newValidator(repo)

			err := val.validate(ctx, user, nil, request, nil)

			assert.Error(t, err)
		})

		t.Run("invalid export data access level fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			user, request := validUser()
			user.Permissions.ExportData = AccessLevel("INVALID")
			request.Permissions.ExportData = AccessLevel("INVALID")
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(ctx, "testuser").Return(nil, nil)
			val := newValidator(repo)

			err := val.validate(ctx, user, nil, request, nil)

			assert.Error(t, err)
		})

		t.Run("invalid VPNs access level fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			user, request := validUser()
			user.Permissions.VPNs = AccessLevel("INVALID")
			request.Permissions.VPNs = AccessLevel("INVALID")
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(ctx, "testuser").Return(nil, nil)
			val := newValidator(repo)

			err := val.validate(ctx, user, nil, request, nil)

			assert.Error(t, err)
		})

		t.Run("invalid caches access level fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			user, request := validUser()
			user.Permissions.Caches = AccessLevel("INVALID")
			request.Permissions.Caches = AccessLevel("INVALID")
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(ctx, "testuser").Return(nil, nil)
			val := newValidator(repo)

			err := val.validate(ctx, user, nil, request, nil)

			assert.Error(t, err)
		})

		t.Run("nginxServer no access fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			user, request := validUser()
			user.Permissions.NginxServer = NoAccessAccessLevel
			request.Permissions.NginxServer = NoAccessAccessLevel
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(ctx, "testuser").Return(nil, nil)
			val := newValidator(repo)

			err := val.validate(ctx, user, nil, request, nil)

			assert.Error(t, err)
		})

		t.Run("logs read-write access fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			user, request := validUser()
			user.Permissions.Logs = ReadWriteAccessLevel
			request.Permissions.Logs = ReadWriteAccessLevel
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(ctx, "testuser").Return(nil, nil)
			val := newValidator(repo)

			err := val.validate(ctx, user, nil, request, nil)

			assert.Error(t, err)
		})

		t.Run("export data read-write access fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			user, request := validUser()
			user.Permissions.ExportData = ReadWriteAccessLevel
			request.Permissions.ExportData = ReadWriteAccessLevel
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(ctx, "testuser").Return(nil, nil)
			val := newValidator(repo)

			err := val.validate(ctx, user, nil, request, nil)

			assert.Error(t, err)
		})
	})
}
