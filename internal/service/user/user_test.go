package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	authModel "github.com/lin-snow/ech0/internal/model/auth"
	commonModel "github.com/lin-snow/ech0/internal/model/common"
	settingModel "github.com/lin-snow/ech0/internal/model/setting"
	model "github.com/lin-snow/ech0/internal/model/user"
	cryptoUtil "github.com/lin-snow/ech0/internal/util/crypto"
)

// MockUserRepository 模拟用户仓库接口
type MockUserRepository struct{ mock.Mock }

func (m *MockUserRepository) GetAllUsers() ([]model.User, error) {
	args := m.Called()
	return args.Get(0).([]model.User), args.Error(1)
}
func (m *MockUserRepository) GetUserByUsername(username string) (model.User, error) {
	args := m.Called(username)
	return args.Get(0).(model.User), args.Error(1)
}
func (m *MockUserRepository) CreateUser(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}
func (m *MockUserRepository) GetUserByID(id int) (model.User, error) { return model.User{}, nil }
func (m *MockUserRepository) UpdateUser(user *model.User) error      { return nil }
func (m *MockUserRepository) DeleteUser(id uint) error               { return nil }
func (m *MockUserRepository) GetSysAdmin() (model.User, error)       { return model.User{}, nil }

// MockSettingService 模拟设置服务接口
type MockSettingService struct{ mock.Mock }

func (m *MockSettingService) GetSetting(setting *settingModel.SystemSetting) error {
	args := m.Called(setting)
	if args.Error(0) == nil {
		if s, ok := args.Get(1).(settingModel.SystemSetting); ok {
			*setting = s
		}
	}
	return args.Error(0)
}
func (m *MockSettingService) GetCommentSetting(setting *settingModel.CommentSetting) error { return nil }
func (m *MockSettingService) UpdateCommentSetting(id uint, setting *settingModel.CommentSettingDto) error {
	return nil
}
func (m *MockSettingService) UpdateSetting(id uint, setting *settingModel.SystemSettingDto) error { return nil }

// 测试套件
type UserServiceTestSuite struct {
	suite.Suite
	userService    *UserService
	mockUserRepo   *MockUserRepository
	mockSettingSvc *MockSettingService
}

func (suite *UserServiceTestSuite) SetupTest() {
	suite.mockUserRepo = new(MockUserRepository)
	suite.mockSettingSvc = new(MockSettingService)
	suite.userService = &UserService{
		userRepository: suite.mockUserRepo,
		settingService: suite.mockSettingSvc,
	}
}

// ✅ 测试首个用户注册 → 自动成为管理员
func (suite *UserServiceTestSuite) TestRegister_FirstUser_ShouldBeAdmin() {
	registerDto := &authModel.RegisterDto{Username: "admin", Password: "password123"}

	// Mock: 没有现有用户
	suite.mockUserRepo.On("GetAllUsers").Return([]model.User{}, nil)

	// Mock: 用户名不存在
	suite.mockUserRepo.On("GetUserByUsername", "admin").Return(
		model.User{ID: model.USER_NOT_EXISTS_ID}, errors.New("user not found"),
	)

	// Mock: GetSetting 即使没用到，也要返回默认值
	suite.mockSettingSvc.On("GetSetting", mock.Anything).Return(nil, settingModel.SystemSetting{})

	// Mock: 成功创建用户
	suite.mockUserRepo.On("CreateUser", mock.MatchedBy(func(user *model.User) bool {
		return user.Username == "admin" &&
			user.Password == cryptoUtil.MD5Encrypt("password123") &&
			user.IsAdmin
	})).Return(nil)

	err := suite.userService.Register(registerDto)

	assert.NoError(suite.T(), err)
	suite.mockUserRepo.AssertExpectations(suite.T())
	suite.mockSettingSvc.AssertExpectations(suite.T())
}

// 🚫 测试已有用户时禁止注册
func (suite *UserServiceTestSuite) TestRegister_RegistrationNotAllowed() {
	registerDto := &authModel.RegisterDto{Username: "user1", Password: "password123"}
	existingUsers := []model.User{{ID: 1, Username: "admin", IsAdmin: true}}
	setting := settingModel.SystemSetting{AllowRegister: false}

	// Mock: 已有用户
	suite.mockUserRepo.On("GetAllUsers").Return(existingUsers, nil)
	// Mock: 用户名不存在
	suite.mockUserRepo.On("GetUserByUsername", "user1").Return(
		model.User{ID: model.USER_NOT_EXISTS_ID}, errors.New("user not found"),
	)
	// Mock: 不允许注册（放宽匹配条件）
	suite.mockSettingSvc.On("GetSetting", mock.Anything).Return(nil, setting)

	err := suite.userService.Register(registerDto)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), commonModel.USER_REGISTER_NOT_ALLOW, err.Error())
	suite.mockUserRepo.AssertExpectations(suite.T())
	suite.mockSettingSvc.AssertExpectations(suite.T())
}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}
