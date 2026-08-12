package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/minhhoccode111/go-clean-template-gin/internal/entity"
	"github.com/minhhoccode111/go-clean-template-gin/internal/usecase"
	"github.com/minhhoccode111/go-clean-template-gin/internal/usecase/user"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

const (
	testUsername = "testuser"
	testEmail    = "test@example.com"
	userTestID   = "user-id-123"
)

func newUserUseCase(t *testing.T) (usecase.User, *MockUserRepo) {
	t.Helper()

	ctrl := gomock.NewController(t)

	repo := NewMockUserRepo(ctrl)
	jwtManager := jwt.New("test-secret", time.Hour)
	useCase := user.New(repo, jwtManager)

	return useCase, repo
}

func TestRegister(t *testing.T) {
	t.Parallel()

	t.Run("register success", func(t *testing.T) {
		t.Parallel()

		uc, repo := newUserUseCase(t)
		repo.EXPECT().Store(gomock.Any(), gomock.Any()).Return(nil)

		u, err := uc.Register(context.Background(), testUsername, testEmail, "password123")

		require.NoError(t, err)
		assert.NotEmpty(t, u.ID)
		assert.Equal(t, testUsername, u.Username)
		assert.Equal(t, testEmail, u.Email)
	})

	t.Run("register duplicate", func(t *testing.T) {
		t.Parallel()

		uc, repo := newUserUseCase(t)
		repo.EXPECT().Store(gomock.Any(), gomock.Any()).Return(entity.ErrUserAlreadyExists)

		_, err := uc.Register(context.Background(), testUsername, testEmail, "password123")

		require.ErrorIs(t, err, entity.ErrUserAlreadyExists)
	})
}

func TestLogin(t *testing.T) {
	t.Parallel()

	t.Run("login success", func(t *testing.T) {
		t.Parallel()

		uc, repo := newUserUseCase(t)
		hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		require.NoError(t, err)

		storedUser := entity.User{
			ID: userTestID, Username: testUsername,
			Email: testEmail, PasswordHash: string(hash),
		}
		repo.EXPECT().GetByEmail(gomock.Any(), testEmail).Return(storedUser, nil)

		token, err := uc.Login(context.Background(), testEmail, "password123")

		require.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("login wrong password", func(t *testing.T) {
		t.Parallel()

		uc, repo := newUserUseCase(t)
		hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		require.NoError(t, err)

		storedUser := entity.User{
			ID: userTestID, Username: testUsername,
			Email: testEmail, PasswordHash: string(hash),
		}
		repo.EXPECT().GetByEmail(gomock.Any(), testEmail).Return(storedUser, nil)

		token, err := uc.Login(context.Background(), testEmail, "wrongpassword")

		require.ErrorIs(t, err, entity.ErrInvalidCredentials)
		assert.Empty(t, token)
	})

	t.Run("login user not found", func(t *testing.T) {
		t.Parallel()

		uc, repo := newUserUseCase(t)
		repo.EXPECT().
			GetByEmail(gomock.Any(), "notfound@example.com").
			Return(entity.User{}, entity.ErrUserNotFound)

		token, err := uc.Login(context.Background(), "notfound@example.com", "password123")

		require.ErrorIs(t, err, entity.ErrInvalidCredentials)
		assert.Empty(t, token)
	})
}

func TestGetUser(t *testing.T) {
	t.Parallel()

	expectedUser := entity.User{
		ID:       userTestID,
		Username: testUsername,
		Email:    testEmail,
	}

	t.Run("get user success", func(t *testing.T) {
		t.Parallel()

		uc, repo := newUserUseCase(t)
		repo.EXPECT().GetByID(gomock.Any(), userTestID).Return(expectedUser, nil)

		u, err := uc.GetUser(context.Background(), userTestID)

		require.NoError(t, err)
		assert.Equal(t, expectedUser, u)
	})

	t.Run("get user not found", func(t *testing.T) {
		t.Parallel()

		uc, repo := newUserUseCase(t)
		repo.EXPECT().
			GetByID(gomock.Any(), "missing-id").
			Return(entity.User{}, entity.ErrUserNotFound)

		_, err := uc.GetUser(context.Background(), "missing-id")

		require.ErrorIs(t, err, entity.ErrUserNotFound)
	})
}

func TestGetUser_GenericError(t *testing.T) {
	t.Parallel()

	uc, repo := newUserUseCase(t)

	repo.EXPECT().GetByID(gomock.Any(), userTestID).Return(entity.User{}, errInternalServErr)

	_, err := uc.GetUser(context.Background(), userTestID)

	require.Error(t, err)
	require.ErrorIs(t, err, errInternalServErr)
}
