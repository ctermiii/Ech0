package repository

import (
	"errors"
	"strings"
	"time"

	"github.com/lin-snow/ech0/internal/cache"
	commonModel "github.com/lin-snow/ech0/internal/model/common"
	model "github.com/lin-snow/ech0/internal/model/echo"
	"gorm.io/gorm"
)

type EchoRepository struct {
	db    *gorm.DB
	cache cache.ICache[string, commonModel.PageQueryResult[[]model.Echo]]
}

func NewEchoRepository(db *gorm.DB, cache cache.ICache[string, commonModel.PageQueryResult[[]model.Echo]]) EchoRepositoryInterface {
	return &EchoRepository{db: db, cache: cache}
}

func (echoRepository *EchoRepository) CreateEcho(echo *model.Echo) error {
	echo.Content = strings.TrimSpace(echo.Content)

	result := echoRepository.db.Create(echo)
	if result.Error != nil {
		return result.Error
	}

	ClearEchoPageCache(echoRepository.cache)

	return nil
}

// GetEchosByPage 获取分页的 Echo 列表
func (echoRepository *EchoRepository) GetEchosByPage(page, pageSize int, search string, showPrivate bool) ([]model.Echo, int64) {
	// 查找缓存
	cacheKey := GetEchoPageCacheKey(page, pageSize, search, showPrivate)
	if cachedResult, err := echoRepository.cache.Get(cacheKey); err == nil {
		return cachedResult.Items, cachedResult.Total
	}

	// 如果缓存未命中，进行数据库查询

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 查询数据库
	var echos []model.Echo
	var total int64

	query := echoRepository.db.Model(&model.Echo{})

	// 如果 search 不为空，添加模糊查询条件
	if search != "" {
		searchPattern := "%" + search + "%" // 模糊匹配模式
		query = query.Where("content LIKE ?", searchPattern)
	}

	// 如果不是管理员，过滤私密Echo
	if !showPrivate {
		query = query.Where("private = ?", false)
	}

	// 获取总数并进行分页查询
	query.Count(&total).
		Preload("Images").
		Limit(pageSize).
		Offset(offset).
		Order("created_at DESC").
		Find(&echos)

	// 保存到缓存
	echoKeyList = append(echoKeyList, cacheKey) // 记录缓存键
	echoRepository.cache.Set(cacheKey, commonModel.PageQueryResult[[]model.Echo]{
		Items: echos,
		Total: total,
	}, 1)

	// 返回结果
	return echos, total
}

// GetEchosById 根据 ID 获取 Echo
func (echoRepository *EchoRepository) GetEchosById(id uint) (*model.Echo, error) {
	var echo model.Echo
	result := echoRepository.db.Preload("Images").First(&echo, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // 如果未找到记录，则返回 nil
		}
		return nil, result.Error // 其他错误返回
	}

	return &echo, nil
}

// DeleteEchoById 删除 Echo
func (echoRepository *EchoRepository) DeleteEchoById(id uint) error {
	var echo model.Echo
	// 删除外键images
	echoRepository.db.Where("message_id = ?", id).Delete(&model.Image{})

	result := echoRepository.db.Delete(&echo, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // 如果没有找到记录
	}

	// 清除相关缓存
	ClearEchoPageCache(echoRepository.cache)

	return nil
}

// GetTodayEchos 获取今天的 Echo 列表
func (echoRepository *EchoRepository) GetTodayEchos(showPrivate bool) []model.Echo {
	// 查询数据库
	var echos []model.Echo

	// 获取当天开始和结束时间
	today := time.Now()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	query := echoRepository.db.Model(&model.Echo{})
	// 如果不是管理员，过滤私密Echo
	if !showPrivate {
		query = query.Where("private = ?", false)
	}

	// 添加当天的时间过滤
	query = query.Where("created_at >= ? AND created_at < ?", startOfDay, endOfDay)

	// 获取总数并进行分页查询
	query.
		Preload("Images").
		Order("created_at DESC").
		Find(&echos)

	// 返回结果
	return echos
}

// UpdateEcho 更新 Echo
func (echoRepository *EchoRepository) UpdateEcho(echo *model.Echo) error {
	// 清空缓存
	ClearEchoPageCache(echoRepository.cache)

	// 开启事务确保数据一致性
	tx := echoRepository.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 先删除该 Echo 关联的所有旧图片
	if err := tx.Where("message_id = ?", echo.ID).Delete(&model.Image{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 2. 更新 Echo 内容（包括关联的新图片）
	if err := tx.Model(&model.Echo{}).
		Where("id = ?", echo.ID).
		Updates(map[string]interface{}{
			"content":        echo.Content,
			"private":        echo.Private,
			"extension":      echo.Extension,
			"extension_type": echo.ExtensionType,
		}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 3. 重新添加Images
	if len(echo.Images) > 0 {
		var images []model.Image
		for _, img := range echo.Images {
			// 确保每个图片都关联到正确的 Echo ID
			img.MessageID = echo.ID
			images = append(images, img)
		}
		// 批量插入新图片
		if err := tx.Create(&images).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// 提交事务
	return tx.Commit().Error
}

// LikeEcho 点赞 Echo
func (echoRepository *EchoRepository) LikeEcho(id uint) error {
	// 检查是否存在（可选，防止无效点赞）
	var exists bool
	if err := echoRepository.db.
		Model(&model.Echo{}).
		Select("count(*) > 0").
		Where("id = ?", id).
		Find(&exists).Error; err != nil {
		return err
	}
	if !exists {
		return errors.New(commonModel.ECHO_NOT_FOUND)
	}

	// 原子自增点赞数
	if err := echoRepository.db.
		Model(&model.Echo{}).
		Where("id = ?", id).
		UpdateColumn("fav_count", gorm.Expr("fav_count + ?", 1)).Error; err != nil {
		return err
	}

	// 清除相关缓存
	ClearEchoPageCache(echoRepository.cache)

	return nil
}
