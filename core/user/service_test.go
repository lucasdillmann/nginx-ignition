package user

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/user/passwordhash"
)

func Test_service(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		t.Run("returns user when found", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			expected := newUser()

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByID(t.Context(), expected.ID).Return(expected, nil)

			cfg := &configuration.Configuration{}
			svc, _ := newCommands(repo, cfg)
			result, err := svc.Get(t.Context(), expected.ID)

			assert.NoError(t, err)
			assert.Equal(t, expected, result)
		})

		t.Run("returns error when repository fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := uuid.New()
			expectedErr := errors.New("not found")

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByID(t.Context(), id).Return(nil, expectedErr)

			cfg := &configuration.Configuration{}
			svc, _ := newCommands(repo, cfg)
			result, err := svc.Get(t.Context(), id)

			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Equal(t, expectedErr, err)
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("deletes successfully", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := uuid.New()

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().DeleteByID(t.Context(), id).Return(nil)

			cfg := &configuration.Configuration{}
			svc, _ := newCommands(repo, cfg)
			err := svc.Delete(t.Context(), id)

			assert.NoError(t, err)
		})
	})

	t.Run("List", func(t *testing.T) {
		t.Run("returns paginated results", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			expectedPage := pagination.Of([]User{*newUser()})
			searchTerms := "test"

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindPage(t.Context(), 10, 1, &searchTerms).Return(expectedPage, nil)

			cfg := &configuration.Configuration{}
			svc, _ := newCommands(repo, cfg)
			result, err := svc.List(t.Context(), 10, 1, &searchTerms)

			assert.NoError(t, err)
			assert.Equal(t, expectedPage, result)
		})
	})

	t.Run("GetCount", func(t *testing.T) {
		t.Run("returns count", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			expectedCount := 5

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().Count(t.Context()).Return(expectedCount, nil)

			cfg := &configuration.Configuration{}
			svc, _ := newCommands(repo, cfg)
			count, err := svc.GetCount(t.Context())

			assert.NoError(t, err)
			assert.Equal(t, expectedCount, count)
		})
	})

	t.Run("OnboardingCompleted", func(t *testing.T) {
		t.Run("returns true when users exist", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().Count(t.Context()).Return(1, nil)

			cfg := &configuration.Configuration{}
			svc, _ := newCommands(repo, cfg)
			completed, err := svc.OnboardingCompleted(t.Context())

			assert.NoError(t, err)
			assert.True(t, completed)
		})

		t.Run("returns false when no users exist", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().Count(t.Context()).Return(0, nil)

			cfg := &configuration.Configuration{}
			svc, _ := newCommands(repo, cfg)
			completed, err := svc.OnboardingCompleted(t.Context())

			assert.NoError(t, err)
			assert.False(t, completed)
		})
	})

	t.Run("Save", func(t *testing.T) {
		t.Run("saves user successfully", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			request := newSaveRequest()
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByID(t.Context(), request.ID).Return(nil, nil)
			repo.EXPECT().FindByUsername(t.Context(), request.Username).Return(nil, nil)
			repo.EXPECT().Save(t.Context(), gomock.Any()).Return(nil)

			cfg := &configuration.Configuration{}
			svc, _ := newCommands(repo, cfg)
			err := svc.Save(t.Context(), request, nil)

			assert.NoError(t, err)
		})

		t.Run("removes TOTP when requested", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usr := newUser()
			secret := "secret"
			usr.TOTP = TOTP{Secret: &secret, Validated: true}

			request := newSaveRequest()
			request.ID = usr.ID
			request.Password = nil
			request.RemoveTOTP = true

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByID(t.Context(), request.ID).Return(usr, nil)
			repo.EXPECT().FindByUsername(t.Context(), request.Username).Return(usr, nil)
			repo.EXPECT().Save(t.Context(), gomock.Any()).DoAndReturn(func(_ any, u *User) error {
				assert.Nil(t, u.TOTP.Secret)
				assert.False(t, u.TOTP.Validated)
				return nil
			})

			cfg := &configuration.Configuration{}
			svc, _ := newCommands(repo, cfg)
			err := svc.Save(t.Context(), request, nil)

			assert.NoError(t, err)
		})
	})

	t.Run("Authenticate", func(t *testing.T) {
		cfg := &configuration.Configuration{}
		ph := passwordhash.New(cfg)
		password := "password"
		hash, salt, _ := ph.Hash(password)

		t.Run("returns error when user not found", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(t.Context(), "nonexistent").Return(nil, nil)

			svc, _ := newCommands(repo, cfg)
			outcome, result, err := svc.Authenticate(t.Context(), "nonexistent", "password", "")

			require.Error(t, err)
			assert.Nil(t, result)
			assert.Equal(t, outcome, AuthenticationFailed)

			var coreErr *coreerror.CoreError
			require.ErrorAs(t, err, &coreErr)
			assert.Equal(t, i18n.K.CoreUserInvalidCredentials, coreErr.Message.Key)
		})

		t.Run("returns success when password matches and TOTP disabled", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usr := newUser()
			usr.PasswordHash = hash
			usr.PasswordSalt = salt
			usr.TOTP.Validated = false

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(t.Context(), usr.Username).Return(usr, nil)

			svc, _ := newCommands(repo, cfg)
			outcome, result, err := svc.Authenticate(t.Context(), usr.Username, password, "")

			assert.NoError(t, err)
			assert.Equal(t, usr, result)
			assert.Equal(t, AuthenticationSuccessful, outcome)
		})

		t.Run("returns failure when password does not match", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usr := newUser()
			usr.PasswordHash = hash
			usr.PasswordSalt = salt

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(t.Context(), usr.Username).Return(usr, nil)

			svc, _ := newCommands(repo, cfg)
			outcome, result, err := svc.Authenticate(t.Context(), usr.Username, "wrongpassword", "")

			require.Error(t, err)
			assert.Nil(t, result)
			assert.Equal(t, outcome, AuthenticationFailed)

			var coreErr *coreerror.CoreError
			require.ErrorAs(t, err, &coreErr)
			assert.Equal(t, i18n.K.CoreUserInvalidCredentials, coreErr.Message.Key)
		})

		t.Run("returns missing TOTP when enabled but code empty", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			secret := "JBSWY3DPEHPK3PXP"
			usr := newUser()
			usr.PasswordHash = hash
			usr.PasswordSalt = salt
			usr.TOTP = TOTP{Secret: &secret, Validated: true}

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(t.Context(), usr.Username).Return(usr, nil)

			svc, _ := newCommands(repo, cfg)
			outcome, result, err := svc.Authenticate(t.Context(), usr.Username, password, "")

			assert.NoError(t, err)
			assert.Nil(t, result)
			assert.Equal(t, AuthenticationMissingTOTP, outcome)
		})

		t.Run("returns failure when TOTP enabled but code invalid", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			secret := "JBSWY3DPEHPK3PXP"
			usr := newUser()
			usr.PasswordHash = hash
			usr.PasswordSalt = salt
			usr.TOTP = TOTP{Secret: &secret, Validated: true}

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(t.Context(), usr.Username).Return(usr, nil)

			svc, _ := newCommands(repo, cfg)
			outcome, result, err := svc.Authenticate(t.Context(), usr.Username, password, "000000")

			assert.NoError(t, err)
			assert.Nil(t, result)
			assert.Equal(t, AuthenticationFailed, outcome)
		})

		t.Run("returns success when TOTP enabled and code valid", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			secret := "JBSWY3DPEHPK3PXP"
			code, _ := totp.GenerateCode(secret, time.Now())

			usr := newUser()
			usr.PasswordHash = hash
			usr.PasswordSalt = salt
			usr.TOTP = TOTP{Secret: &secret, Validated: true}

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(t.Context(), usr.Username).Return(usr, nil)

			svc, _ := newCommands(repo, cfg)
			outcome, result, err := svc.Authenticate(t.Context(), usr.Username, password, code)

			assert.NoError(t, err)
			assert.Equal(t, usr, result)
			assert.Equal(t, AuthenticationSuccessful, outcome)
		})
	})

	t.Run("GetStatus", func(t *testing.T) {
		t.Run("returns true when enabled", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := uuid.New()

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().IsEnabledByID(t.Context(), id).Return(true, nil)

			cfg := &configuration.Configuration{}
			svc, _ := newCommands(repo, cfg)
			enabled, err := svc.GetStatus(t.Context(), id)

			assert.NoError(t, err)
			assert.True(t, enabled)
		})
	})

	t.Run("GetTOTPStatus", func(t *testing.T) {
		t.Run("returns true when validated and secret exists", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			secret := "JBSWY3DPEHPK3PXP"
			usr := newUser()
			usr.TOTP = TOTP{Secret: &secret, Validated: true}

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByID(t.Context(), usr.ID).Return(usr, nil)

			cfg := &configuration.Configuration{}
			svc, _ := newCommands(repo, cfg)
			status, err := svc.GetTOTPStatus(t.Context(), usr.ID)

			assert.NoError(t, err)
			assert.True(t, status)
		})

		t.Run("returns false when not validated", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			secret := "JBSWY3DPEHPK3PXP"
			usr := newUser()
			usr.TOTP = TOTP{Secret: &secret, Validated: false}

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByID(t.Context(), usr.ID).Return(usr, nil)

			cfg := &configuration.Configuration{}
			svc, _ := newCommands(repo, cfg)
			status, err := svc.GetTOTPStatus(t.Context(), usr.ID)

			assert.NoError(t, err)
			assert.False(t, status)
		})

		t.Run("returns false when secret is missing", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usr := newUser()
			usr.TOTP = TOTP{Secret: nil, Validated: true}

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByID(t.Context(), usr.ID).Return(usr, nil)

			cfg := &configuration.Configuration{}
			svc, _ := newCommands(repo, cfg)
			status, err := svc.GetTOTPStatus(t.Context(), usr.ID)

			assert.NoError(t, err)
			assert.False(t, status)
		})
	})

	t.Run("DisableTOTP", func(t *testing.T) {
		t.Run("disables successfully", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			secret := "JBSWY3DPEHPK3PXP"
			usr := newUser()
			usr.TOTP = TOTP{Secret: &secret, Validated: true}

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByID(t.Context(), usr.ID).Return(usr, nil)
			repo.EXPECT().Save(t.Context(), gomock.Any()).DoAndReturn(func(_ any, u *User) error {
				assert.Nil(t, u.TOTP.Secret)
				assert.False(t, u.TOTP.Validated)
				return nil
			})

			cfg := &configuration.Configuration{}
			svc, _ := newCommands(repo, cfg)
			err := svc.DisableTOTP(t.Context(), usr.ID)

			assert.NoError(t, err)
		})
	})

	t.Run("EnableTOTP", func(t *testing.T) {
		t.Run("enables and returns provisioning URL", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usr := newUser()

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByID(t.Context(), usr.ID).Return(usr, nil)
			repo.EXPECT().Save(t.Context(), gomock.Any()).DoAndReturn(func(_ any, u *User) error {
				assert.NotNil(t, u.TOTP.Secret)
				assert.False(t, u.TOTP.Validated)
				return nil
			})

			cfg := &configuration.Configuration{}
			svc, _ := newCommands(repo, cfg)
			url, err := svc.EnableTOTP(t.Context(), usr.ID)

			assert.NoError(t, err)
			assert.Contains(t, url, "otpauth://totp/nginx-ignition")
			assert.Contains(t, url, usr.Username)
		})
	})

	t.Run("ActivateTOTP", func(t *testing.T) {
		t.Run("activates with valid code", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Use a fixed secret and generate a valid code for it
			secret := "JBSWY3DPEHPK3PXP"
			code, _ := totp.GenerateCode(secret, time.Now())

			usr := newUser()
			usr.TOTP = TOTP{Secret: &secret, Validated: false}

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByID(t.Context(), usr.ID).Return(usr, nil)
			repo.EXPECT().Save(t.Context(), gomock.Any()).DoAndReturn(func(_ any, u *User) error {
				assert.True(t, u.TOTP.Validated)
				return nil
			})

			cfg := &configuration.Configuration{}
			svc, _ := newCommands(repo, cfg)
			ok, err := svc.ActivateTOTP(t.Context(), usr.ID, code)

			assert.NoError(t, err)
			assert.True(t, ok)
		})

		t.Run("returns false with invalid code", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			secret := "JBSWY3DPEHPK3PXP"
			usr := newUser()
			usr.TOTP = TOTP{Secret: &secret, Validated: false}

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByID(t.Context(), usr.ID).Return(usr, nil)

			cfg := &configuration.Configuration{}
			svc, _ := newCommands(repo, cfg)
			ok, err := svc.ActivateTOTP(t.Context(), usr.ID, "000000")

			assert.NoError(t, err)
			assert.False(t, ok)
		})
	})
}
