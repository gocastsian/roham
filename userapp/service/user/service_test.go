package user_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"mime/multipart"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gocastsian/roham/types"
	"github.com/gocastsian/roham/userapp/service/guard"
	"github.com/gocastsian/roham/userapp/service/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetAllUsers(ctx context.Context) ([]user.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]user.User), args.Error(1)
}

func (m *MockRepository) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (user.User, error) {
	args := m.Called(ctx, phoneNumber)
	return args.Get(0).(user.User), args.Error(1)
}

func (m *MockRepository) CheckUserUniquness(ctx context.Context, email string, username string) (bool, error) {
	args := m.Called(ctx, email, username)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) GetUser(ctx context.Context, ID types.ID) (user.User, error) {
	args := m.Called(ctx, ID)
	return args.Get(0).(user.User), args.Error(1)
}

func (m *MockRepository) CheckUserExistByID(ctx context.Context, ID types.ID) (bool, error) {
	args := m.Called(ctx, ID)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) RegisterUser(ctx context.Context, u user.User) (types.ID, error) {
	args := m.Called(ctx, u)
	return args.Get(0).(types.ID), args.Error(1)
}
func (m *MockRepository) UpdateAvatar(ctx context.Context, ID types.ID, uploadAddress string) error {
	fmt.Printf("userId:%d destAddre:%s\n", ID, uploadAddress)
	args := m.Called(ctx, ID, uploadAddress)
	return args.Error(0)
}

type mockMultipartFile struct {
	*bytes.Reader
}

func (m *mockMultipartFile) Close() error {
	return nil
}

func TestRegisterUser_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	userValidator := user.NewValidator(mockRepo)
	userConf := user.Config{}
	guardSvc := &guard.Service{}
	service := user.NewService(mockRepo, userValidator, nil, guardSvc, userConf)

	regReq := user.RegisterRequest{
		Username:        "testuser",
		FirstName:       "Test",
		LastName:        "User",
		Email:           "test@example.com",
		Password:        "s2Securepassword",
		ConfirmPassword: "s2Securepassword",
	}

	mockRepo.On("CheckUserUniquness", mock.Anything, regReq.Email, regReq.Username).Return(false, nil)
	mockRepo.On("RegisterUser", mock.Anything, mock.AnythingOfType("user.User")).Return(types.ID(1), nil)

	resp, err := service.RegisterUser(context.Background(), regReq)
	assert.NoError(t, err)
	assert.Equal(t, types.ID(1), resp.ID)
	mockRepo.AssertExpectations(t)
}

func TestRegisterUser_UserAlreadyExists(t *testing.T) {
	mockRepo := new(MockRepository)
	userConf := user.Config{}
	userValidator := user.NewValidator(mockRepo)
	service := user.NewService(mockRepo, userValidator, nil, nil, userConf)

	regReq := user.RegisterRequest{
		Username:        "testuser",
		FirstName:       "Test",
		LastName:        "User",
		Email:           "test@example.com",
		Password:        "se2Scurepassword",
		ConfirmPassword: "se2Scurepassword",
	}

	mockRepo.On("CheckUserUniquness", mock.Anything, regReq.Email, regReq.Username).Return(true, nil)

	_, err := service.RegisterUser(context.Background(), regReq)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user already exist")
	mockRepo.AssertExpectations(t)
}

func TestGetUser(t *testing.T) {
	mockRepo := new(MockRepository)
	userValidator := user.NewValidator(mockRepo)
	guardSvc := &guard.Service{}
	userConf := user.Config{}
	service := user.NewService(mockRepo, userValidator, nil, guardSvc, userConf)

	testUser := user.User{
		ID:        1,
		Username:  "test",
		FirstName: "firstname",
		LastName:  "lastname",
		Email:     "email@gmail.com",
		Avatar:    "",
		Role:      0,
	}

	type testCase struct {
		name   string
		userId types.ID
		err    error
		user   user.User
	}
	testCases := []testCase{
		{
			name:   "not found a user",
			userId: types.ID(0),
			err:    fmt.Errorf("the user not found"),
			user:   user.User{},
		},
		{
			name:   "successfully get a user",
			userId: types.ID(1),
			err:    nil,
			user:   testUser,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.On("GetUser", mock.Anything, tc.userId).Return(tc.user, tc.err)

			user, err := service.GetUser(context.Background(), tc.userId)
			if tc.err != nil {
				assert.Error(t, tc.err)
				assert.Contains(t, err.Error(), tc.err.Error())
			} else {
				assert.Equal(t, user.ID, user.ID)
			}
		})
	}
}

func createTestMultipartFile(t *testing.T, fieldName, filename string, content []byte) (multipart.File, *multipart.FileHeader) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(fieldName, filename)
	require.NoError(t, err)

	_, err = part.Write(content)
	require.NoError(t, err)

	err = writer.Close()
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	err = req.ParseMultipartForm(32 << 20)
	require.NoError(t, err)

	file, fileHeader, err := req.FormFile(fieldName)
	require.NoError(t, err)

	return file, fileHeader
}
func TestUpdateUserAvatar(t *testing.T) {
	mockRepo := new(MockRepository)
	tmpDir := t.TempDir()
	conf := user.Config{
		AvatarConfig: user.AvatarConfig{
			MaxSize:       1, // 1 MB
			ValidFormats:  []string{"image/png"},
			UploadFileDir: tmpDir,
		},
	}
	service := user.NewService(
		mockRepo,
		user.NewValidator(mockRepo),
		nil,
		&guard.Service{},
		conf,
	)

	type TestCase struct {
		name           string
		userId         types.ID
		avatar         *user.Avatar
		expectError    bool
		expectRepoCall bool
	}

	pngContent := []byte("\x89PNG\r\n\x1a\n" + "some image content here")
	file, fileHeader := createTestMultipartFile(t, "avatar", "avatar.png", pngContent)
	validAvatar := user.Avatar{
		FileHandler: fileHeader,
		File:        file,
	}

	exeContent := []byte("MZP...not really an image")
	badFile, badFileHeader := createTestMultipartFile(t, "avatar", "bad.exe", exeContent)
	invalidAvatar := user.Avatar{
		FileHandler: badFileHeader,
		File:        badFile,
	}
	largeContent := bytes.Repeat([]byte("A"), 2*1024*1024) // 2 MB
	largeFile, largeFileHeader := createTestMultipartFile(t, "avatar", "large.png", largeContent)
	largeAvatar := user.Avatar{
		FileHandler: largeFileHeader,
		File:        largeFile,
	}

	testCases := []TestCase{
		{
			name:           "success upload avatar",
			userId:         types.ID(1),
			avatar:         &validAvatar,
			expectError:    false,
			expectRepoCall: true,
		},
		{
			name:           "invalid file (wrong MIME)",
			userId:         types.ID(1),
			avatar:         &invalidAvatar,
			expectError:    true,
			expectRepoCall: false,
		},
		{
			name:           "file too large",
			userId:         types.ID(1),
			avatar:         &largeAvatar,
			expectError:    true,
			expectRepoCall: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dstPath := filepath.Join(tmpDir, tc.avatar.FileHandler.Filename)

			if tc.expectRepoCall {
				mockRepo.On("UpdateAvatar", mock.Anything, tc.userId, dstPath).Return(nil).Once()
			}

			err := service.UpdateUserAvatar(context.Background(), tc.userId, *tc.avatar)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			fmt.Println("Recorded calls:", mockRepo.Calls)
			if tc.expectRepoCall {
				mockRepo.AssertCalled(t, "UpdateAvatar", mock.Anything, tc.userId, dstPath)
				mockRepo.ExpectedCalls = nil
				mockRepo.Calls = nil
			} else {
				mockRepo.AssertNotCalled(t, "UpdateAvatar", mock.Anything, tc.userId, mock.Anything)
			}
		})
	}
}
