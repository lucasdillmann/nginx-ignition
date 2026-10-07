package user

import (
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lucasdillmann/nginx-ignition/internal/database/core/database"
	"github.com/lucasdillmann/nginx-ignition/internal/database/core/testutils"
)

func Test_Repository(t *testing.T) {
	testutils.RunWithMockedDatabases(t, runRepositoryTests)
}

func Test_Repository_TryCreateInitialUser(t *testing.T) {
	t.Run("successfully creates the first user", func(t *testing.T) {
		testutils.RunWithMockedDatabases(t, func(t *testing.T, db *database.Database) {
			repo := New(db)
			cmd := newUser()

			created, err := repo.TryCreateInitialUser(t.Context(), cmd)
			require.NoError(t, err)
			assert.True(t, created)

			count, err := repo.Count(t.Context())
			require.NoError(t, err)
			assert.Equal(t, 1, count)

			saved, err := repo.FindByID(t.Context(), cmd.ID)
			require.NoError(t, err)
			require.NotNil(t, saved)
			assert.Equal(t, cmd.Username, saved.Username)
		})
	})

	t.Run("returns false when onboarding already completed", func(t *testing.T) {
		testutils.RunWithMockedDatabases(t, func(t *testing.T, db *database.Database) {
			repo := New(db)
			existing := newUser()
			created, err := repo.TryCreateInitialUser(t.Context(), existing)
			require.NoError(t, err)
			require.True(t, created)

			another := newUser()
			another.ID = uuid.New()
			another.Username = uuid.New().String()

			created, err = repo.TryCreateInitialUser(t.Context(), another)
			require.NoError(t, err)
			assert.False(t, created)

			count, err := repo.Count(t.Context())
			require.NoError(t, err)
			assert.Equal(t, 1, count)
		})
	})

	t.Run("allows only one winner under concurrent creation", func(t *testing.T) {
		testutils.RunWithMockedDatabases(t, func(t *testing.T, db *database.Database) {
			repo := New(db)
			const goroutines = 10
			var successCount int32
			var failureCount int32
			waitGroup := sync.WaitGroup{}
			waitGroup.Add(goroutines)

			for range goroutines {
				go func() {
					defer waitGroup.Done()

					candidate := newUser()
					candidate.ID = uuid.New()
					candidate.Username = uuid.New().String()

					created, err := repo.TryCreateInitialUser(t.Context(), candidate)
					if err != nil {
						if !strings.Contains(err.Error(), "SQLITE_BUSY") {
							t.Errorf("TryCreateInitialUser returned error: %v", err)
							return
						}
					}

					if created {
						atomic.AddInt32(&successCount, 1)
					} else {
						atomic.AddInt32(&failureCount, 1)
					}
				}()
			}

			waitGroup.Wait()

			assert.Equal(t, int32(1), successCount)
			assert.Equal(t, int32(goroutines-1), failureCount)

			count, err := repo.Count(t.Context())
			require.NoError(t, err)
			assert.Equal(t, 1, count)
		})
	})
}

func runRepositoryTests(t *testing.T, db *database.Database) {
	repo := New(db)

	t.Run("Save", func(t *testing.T) {
		t.Run("successfully saves a new user", func(t *testing.T) {
			cmd := newUser()

			err := repo.Save(t.Context(), cmd)
			require.NoError(t, err)

			saved, err := repo.FindByID(t.Context(), cmd.ID)
			require.NoError(t, err)
			require.NotNil(t, saved)
			assert.Equal(t, cmd.Name, saved.Name)
			assert.Equal(t, cmd.Username, saved.Username)
			assert.Equal(t, cmd.PasswordHash, saved.PasswordHash)
			assert.Equal(t, cmd.PasswordSalt, saved.PasswordSalt)
			assert.Equal(t, cmd.Permissions, saved.Permissions)
			assert.Equal(t, cmd.Enabled, saved.Enabled)
			assert.Equal(t, cmd.TOTP, saved.TOTP)
		})

		t.Run("successfully updates an existing user", func(t *testing.T) {
			id := uuid.New()
			cmd := newUser()
			cmd.ID = id
			require.NoError(t, repo.Save(t.Context(), cmd))

			cmd.Name = "Updated User"
			cmd.Enabled = false
			err := repo.Save(t.Context(), cmd)
			require.NoError(t, err)

			saved, err := repo.FindByID(t.Context(), id)
			require.NoError(t, err)
			assert.Equal(t, "Updated User", saved.Name)
			assert.False(t, saved.Enabled)
		})
	})

	t.Run("FindByUsername", func(t *testing.T) {
		t.Run("returns user by exact username", func(t *testing.T) {
			cmd := newUser()
			require.NoError(t, repo.Save(t.Context(), cmd))

			saved, err := repo.FindByUsername(t.Context(), cmd.Username)
			require.NoError(t, err)
			require.NotNil(t, saved)
			assert.Equal(t, cmd.ID, saved.ID)
		})

		t.Run("returns nil if not found", func(t *testing.T) {
			saved, err := repo.FindByUsername(t.Context(), "nonexistent")
			require.NoError(t, err)
			assert.Nil(t, saved)
		})
	})

	t.Run("FindPage", func(t *testing.T) {
		t.Run("returns a page of users filtered by name or username", func(t *testing.T) {
			prefix := uuid.New().String()
			names := []string{
				prefix + "Alpha",
				prefix + "Beta",
				prefix + "Gamma",
			}

			for _, name := range names {
				cmd := newUser()
				cmd.ID = uuid.New()
				cmd.Name = name
				cmd.Username = uuid.New().String()
				require.NoError(t, repo.Save(t.Context(), cmd))
			}

			byUsername := newUser()
			byUsername.ID = uuid.New()
			byUsername.Name = "Hidden Name"
			byUsername.Username = prefix + "User"
			require.NoError(t, repo.Save(t.Context(), byUsername))

			other := newUser()
			other.ID = uuid.New()
			other.Name = "Other" + uuid.New().String()
			other.Username = "Other" + uuid.New().String()
			require.NoError(t, repo.Save(t.Context(), other))

			page, err := repo.FindPage(t.Context(), 10, 0, new(prefix))
			require.NoError(t, err)

			assert.GreaterOrEqual(t, page.TotalItems, 4)

			for _, item := range page.Contents {
				found := false
				if strings.Contains(item.Name, prefix) || strings.Contains(item.Username, prefix) {
					found = true
				}
				assert.True(
					t,
					found,
					"Item %s (user: %s) should match %s",
					item.Name,
					item.Username,
					prefix,
				)
			}
		})
	})

	t.Run("IsEnabledByID", func(t *testing.T) {
		t.Run("returns true when enabled", func(t *testing.T) {
			cmd := newUser()
			cmd.Enabled = true
			require.NoError(t, repo.Save(t.Context(), cmd))

			enabled, err := repo.IsEnabledByID(t.Context(), cmd.ID)
			require.NoError(t, err)
			assert.True(t, enabled)
		})

		t.Run("returns false when disabled", func(t *testing.T) {
			cmd := newUser()
			cmd.Enabled = false
			require.NoError(t, repo.Save(t.Context(), cmd))

			enabled, err := repo.IsEnabledByID(t.Context(), cmd.ID)
			require.NoError(t, err)
			assert.False(t, enabled)
		})

		t.Run("returns false when not exists", func(t *testing.T) {
			enabled, err := repo.IsEnabledByID(t.Context(), uuid.New())
			require.NoError(t, err)
			assert.False(t, enabled)
		})
	})

	t.Run("Count", func(t *testing.T) {
		t.Run("returns total user count", func(t *testing.T) {
			initial, err := repo.Count(t.Context())
			require.NoError(t, err)

			cmd := newUser()
			require.NoError(t, repo.Save(t.Context(), cmd))

			cmd2 := newUser()
			cmd2.ID = uuid.New()
			cmd2.Username = uuid.New().String()
			require.NoError(t, repo.Save(t.Context(), cmd2))

			current, err := repo.Count(t.Context())
			require.NoError(t, err)
			assert.Equal(t, initial+2, current)
		})
	})

	t.Run("DeleteByID", func(t *testing.T) {
		t.Run("removes the user", func(t *testing.T) {
			cmd := newUser()
			require.NoError(t, repo.Save(t.Context(), cmd))

			err := repo.DeleteByID(t.Context(), cmd.ID)
			require.NoError(t, err)

			saved, err := repo.FindByID(t.Context(), cmd.ID)
			require.NoError(t, err)
			assert.Nil(t, saved)
		})

		t.Run("also removes the user tokens", func(t *testing.T) {
			cmd := newUser()
			require.NoError(t, repo.Save(t.Context(), cmd))
			token := newAPIToken(cmd)
			require.NoError(t, repo.CreateToken(t.Context(), token))

			require.NoError(t, repo.DeleteByID(t.Context(), cmd.ID))

			saved, err := repo.FindTokenByID(t.Context(), cmd.ID, token.ID)
			require.NoError(t, err)
			assert.Nil(t, saved)
		})
	})

	t.Run("CreateToken", func(t *testing.T) {
		t.Run("persists the token attributes", func(t *testing.T) {
			usr := newUser()
			require.NoError(t, repo.Save(t.Context(), usr))
			token := newAPIToken(usr)

			err := repo.CreateToken(t.Context(), token)
			require.NoError(t, err)

			saved, err := repo.FindTokenByID(t.Context(), usr.ID, token.ID)
			require.NoError(t, err)
			require.NotNil(t, saved)
			assert.Equal(t, token.ID, saved.ID)
			assert.Equal(t, usr.ID, saved.UserID)
			assert.Equal(t, token.Name, saved.Name)
			require.NotNil(t, saved.Expiration)
			assert.True(t, saved.CreatedAt.Equal(token.CreatedAt))
		})

		t.Run("supports tokens without expiration", func(t *testing.T) {
			usr := newUser()
			require.NoError(t, repo.Save(t.Context(), usr))
			token := newAPIToken(usr)
			token.Expiration = nil

			require.NoError(t, repo.CreateToken(t.Context(), token))

			saved, err := repo.FindTokenByID(t.Context(), usr.ID, token.ID)
			require.NoError(t, err)
			require.NotNil(t, saved)
			assert.Nil(t, saved.Expiration)
		})

		t.Run("rejects duplicated names for the same user", func(t *testing.T) {
			usr := newUser()
			require.NoError(t, repo.Save(t.Context(), usr))
			token := newAPIToken(usr)
			require.NoError(t, repo.CreateToken(t.Context(), token))

			duplicated := newAPIToken(usr)
			duplicated.Name = token.Name

			err := repo.CreateToken(t.Context(), duplicated)
			assert.Error(t, err)
		})

		t.Run("allows the same name for different users", func(t *testing.T) {
			firstUser := newUser()
			require.NoError(t, repo.Save(t.Context(), firstUser))
			secondUser := newUser()
			require.NoError(t, repo.Save(t.Context(), secondUser))

			token := newAPIToken(firstUser)
			require.NoError(t, repo.CreateToken(t.Context(), token))

			another := newAPIToken(secondUser)
			another.Name = token.Name

			require.NoError(t, repo.CreateToken(t.Context(), another))
		})
	})

	t.Run("FindTokensByUserID", func(t *testing.T) {
		t.Run("returns only the tokens owned by the user, ordered by name", func(t *testing.T) {
			usr := newUser()
			require.NoError(t, repo.Save(t.Context(), usr))
			other := newUser()
			require.NoError(t, repo.Save(t.Context(), other))
			require.NoError(t, repo.CreateToken(t.Context(), newAPIToken(other)))

			second := newAPIToken(usr)
			second.Name = "zeta-token"
			require.NoError(t, repo.CreateToken(t.Context(), second))

			first := newAPIToken(usr)
			first.Name = "alpha-token"
			require.NoError(t, repo.CreateToken(t.Context(), first))

			page, err := repo.FindTokensByUserID(t.Context(), usr.ID, 0, 10, nil)
			require.NoError(t, err)
			assert.Equal(t, 2, page.TotalItems)
			require.Len(t, page.Contents, 2)
			assert.Equal(t, first.ID, page.Contents[0].ID)
			assert.Equal(t, second.ID, page.Contents[1].ID)
		})

		t.Run(
			"returns a single page when there are more tokens than the page size",
			func(t *testing.T) {
				usr := newUser()
				require.NoError(t, repo.Save(t.Context(), usr))

				for _, name := range []string{"a-token", "b-token", "c-token"} {
					token := newAPIToken(usr)
					token.Name = name
					require.NoError(t, repo.CreateToken(t.Context(), token))
				}

				page, err := repo.FindTokensByUserID(t.Context(), usr.ID, 0, 2, nil)
				require.NoError(t, err)
				assert.Equal(t, 3, page.TotalItems)
				assert.Equal(t, 2, page.PageSize)
				assert.Equal(t, 0, page.PageNumber)
				require.Len(t, page.Contents, 2)
				assert.Equal(t, "a-token", page.Contents[0].Name)

				next, err := repo.FindTokensByUserID(t.Context(), usr.ID, 1, 2, nil)
				require.NoError(t, err)
				assert.Equal(t, 3, next.TotalItems)
				require.Len(t, next.Contents, 1)
				assert.Equal(t, "c-token", next.Contents[0].Name)
			},
		)

		t.Run("filters by search terms", func(t *testing.T) {
			usr := newUser()
			require.NoError(t, repo.Save(t.Context(), usr))

			matching := newAPIToken(usr)
			matching.Name = "production-automation"
			require.NoError(t, repo.CreateToken(t.Context(), matching))

			other := newAPIToken(usr)
			other.Name = "staging-automation"
			require.NoError(t, repo.CreateToken(t.Context(), other))

			searchTerms := "PRODUCTION"
			page, err := repo.FindTokensByUserID(t.Context(), usr.ID, 0, 10, &searchTerms)
			require.NoError(t, err)
			assert.Equal(t, 1, page.TotalItems)
			require.Len(t, page.Contents, 1)
			assert.Equal(t, matching.ID, page.Contents[0].ID)
		})

		t.Run("returns an empty page when the user has no tokens", func(t *testing.T) {
			usr := newUser()
			require.NoError(t, repo.Save(t.Context(), usr))

			page, err := repo.FindTokensByUserID(t.Context(), usr.ID, 0, 10, nil)
			require.NoError(t, err)
			assert.Equal(t, 0, page.TotalItems)
			assert.Empty(t, page.Contents)
		})
	})

	t.Run("ExistsTokenByName", func(t *testing.T) {
		t.Run("returns true for an existing name", func(t *testing.T) {
			usr := newUser()
			require.NoError(t, repo.Save(t.Context(), usr))

			token := newAPIToken(usr)
			token.Name = "automation"
			require.NoError(t, repo.CreateToken(t.Context(), token))

			exists, err := repo.ExistsTokenByName(t.Context(), usr.ID, "AUTOMATION")
			require.NoError(t, err)
			assert.True(t, exists)
		})

		t.Run("returns false for an unknown name", func(t *testing.T) {
			usr := newUser()
			require.NoError(t, repo.Save(t.Context(), usr))
			require.NoError(t, repo.CreateToken(t.Context(), newAPIToken(usr)))

			exists, err := repo.ExistsTokenByName(t.Context(), usr.ID, "unknown")
			require.NoError(t, err)
			assert.False(t, exists)
		})

		t.Run("returns false for a name owned by another user", func(t *testing.T) {
			usr := newUser()
			require.NoError(t, repo.Save(t.Context(), usr))
			other := newUser()
			require.NoError(t, repo.Save(t.Context(), other))

			token := newAPIToken(other)
			token.Name = "automation"
			require.NoError(t, repo.CreateToken(t.Context(), token))

			exists, err := repo.ExistsTokenByName(t.Context(), usr.ID, "automation")
			require.NoError(t, err)
			assert.False(t, exists)
		})
	})

	t.Run("FindTokenByID", func(t *testing.T) {
		t.Run("returns nil when the token belongs to another user", func(t *testing.T) {
			usr := newUser()
			require.NoError(t, repo.Save(t.Context(), usr))
			other := newUser()
			require.NoError(t, repo.Save(t.Context(), other))
			token := newAPIToken(other)
			require.NoError(t, repo.CreateToken(t.Context(), token))

			saved, err := repo.FindTokenByID(t.Context(), usr.ID, token.ID)
			require.NoError(t, err)
			assert.Nil(t, saved)
		})

		t.Run("returns nil when not found", func(t *testing.T) {
			usr := newUser()
			require.NoError(t, repo.Save(t.Context(), usr))

			saved, err := repo.FindTokenByID(t.Context(), usr.ID, uuid.New())
			require.NoError(t, err)
			assert.Nil(t, saved)
		})
	})

	t.Run("DeleteTokenByID", func(t *testing.T) {
		t.Run("removes the token", func(t *testing.T) {
			usr := newUser()
			require.NoError(t, repo.Save(t.Context(), usr))
			token := newAPIToken(usr)
			require.NoError(t, repo.CreateToken(t.Context(), token))

			err := repo.DeleteTokenByID(t.Context(), usr.ID, token.ID)
			require.NoError(t, err)

			saved, err := repo.FindTokenByID(t.Context(), usr.ID, token.ID)
			require.NoError(t, err)
			assert.Nil(t, saved)
		})

		t.Run("keeps tokens owned by other users", func(t *testing.T) {
			usr := newUser()
			require.NoError(t, repo.Save(t.Context(), usr))
			other := newUser()
			require.NoError(t, repo.Save(t.Context(), other))
			token := newAPIToken(other)
			require.NoError(t, repo.CreateToken(t.Context(), token))

			require.NoError(t, repo.DeleteTokenByID(t.Context(), usr.ID, token.ID))

			saved, err := repo.FindTokenByID(t.Context(), other.ID, token.ID)
			require.NoError(t, err)
			require.NotNil(t, saved)
		})
	})

	t.Run("TryUpdateLastUsedTOTPCode", func(t *testing.T) {
		t.Run("successfully updates on first code", func(t *testing.T) {
			usr := newUser()
			require.NoError(t, repo.Save(t.Context(), usr))

			ok, err := repo.TryUpdateLastUsedTOTPCode(t.Context(), usr.ID, "123456")
			require.NoError(t, err)
			assert.True(t, ok)

			saved, _ := repo.FindByID(t.Context(), usr.ID)
			assert.Equal(t, []string{"123456"}, saved.TOTP.LastUsedCodes)
		})

		t.Run("fails on exact replay", func(t *testing.T) {
			usr := newUser()
			require.NoError(t, repo.Save(t.Context(), usr))

			ok, err := repo.TryUpdateLastUsedTOTPCode(t.Context(), usr.ID, "111111")
			require.NoError(t, err)
			assert.True(t, ok)

			ok, err = repo.TryUpdateLastUsedTOTPCode(t.Context(), usr.ID, "111111")
			require.NoError(t, err)
			assert.False(t, ok)
		})

		t.Run("successfully prepends new unique codes", func(t *testing.T) {
			usr := newUser()
			require.NoError(t, repo.Save(t.Context(), usr))

			_, _ = repo.TryUpdateLastUsedTOTPCode(t.Context(), usr.ID, "100001")
			_, _ = repo.TryUpdateLastUsedTOTPCode(t.Context(), usr.ID, "100002")

			saved, _ := repo.FindByID(t.Context(), usr.ID)
			assert.Equal(t, []string{"100002", "100001"}, saved.TOTP.LastUsedCodes)
		})

		t.Run("fails on replay of older code in history", func(t *testing.T) {
			usr := newUser()
			require.NoError(t, repo.Save(t.Context(), usr))

			_, _ = repo.TryUpdateLastUsedTOTPCode(t.Context(), usr.ID, "200001")
			_, _ = repo.TryUpdateLastUsedTOTPCode(t.Context(), usr.ID, "200002")
			_, _ = repo.TryUpdateLastUsedTOTPCode(t.Context(), usr.ID, "200003")

			ok, err := repo.TryUpdateLastUsedTOTPCode(t.Context(), usr.ID, "200001")
			require.NoError(t, err)
			assert.False(t, ok)
		})

		t.Run("maintains only the last few codes based on limit", func(t *testing.T) {
			usr := newUser()
			require.NoError(t, repo.Save(t.Context(), usr))

			codes := []string{"300001", "300002", "300003", "300004", "300005"}
			for _, code := range codes {
				_, _ = repo.TryUpdateLastUsedTOTPCode(t.Context(), usr.ID, code)
			}

			saved, err := repo.FindByID(t.Context(), usr.ID)
			require.NoError(t, err)
			assert.Equal(t, 3, len(saved.TOTP.LastUsedCodes))
			assert.Equal(t, "300005", saved.TOTP.LastUsedCodes[0])
			assert.Equal(t, "300004", saved.TOTP.LastUsedCodes[1])
			assert.Equal(t, "300003", saved.TOTP.LastUsedCodes[2])
		})
	})
}
