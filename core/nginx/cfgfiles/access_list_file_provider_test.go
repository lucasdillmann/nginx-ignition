package cfgfiles

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/accesslist"
	"dillmann.com.br/nginx-ignition/core/host"
)

func Test_accessListFileProvider(t *testing.T) {
	t.Run("Provide", func(t *testing.T) {
		t.Run("generate the file successfully", func(t *testing.T) {
			provider := &accessListFileProvider{}
			id := uuid.New()
			ctx := newProviderContext(t)
			ctx.hosts = []host.Host{
				{
					AccessListID: &id,
				},
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			commands := accesslist.NewMockedCommands(ctrl)
			accList := newAccessList()
			accList.ID = id
			commands.EXPECT().
				GetAll(gomock.Any()).
				Return([]accesslist.AccessList{accList}, nil)

			provider.commands = commands

			files, err := provider.provide(ctx)
			assert.NoError(t, err)
			assert.Len(t, files, 2)
			assert.Equal(t, fmt.Sprintf("access-list-%s.conf", id), files[0].Name)
			assert.Equal(t, fmt.Sprintf("access-list-%s.htpasswd", id), files[1].Name)
		})

		t.Run("returns error when commands.GetAll fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			commands := accesslist.NewMockedCommands(ctrl)
			commands.EXPECT().
				GetAll(gomock.Any()).
				Return(nil, assert.AnError)

			ctx := newProviderContext(t)
			provider := &accessListFileProvider{
				commands: commands,
			}
			_, err := provider.provide(ctx)
			assert.ErrorIs(t, err, assert.AnError)
		})
	})

	t.Run("BuildConfFile", func(t *testing.T) {
		id := uuid.New()
		paths := newPaths()
		provider := &accessListFileProvider{}

		t.Run("generates correct content for IP based access list", func(t *testing.T) {
			accessList := newAccessList()
			accessList.ID = id
			accessList.Credentials = nil
			accessList.Entries = []accesslist.Entry{
				{
					Outcome: accesslist.AllowOutcome,
					SourceAddress: []string{
						"10.0.0.1",
						"10.0.0.2",
					},
				},
			}

			file := provider.buildConfFile(&accessList, paths)
			assert.Equal(t, fmt.Sprintf("access-list-%s.conf", id), file.Name)
			assert.Contains(t, file.Contents, "allow 10.0.0.1;")
			assert.Contains(t, file.Contents, "allow 10.0.0.2;")
			assert.Contains(t, file.Contents, "deny all;")
			assert.Contains(t, file.Contents, "satisfy any;")
		})

		t.Run("generates correct content for credentials", func(t *testing.T) {
			accessList := newAccessList()
			accessList.ID = id
			accessList.Realm = "Restricted"

			file := provider.buildConfFile(&accessList, paths)
			assert.Contains(t, file.Contents, `auth_basic "Restricted";`)
			assert.Contains(
				t,
				file.Contents,
				fmt.Sprintf("auth_basic_user_file \"/etc/nginx/access-list-%s.htpasswd\";", id),
			)
		})

		t.Run("handles satisfy all mode", func(t *testing.T) {
			accessList := newAccessList()
			accessList.ID = id
			accessList.SatisfyAll = true
			accessList.Entries = []accesslist.Entry{
				{
					Outcome:       accesslist.AllowOutcome,
					SourceAddress: []string{"10.0.0.1"},
				},
			}

			file := provider.buildConfFile(&accessList, paths)
			assert.Contains(t, file.Contents, "satisfy all;")
		})

		t.Run("handles satisfy any mode when requested", func(t *testing.T) {
			accessList := newAccessList()
			accessList.ID = id
			accessList.SatisfyAll = false
			accessList.Entries = []accesslist.Entry{
				{
					Outcome:       accesslist.AllowOutcome,
					SourceAddress: []string{"10.0.0.1"},
				},
			}

			file := provider.buildConfFile(&accessList, paths)
			assert.Contains(t, file.Contents, "satisfy any;")
		})

		t.Run("removes Authorization header when forwarding is disabled", func(t *testing.T) {
			accessList := newAccessList()
			accessList.ForwardAuthenticationHeader = false
			accessList.Credentials = nil

			file := provider.buildConfFile(&accessList, paths)
			assert.Contains(t, file.Contents, `proxy_set_header Authorization "";`)
		})

		t.Run("keeps Authorization header when forwarding is enabled", func(t *testing.T) {
			accessList := newAccessList()
			accessList.ForwardAuthenticationHeader = true
			accessList.Credentials = nil

			file := provider.buildConfFile(&accessList, paths)
			assert.NotContains(t, file.Contents, `proxy_set_header Authorization "";`)
		})
	})

	t.Run("BuildHtpasswdFile", func(t *testing.T) {
		provider := &accessListFileProvider{}

		t.Run("returns nil for no credentials", func(t *testing.T) {
			accessList := newAccessList()
			accessList.Credentials = nil
			assert.Nil(t, provider.buildHtpasswdFile(&accessList))
		})

		t.Run("generates htpasswd entries", func(t *testing.T) {
			accessList := newAccessList()
			accessList.Credentials = []accesslist.Credentials{
				{
					Username: "user1",
					Password: "password1",
				},
			}
			file := provider.buildHtpasswdFile(&accessList)
			assert.NotNil(t, file)
			assert.Contains(t, file.Contents, "user1:")
		})
	})
}

func Test_toNginxOperation(t *testing.T) {
	assert.Equal(t, "allow", toNginxOperation(accesslist.AllowOutcome))
	assert.Equal(t, "deny", toNginxOperation(accesslist.DenyOutcome))
	assert.Equal(t, "", toNginxOperation("INVALID"))
}
