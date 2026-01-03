package cfgfiles

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/accesslist"
	"dillmann.com.br/nginx-ignition/core/host"
)

func Test_AccessListFileProvider_Provide(t *testing.T) {
	p := &accessListFileProvider{}
	paths := &Paths{Config: "/etc/nginx/"}
	id := uuid.New()
	ctx := &providerContext{
		context: context.Background(),
		paths:   paths,
		hosts: []host.Host{
			{
				AccessListID: &id,
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commands := accesslist.NewMockedCommands(ctrl)
	commands.EXPECT().
		GetAll(gomock.Any()).
		Return([]accesslist.AccessList{
			{
				ID:             id,
				DefaultOutcome: accesslist.DenyOutcome,
				Credentials: []accesslist.Credentials{
					{Username: "user", Password: "pwd"},
				},
			},
		}, nil)

	p.commands = commands

	files, err := p.provide(ctx)
	assert.NoError(t, err)
	assert.Len(t, files, 2)
	assert.Equal(t, fmt.Sprintf("access-list-%s.conf", id), files[0].Name)
	assert.Equal(t, fmt.Sprintf("access-list-%s.htpasswd", id), files[1].Name)

	t.Run("returns error when commands.GetAll fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		commands := accesslist.NewMockedCommands(ctrl)
		commands.EXPECT().
			GetAll(gomock.Any()).
			Return(nil, assert.AnError)

		p.commands = commands
		_, err := p.provide(ctx)
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func Test_ToNginxOperation(t *testing.T) {
	assert.Equal(t, "allow", toNginxOperation(accesslist.AllowOutcome))
	assert.Equal(t, "deny", toNginxOperation(accesslist.DenyOutcome))
	assert.Equal(t, "", toNginxOperation("INVALID"))
}

func Test_AccessListFileProvider_BuildConfFile(t *testing.T) {
	id := uuid.New()
	paths := &Paths{
		Config: "/etc/nginx/",
	}
	p := &accessListFileProvider{}

	t.Run("generates correct content for IP based access list", func(t *testing.T) {
		al := &accesslist.AccessList{
			ID:             id,
			DefaultOutcome: accesslist.DenyOutcome,
			Entries: []accesslist.Entry{
				{
					Outcome: accesslist.AllowOutcome,
					SourceAddress: []string{
						"10.0.0.1",
						"10.0.0.2",
					},
				},
			},
		}

		file := p.buildConfFile(al, paths)
		assert.Equal(t, fmt.Sprintf("access-list-%s.conf", id), file.Name)
		assert.Contains(t, file.Contents, "allow 10.0.0.1;")
		assert.Contains(t, file.Contents, "allow 10.0.0.2;")
		assert.Contains(t, file.Contents, "deny all;")
		assert.Contains(t, file.Contents, "satisfy any;")
	})

	t.Run("generates correct content for credentials", func(t *testing.T) {
		al := &accesslist.AccessList{
			ID:    id,
			Realm: "Restricted",
			Credentials: []accesslist.Credentials{
				{
					Username: "user",
					Password: "pwd",
				},
			},
		}

		file := p.buildConfFile(al, paths)
		assert.Contains(t, file.Contents, `auth_basic "Restricted";`)
		assert.Contains(
			t,
			file.Contents,
			fmt.Sprintf("auth_basic_user_file /etc/nginx/access-list-%s.htpasswd;", id),
		)
	})

	t.Run("handles satisfy all mode", func(t *testing.T) {
		al := &accesslist.AccessList{
			ID:         id,
			SatisfyAll: true,
			Credentials: []accesslist.Credentials{
				{
					Username: "user",
					Password: "pwd",
				},
			},
			Entries: []accesslist.Entry{
				{
					Outcome:       accesslist.AllowOutcome,
					SourceAddress: []string{"10.0.0.1"},
				},
			},
		}

		file := p.buildConfFile(al, paths)
		assert.Contains(t, file.Contents, "satisfy all;")
	})

	t.Run("handles satisfy any mode when requested", func(t *testing.T) {
		al := &accesslist.AccessList{
			ID:         id,
			SatisfyAll: false,
			Credentials: []accesslist.Credentials{
				{
					Username: "user",
					Password: "pwd",
				},
			},
			Entries: []accesslist.Entry{
				{
					Outcome:       accesslist.AllowOutcome,
					SourceAddress: []string{"10.0.0.1"},
				},
			},
		}

		file := p.buildConfFile(al, paths)
		assert.Contains(t, file.Contents, "satisfy any;")
	})

	t.Run("removes Authorization header when forwarding is disabled", func(t *testing.T) {
		al := &accesslist.AccessList{
			ForwardAuthenticationHeader: false,
		}

		file := p.buildConfFile(al, paths)
		assert.Contains(t, file.Contents, `proxy_set_header Authorization "";`)
	})

	t.Run("keeps Authorization header when forwarding is enabled", func(t *testing.T) {
		al := &accesslist.AccessList{
			ForwardAuthenticationHeader: true,
		}

		file := p.buildConfFile(al, paths)
		assert.NotContains(t, file.Contents, `proxy_set_header Authorization "";`)
	})
}

func Test_AccessListFileProvider_BuildHtpasswdFile(t *testing.T) {
	p := &accessListFileProvider{}

	t.Run("returns nil for no credentials", func(t *testing.T) {
		al := &accesslist.AccessList{
			Credentials: []accesslist.Credentials{},
		}
		assert.Nil(t, p.buildHtpasswdFile(al))
	})

	t.Run("generates htpasswd entries", func(t *testing.T) {
		al := &accesslist.AccessList{
			ID: uuid.New(),
			Credentials: []accesslist.Credentials{
				{
					Username: "user1",
					Password: "password1",
				},
			},
		}
		file := p.buildHtpasswdFile(al)
		assert.NotNil(t, file)
		assert.Contains(t, file.Contents, "user1:")
	})
}
