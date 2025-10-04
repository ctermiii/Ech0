package database

import (
	"errors"
	"log"
	"time"

	commonModel "github.com/lin-snow/ech0/internal/model/common"
	echoModel "github.com/lin-snow/ech0/internal/model/echo"
	"gorm.io/gorm"
)

// UpdateMigration 执行数据库迁移，将旧版 Message 表的数据迁移到新版 Echo 表
func UpdateMigration() {
	// 确认旧表存在，否则无需迁移
	if !GetDB().Migrator().HasTable(&echoModel.Message{}) {
		// log.Println("未发现旧版 Message 表，无需迁移。")
		return
	}

	log.Println("警告：将执行“清空并重建”迁移。此操作会删除 echos 和 images 表中的所有数据！")
	log.Println("开始执行迁移...")

	var kvFlag commonModel.KeyValue
	result := GetDB().First(&kvFlag, "key = ?", commonModel.MigrationKey).Error

	if result == nil {
		// 如果找到了记录，说明已迁移
		log.Printf("发现迁移标记 Key: '%s' (Value: '%s')，任务已完成，跳过执行。", kvFlag.Key, kvFlag.Value)
		return
	}

	// 如果错误不是 "record not found"，说明可能发生了其他数据库问题
	if !errors.Is(result, gorm.ErrRecordNotFound) {
		log.Fatalf("查询迁移标记时发生意外错误: %v", result)
		return
	}

	// 如果错误是 gorm.ErrRecordNotFound，说明标记不存在，可以开始迁移
	log.Println("未发现迁移标记，准备执行迁移任务...")

	// =================================================================
	// 步骤 1: 一次性从旧表中加载所有数据到内存
	// =================================================================
	var messages []echoModel.Message
	// 使用 Preload 一并加载关联的 Images
	if err := GetDB().Preload("Images").Find(&messages).Error; err != nil {
		log.Fatalf("从旧表加载数据失败: %v", err)
		return
	}

	if len(messages) == 0 {
		log.Println("旧表中没有数据，迁移完成。")
		return
	}
	log.Printf("成功从旧表加载 %d 条 message 记录。", len(messages))

	// =================================================================
	// 步骤 2: 启动一个事务来执行所有数据库写操作
	// =================================================================
	err := GetDB().Transaction(func(tx *gorm.DB) error {
		// 步骤 3: 清空新表
		// 注意删除顺序：先删除带有外键的表 (images)，再删除被引用的表 (echos)
		log.Println("在事务中清空 echos 和 images 表...")

		// 使用 Exec("DELETE FROM ...") 而不是 Truncate，因为 DELETE 可以在事务中回滚
		if err := tx.Exec("DELETE FROM images").Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM echos").Error; err != nil {
			return err
		}
		log.Println("新表已清空。")

		// 步骤 4: 遍历内存中的旧数据，一个一个地重建到新表中
		log.Println("开始逐条重建数据...")
		for _, msg := range messages {
			// a. 创建新的 Echo 对象
			echo := echoModel.Echo{
				ID:            msg.ID, // <-- 关键：直接使用旧的ID
				Content:       msg.Content,
				Username:      msg.Username,
				Private:       msg.Private,
				UserID:        msg.UserID,
				Extension:     msg.Extension,
				ExtensionType: msg.ExtensionType,
				CreatedAt:     msg.CreatedAt, // 保留原始创建时间
			}

			// b. 在事务中创建 Echo 记录
			if err := tx.Create(&echo).Error; err != nil {
				log.Printf("创建 echo 记录失败 (ID: %d): %v", msg.ID, err)
				return err // 任何错误都会导致事务回滚
			}

			// c. 如果有关联的图片，也一并创建
			if len(msg.Images) > 0 {
				newImages := make([]echoModel.Image, len(msg.Images))
				for i, img := range msg.Images {
					newImages[i] = echoModel.Image{
						// ID 会自增，我们不关心，重要的是 MessageID
						MessageID:   msg.ID, // <-- 关键：直接使用旧的ID作为外键
						ImageURL:    img.ImageURL,
						ImageSource: img.ImageSource,
					}
				}
				// d. 在事务中批量创建这组图片
				if err := tx.Create(&newImages).Error; err != nil {
					log.Printf("为 echo ID %d 创建 images 失败: %v", msg.ID, err)
					return err // 任何错误都会导致事务回滚
				}
			}
		}

		log.Println("所有数据重建完毕，正在写入迁移标记...")
		completionValue := "completed_at_" + time.Now().Format(time.RFC3339)
		migrationFlag := commonModel.KeyValue{
			Key:   commonModel.MigrationKey,
			Value: completionValue,
		}

		if err := tx.Create(&migrationFlag).Error; err != nil {
			return err // 如果写入标记失败，整个事务回滚
		}

		log.Println("所有数据已成功准备好，提交事务...")
		return nil // 返回 nil 表示事务成功，可以提交
	})

	if err != nil {
		log.Fatalf("迁移过程中发生错误，事务已回滚，数据库未发生任何改变: %v", err)
	} else {
		log.Println("迁移成功完成！")
	}
}
