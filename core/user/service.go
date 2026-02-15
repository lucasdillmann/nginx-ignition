package user

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"

	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/common/validation"
	"dillmann.com.br/nginx-ignition/core/user/passwordhash"
)

type service struct {
	repository    Repository
	configuration *configuration.Configuration
}

func newService(repository Repository, cfg *configuration.Configuration) *service {
	return &service{
		repository:    repository,
		configuration: cfg,
	}
}

func (s *service) Authenticate(
	ctx context.Context,
	username, password, code string,
) (AuthenticationOutcome, *User, error) {
	usr, err := s.repository.FindByUsername(ctx, username)
	if err != nil {
		return AuthenticationFailed, nil, err
	}

	if usr == nil {
		return AuthenticationFailed, nil, coreerror.New(
			i18n.M(ctx, i18n.K.CoreUserInvalidCredentials),
			true,
		)
	}

	passwordMatches, err := passwordhash.
		New(s.configuration).
		Verify(password, usr.PasswordHash, usr.PasswordSalt)
	if err != nil {
		return AuthenticationFailed, nil, err
	}

	if !passwordMatches {
		return AuthenticationFailed, nil, coreerror.New(
			i18n.M(ctx, i18n.K.CoreUserInvalidCredentials),
			true,
		)
	}

	totpData := usr.TOTP
	if totpData.Validated && totpData.Secret != nil && strings.TrimSpace(*totpData.Secret) != "" {
		if strings.TrimSpace(code) == "" {
			return AuthenticationMissingTOTP, nil, nil
		}

		if !totp.Validate(code, *totpData.Secret) {
			return AuthenticationFailed, nil, nil
		}
	}

	return AuthenticationSuccessful, usr, nil
}

func (s *service) UpdatePassword(
	ctx context.Context,
	id uuid.UUID,
	currentPassword, newPassword string,
) error {
	hash := passwordhash.New(s.configuration)
	databaseState, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if databaseState == nil {
		return coreerror.New(i18n.M(ctx, i18n.K.CoreUserNotFoundById), true)
	}

	passwordMatches, err := hash.Verify(
		currentPassword,
		databaseState.PasswordHash,
		databaseState.PasswordSalt,
	)
	if err != nil {
		return err
	}

	if !passwordMatches {
		return validation.SingleFieldError(
			"currentPassword",
			i18n.M(ctx, i18n.K.CoreUserCurrentPasswordMismatch),
		)
	}

	if len(newPassword) < minimumPasswordLength {
		return validation.SingleFieldError(
			"newPassword",
			i18n.M(ctx, i18n.K.CoreUserTooShort).V("min", minimumPasswordLength),
		)
	}

	updatedHash, updatedSalt, err := hash.Hash(newPassword)
	if err != nil {
		return err
	}

	databaseState.PasswordHash = updatedHash
	databaseState.PasswordSalt = updatedSalt
	return s.repository.Save(ctx, databaseState)
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*User, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repository.DeleteByID(ctx, id)
}

func (s *service) GetCount(ctx context.Context) (int, error) {
	return s.repository.Count(ctx)
}

func (s *service) OnboardingCompleted(ctx context.Context) (bool, error) {
	count, err := s.repository.Count(ctx)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *service) Save(ctx context.Context, request *SaveRequest, currentUserID *uuid.UUID) error {
	var passwordHash, passwordSalt string
	var databaseState *User
	var err error

	databaseState, err = s.repository.FindByID(ctx, request.ID)
	if err != nil {
		return err
	}

	var totpValue TOTP
	if request.Password == nil && databaseState != nil {
		passwordHash = databaseState.PasswordHash
		passwordSalt = databaseState.PasswordSalt
		totpValue = databaseState.TOTP
	} else if request.Password != nil {
		passwordHash, passwordSalt, err = passwordhash.New(s.configuration).Hash(*request.Password)
		if err != nil {
			return err
		}
	}

	if request.RemoveTOTP {
		totpValue.Secret = nil
		totpValue.Validated = false
	}

	updatedState := &User{
		ID:           request.ID,
		Enabled:      request.Enabled,
		Name:         request.Name,
		Username:     request.Username,
		PasswordHash: passwordHash,
		PasswordSalt: passwordSalt,
		Permissions:  request.Permissions,
		TOTP:         totpValue,
	}

	if err := newValidator(s.repository).validate(
		ctx,
		updatedState,
		databaseState,
		request,
		currentUserID,
	); err != nil {
		return err
	}

	return s.repository.Save(ctx, updatedState)
}

func (s *service) GetStatus(ctx context.Context, id uuid.UUID) (bool, error) {
	return s.repository.IsEnabledByID(ctx, id)
}

func (s *service) List(
	ctx context.Context,
	pageSize, pageNumber int,
	searchTerms *string,
) (*pagination.Page[User], error) {
	return s.repository.FindPage(ctx, pageSize, pageNumber, searchTerms)
}

func (s *service) GetTOTPStatus(ctx context.Context, id uuid.UUID) (bool, error) {
	usr, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return false, err
	}

	if usr == nil {
		return false, coreerror.New(i18n.M(ctx, i18n.K.CoreUserNotFoundById), true)
	}

	return usr.TOTP.Validated && usr.TOTP.Secret != nil && *usr.TOTP.Secret != "", nil
}

func (s *service) DisableTOTP(ctx context.Context, id uuid.UUID) error {
	usr, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if usr == nil {
		return coreerror.New(i18n.M(ctx, i18n.K.CoreUserNotFoundById), true)
	}

	usr.TOTP.Secret = nil
	usr.TOTP.Validated = false
	return s.repository.Save(ctx, usr)
}

func (s *service) EnableTOTP(ctx context.Context, id uuid.UUID) (string, error) {
	usr, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return "", err
	}

	if usr == nil {
		return "", coreerror.New(i18n.M(ctx, i18n.K.CoreUserNotFoundById), true)
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "nginx-ignition",
		AccountName: usr.Username,
	})
	if err != nil {
		return "", err
	}

	usr.TOTP.Secret = new(key.Secret())
	usr.TOTP.Validated = false

	err = s.repository.Save(ctx, usr)
	if err != nil {
		return "", err
	}

	return key.URL(), nil
}

func (s *service) ActivateTOTP(ctx context.Context, id uuid.UUID, code string) (bool, error) {
	usr, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return false, err
	}

	if usr == nil {
		return false, coreerror.New(i18n.M(ctx, i18n.K.CoreUserNotFoundById), true)
	}

	if usr.TOTP.Secret == nil || strings.TrimSpace(*usr.TOTP.Secret) == "" {
		return false, coreerror.New(i18n.M(ctx, i18n.K.CoreUserTotpNotEnabled), true)
	}

	valid := totp.Validate(code, *usr.TOTP.Secret)
	if !valid {
		return false, nil
	}

	usr.TOTP.Validated = true
	return true, s.repository.Save(ctx, usr)
}

func (s *service) resetPassword(ctx context.Context, username string) (string, error) {
	user, err := s.repository.FindByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", coreerror.New(i18n.M(ctx, i18n.K.CoreUserNotFound), true)
	}

	newPassword := uuid.NewString()[:8]
	updatedHash, updatedSalt, err := passwordhash.New(s.configuration).Hash(newPassword)
	if err != nil {
		return "", err
	}

	user.PasswordHash = updatedHash
	user.PasswordSalt = updatedSalt
	return newPassword, s.repository.Save(ctx, user)
}
