package repository

import (
	"github.com/lin-snow/ech0/internal/cache"
	model "github.com/lin-snow/ech0/internal/model/user"
	"gorm.io/gorm"
)

type UserRepository struct {
	db    *gorm.DB
	cache cache.ICache[string, *model.User]
}

func NewUserRepository(db *gorm.DB, cache cache.ICache[string, *model.User]) UserRepositoryInterface {
	return &UserRepository{
		db:    db,
		cache: cache,
	}
}

// GetUserByUsername 根据用户名获取用户
func (userRepository *UserRepository) GetUserByUsername(username string) (model.User, error) {
	// 先查缓存
	cacheKey := GetUsernameKey(username)
	if cachedUser, err := userRepository.cache.Get(cacheKey); err == nil {
		return *cachedUser, nil
	}

	// 缓存未命中，查询数据库
	user := model.User{}
	err := userRepository.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return model.User{}, err
	}

	// 写缓存， cost 设为1
	userRepository.cache.Set(cacheKey, &user, 1)

	return user, nil
}

// GetAllUsers 获取所有用户
func (userRepository *UserRepository) GetAllUsers() ([]model.User, error) {
	var users []model.User
	err := userRepository.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// CreateUser 创建一个新的用户
func (userRepository *UserRepository) CreateUser(user *model.User) error {
	err := userRepository.db.Create(user).Error
	if err != nil {
		return err
	}

	// 加入缓存
	userRepository.cache.Set(GetUserIDKey(user.ID), user, 1)
	userRepository.cache.Set(GetUsernameKey(user.Username), user, 1)

	return nil
}

// GetUserByID 根据用户ID获取用户
func (userRepository *UserRepository) GetUserByID(id int) (model.User, error) {
	cacheKey := GetUserIDKey(uint(id))
	if cachedUser, err := userRepository.cache.Get(cacheKey); err == nil {
		return *cachedUser, nil
	}

	var user model.User
	if err := userRepository.db.First(&user, id).Error; err != nil {
		return user, err
	}

	userRepository.cache.Set(cacheKey, &user, 1)

	return user, nil
}

// GetSysAdmin 获取系统管理员
func (userRepository *UserRepository) GetSysAdmin() (model.User, error) {
	cacheKey := GetSysAdminKey()
	if cachedUser, err := userRepository.cache.Get(cacheKey); err == nil {
		return *cachedUser, nil
	}

	// 获取系统管理员（首个注册的用户）
	user := model.User{}
	err := userRepository.db.Where("is_admin = ?", true).First(&user).Error
	if err != nil {
		return model.User{}, err
	}

	userRepository.cache.Set(cacheKey, &user, 1)

	return user, nil
}

// UpdateUser 更新用户信息
func (userRepository *UserRepository) UpdateUser(user *model.User) error {
	err := userRepository.db.Save(user).Error
	if err != nil {
		return err
	}

	userRepository.cache.Set(GetUserIDKey(user.ID), user, 1)
	userRepository.cache.Set(GetUsernameKey(user.Username), user, 1)
	if user.IsAdmin {
		userRepository.cache.Set(GetAdminKey(user.ID), user, 1)
	}

	return nil
}

// DeleteUser 删除用户
func (userRepository *UserRepository) DeleteUser(id uint) error {
	// 先查找待删除的用户
	userToDel, err := userRepository.GetUserByID(int(id))
	if err != nil {
		return err
	}

	err = userRepository.db.Delete(&model.User{}, id).Error
	if err != nil {
		return err
	}

	// 清空缓存
	userRepository.cache.Delete(GetUserIDKey(userToDel.ID))
	userRepository.cache.Delete(GetUsernameKey(userToDel.Username))
	if userToDel.IsAdmin {
		userRepository.cache.Delete(GetAdminKey(userToDel.ID))
	}

	return nil
}
