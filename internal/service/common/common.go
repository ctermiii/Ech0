package service

import (
	"context"
	"errors"
	"fmt"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"

	"github.com/lin-snow/ech0/internal/config"
	"github.com/lin-snow/ech0/internal/event"
	commonModel "github.com/lin-snow/ech0/internal/model/common"
	echoModel "github.com/lin-snow/ech0/internal/model/echo"
	settingModel "github.com/lin-snow/ech0/internal/model/setting"
	userModel "github.com/lin-snow/ech0/internal/model/user"
	repository "github.com/lin-snow/ech0/internal/repository/common"
	echoRepository "github.com/lin-snow/ech0/internal/repository/echo"
	keyvalueRepository "github.com/lin-snow/ech0/internal/repository/keyvalue"
	"github.com/lin-snow/ech0/internal/transaction"
	httpUtil "github.com/lin-snow/ech0/internal/util/http"
	jsonUtil "github.com/lin-snow/ech0/internal/util/json"
	mdUtil "github.com/lin-snow/ech0/internal/util/md"
	storageUtil "github.com/lin-snow/ech0/internal/util/storage"
)

type CommonService struct {
	txManager          transaction.TransactionManager
	commonRepository   repository.CommonRepositoryInterface
	objStorage         storageUtil.ObjectStorage
	echoRepository     echoRepository.EchoRepositoryInterface
	keyvalueRepository keyvalueRepository.KeyValueRepositoryInterface
	eventBus           event.IEventBus
}

func NewCommonService(
	tm transaction.TransactionManager,
	commonRepository repository.CommonRepositoryInterface,
	echoRepository echoRepository.EchoRepositoryInterface,
	keyvalueRepository keyvalueRepository.KeyValueRepositoryInterface,
	eventBusProvider func() event.IEventBus,
) CommonServiceInterface {
	return &CommonService{
		txManager:          tm,
		commonRepository:   commonRepository,
		echoRepository:     echoRepository,
		keyvalueRepository: keyvalueRepository,
		objStorage:         nil,
		eventBus:           eventBusProvider(),
	}
}

func (commonService *CommonService) CommonGetUserByUserId(userId uint) (userModel.User, error) {
	return commonService.commonRepository.GetUserByUserId(userId)
}

func (commonService *CommonService) UploadImage(userId uint, file *multipart.FileHeader) (string, error) {
	user, err := commonService.commonRepository.GetUserByUserId(userId)
	if err != nil {
		return "", err
	}
	if !user.IsAdmin {
		return "", errors.New(commonModel.NO_PERMISSION_DENIED)
	}

	// 检查文件类型是否合法
	if !storageUtil.IsAllowedType(file.Header.Get("Content-Type"), config.Config.Upload.AllowedTypes) {
		return "", errors.New(commonModel.FILE_TYPE_NOT_ALLOWED)
	}

	// 检查文件大小是否合法
	if file.Size > int64(config.Config.Upload.ImageMaxSize) {
		return "", errors.New(commonModel.FILE_SIZE_EXCEED_LIMIT)
	}

	// 调用存储函数存储图片
	imageUrl, err := storageUtil.UploadFile(file, commonModel.ImageType, commonModel.LOCAL_FILE)
	if err != nil {
		return "", err
	}

	// 触发图片上传事件
	user.Password = "" // 清除密码字段，避免泄露
	commonService.eventBus.Publish(context.Background(), event.NewEvent(
		event.EventTypeResourceUploaded,
		event.EventPayload{
			event.EventPayloadUser: user,
			event.EventPayloadFile: file.Filename,
			event.EventPayloadURL:  imageUrl,
			event.EventPayloadSize: file.Size,
			event.EventPayloadType: commonModel.ImageType,
		},
	))

	return imageUrl, nil
}

func (commonService *CommonService) DeleteImage(userid uint, url, source, object_key string) error {
	user, err := commonService.commonRepository.GetUserByUserId(userid)
	if err != nil {
		return err
	}
	if !user.IsAdmin {
		return errors.New(commonModel.NO_PERMISSION_DENIED)
	}

	// 检查图片是否存在
	if url == "" {
		return errors.New(commonModel.IMAGE_NOT_FOUND)
	}

	switch source {
	case echoModel.ImageSourceLocal:
		// 获取图片名字（去除前面的/images/)
		imageName := url[len("/images/"):]

		// 构造图片路径
		imagePath := fmt.Sprintf("data/images/%s", imageName)

		// 删除图片
		return storageUtil.DeleteFileFromLocal(imagePath)
	case echoModel.ImageSourceURL:
		// 无需处理
	case echoModel.ImageSourceS3:
		if object_key == "" {
			// 如果没有传入 object_key，则无法删除,忽略
			return nil
		}

		_, _, err := commonService.GetS3Client()
		if err != nil {
			// 如果没有配置 S3，则无法删除,忽略
			return nil
		}

		// 删除 S3 上的图片
		return commonService.objStorage.DeleteObject(context.Background(), object_key)

	case echoModel.ImageSourceR2:
		// TODO: 实现R2图片删除
	default:
		// 未知图片来源按本地图片处理
		// 获取图片名字（去除前面的/images/)
		imageName := url[len("/images/"):]

		// 构造图片路径
		imagePath := fmt.Sprintf("data/images/%s", imageName)

		// 删除图片
		return storageUtil.DeleteFileFromLocal(imagePath)
	}

	return nil
}

func (commonService *CommonService) DirectDeleteImage(url, source, object_key string) error {
	// 检查图片是否存在
	if url == "" {
		return errors.New(commonModel.IMAGE_NOT_FOUND)
	}

	switch source {
	case echoModel.ImageSourceLocal:
		// 获取图片名字（去除前面的/images/)
		imageName := url[len("/images/"):]

		// 构造图片路径
		imagePath := fmt.Sprintf("data/images/%s", imageName)

		// 删除图片
		return storageUtil.DeleteFileFromLocal(imagePath)
	case echoModel.ImageSourceURL:
		// 无需处理
	case echoModel.ImageSourceS3:
		cli, _, err := commonService.GetS3Client()
		if err != nil {
			// 如果没有配置 S3，则无法删除,忽略
			return nil
		}
		if object_key == "" {
			// 如果没有传入 object_key，则无法删除,忽略
			return nil
		}
		// 删除 S3 上的图片
		return cli.DeleteObject(context.Background(), object_key)
	case echoModel.ImageSourceR2:
		// TODO: 实现R2图片删除
	default:
		// 未知图片来源按本地图片处理
		// 获取图片名字（去除前面的/images/)
		imageName := url[len("/images/"):]

		// 构造图片路径
		imagePath := fmt.Sprintf("data/images/%s", imageName)

		// 删除图片
		return storageUtil.DeleteFileFromLocal(imagePath)
	}

	return nil
}

func (commonService *CommonService) GetSysAdmin() (userModel.User, error) {
	return commonService.commonRepository.GetSysAdmin()
}

func (commonService *CommonService) GetStatus() (commonModel.Status, error) {
	// 获取系统管理员信息
	sysuser, err := commonService.commonRepository.GetSysAdmin()
	if err != nil {
		return commonModel.Status{}, err
	}

	// 获取所有用户状态信息
	var users []commonModel.UserStatus
	allusers, err := commonService.commonRepository.GetAllUsers()
	if err != nil {
		return commonModel.Status{}, err
	}
	for _, user := range allusers {
		users = append(users, commonModel.UserStatus{
			UserID:   user.ID,
			UserName: user.Username,
			IsAdmin:  user.IsAdmin,
		})
	}

	status := commonModel.Status{}

	echos, err := commonService.commonRepository.GetAllEchos(true)
	if err != nil {
		return status, err
	}

	status.SysAdminID = sysuser.ID     // 管理员ID
	status.Username = sysuser.Username // 管理员用户名
	status.Logo = sysuser.Avatar       // 管理员头像
	status.Users = users               // 所有用户状态
	status.TotalEchos = len(echos)     // Echo总数

	return status, nil
}

func (commonService *CommonService) GetHeatMap() ([]commonModel.Heatmap, error) {
	// 获取当前日期
	today := time.Now()

	// 获取一个月前的日期
	oneMonthAgo := today.AddDate(0, -1, 0)

	// 格式化为YYYY-MM-DD
	startDate := oneMonthAgo.Format("2006-01-02") // 一个月前的日期
	endDate := today.Format("2006-01-02")         // 当前日期

	// 数据库查询 （只返回某天count >= 1的item）
	heatmapData, err := commonService.commonRepository.GetHeatMap(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// 如果不足30天，补齐数据（date为缺的日期，count为0）
	// Create a map for quick lookup of existing heatmap data
	heatmapMap := make(map[string]commonModel.Heatmap)
	for _, item := range heatmapData {
		heatmapMap[item.Date] = item
	}

	var results [30]commonModel.Heatmap
	for i := 0; i < 30; i++ {
		// 计算日期 (from today back to 29 days ago)
		date := today.AddDate(0, 0, -i).Format("2006-01-02")
		resultIndex := 29 - i

		if item, ok := heatmapMap[date]; ok {
			// 找到数据，填充结果
			results[resultIndex] = item
		} else {
			// 未找到数据，填充默认值
			results[resultIndex] = commonModel.Heatmap{
				Date:  date,
				Count: 0,
			}
		}
	}

	return results[:], nil
}

func (commonService *CommonService) GenerateRSS(ctx *gin.Context) (string, error) {
	// 获取所有Echo
	echos, err := commonService.commonRepository.GetAllEchos(false)
	if err != nil {
		return "", err
	}

	// 生成 RSS 订阅链接
	schema := "http"
	if ctx.Request.TLS != nil {
		schema = "https"
	}
	host := ctx.Request.Host
	feed := &feeds.Feed{
		Title: "Ech0",
		Link: &feeds.Link{
			Href: fmt.Sprintf("%s://%s/", schema, host),
		},
		Image: &feeds.Image{
			Url: fmt.Sprintf("%s://%s/favicon.ico", schema, host),
		},
		Description: "Ech0",
		Author: &feeds.Author{
			Name: "Ech0",
		},
		Updated: time.Now(),
	}

	for _, msg := range echos {
		renderedContent := mdUtil.MdToHTML([]byte(msg.Content))

		title := msg.Username + " - " + msg.CreatedAt.Format("2006-01-02")

		// 添加图片链接到正文前(scheme://host/api/ImageURL)
		if len(msg.Images) > 0 {
			for _, image := range msg.Images {
				// 根据图片来源生成链接
				switch image.ImageSource {
				case echoModel.ImageSourceLocal:
					imageURL := fmt.Sprintf("%s://%s/api%s", schema, host, image.ImageURL)
					renderedContent = append(
						[]byte(
							fmt.Sprintf(
								"<img src=\"%s\" alt=\"Image\" style=\"max-width:100%%;height:auto;\" />",
								imageURL,
							),
						),
						renderedContent...)
				case echoModel.ImageSourceS3:
					imageURL := image.ImageURL
					renderedContent = append(
						[]byte(
							fmt.Sprintf(
								"<img src=\"%s\" alt=\"Image\" style=\"max-width:100%%;height:auto;\" />",
								imageURL,
							),
						),
						renderedContent...)
				}
			}
		}

		item := &feeds.Item{
			Title:       title,
			Link:        &feeds.Link{Href: fmt.Sprintf("%s://%s/echo/%d", schema, host, msg.ID)},
			Description: string(renderedContent),
			Author: &feeds.Author{
				Name: msg.Username,
			},
			Created: msg.CreatedAt,
		}
		feed.Items = append(feed.Items, item)
	}

	atom, err := feed.ToAtom()
	if err != nil {
		return "", err
	}

	return atom, nil
}

func (commonService *CommonService) UploadMusic(userId uint, file *multipart.FileHeader) (string, error) {
	user, err := commonService.commonRepository.GetUserByUserId(userId)
	if err != nil {
		return "", err
	}
	if !user.IsAdmin {
		return "", errors.New(commonModel.NO_PERMISSION_DENIED)
	}

	// 检查文件类型是否合法
	if !storageUtil.IsAllowedType(file.Header.Get("Content-Type"), config.Config.Upload.AllowedTypes) {
		return "", errors.New(commonModel.FILE_TYPE_NOT_ALLOWED)
	}

	// 检查文件大小是否合法
	if file.Size > int64(config.Config.Upload.AudioMaxSize) {
		return "", errors.New(commonModel.FILE_SIZE_EXCEED_LIMIT)
	}

	// 调用存储函数存储图片
	audioUrl, err := storageUtil.UploadFile(file, commonModel.AudioType, commonModel.LOCAL_FILE)
	if err != nil {
		return "", err
	}

	return audioUrl, nil
}

func (commonService *CommonService) DeleteMusic(userid uint) error {
	user, err := commonService.commonRepository.GetUserByUserId(userid)
	if err != nil {
		return err
	}
	if !user.IsAdmin {
		return errors.New(commonModel.NO_PERMISSION_DENIED)
	}

	// 支持的音频格式
	audioFiles := []string{"music.flac", "music.m4a", "music.mp3"}

	for _, file := range audioFiles {
		audioPath := fmt.Sprintf("data/audios/%s", file)
		if storageUtil.FileExists(audioPath) {
			return storageUtil.DeleteFileFromLocal(audioPath)
		}
	}

	return nil
}

func (commonService *CommonService) GetPlayMusicUrl() string {
	// 支持的音频格式
	audioFiles := []string{"music.flac", "music.m4a", "music.mp3"}

	for _, file := range audioFiles {
		audioPath := fmt.Sprintf("data/audios/%s", file)
		if storageUtil.FileExists(audioPath) {
			return fmt.Sprintf("/audios/%s", file)
		}
	}

	// 没有找到音频文件
	return ""
}

func (commonService *CommonService) PlayMusic(ctx *gin.Context) {
	// 以文件流的形式返回音乐文件
	musicURL := commonService.GetPlayMusicUrl()
	musicName := ""
	if musicURL != "" {
		// 只保留最后的文件名
		musicName = musicURL[len("/audios/"):]
	}

	// 获取音乐文件的路径
	musicPath := config.Config.Upload.AudioPath + musicName

	// 获取 Content-Type
	contentType := "audio/mpeg"
	lowerName := strings.ToLower(musicName)
	switch {
	case strings.HasSuffix(lowerName, ".flac"):
		contentType = "audio/flac"
	case strings.HasSuffix(lowerName, ".wav"):
		contentType = "audio/wav"
	}

	// 读取文件内容
	data, err := os.ReadFile(musicPath)
	if err != nil {
		ctx.String(500, "读取音乐文件失败")
		return
	}

	// 设置响应头
	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Header("Pragma", "no-cache")
	ctx.Header("Expires", "0")

	// 直接写文件内容，Gin 会自动关闭连接，不会长时间占用文件
	ctx.Data(http.StatusOK, contentType, data)
}

// GetS3PresignURL 获取 S3 预签名 URL
func (commonService *CommonService) GetS3PresignURL(
	userid uint,
	s3Dto *commonModel.GetPresignURLDto,
	method string,
) (commonModel.PresignDto, error) {
	var result commonModel.PresignDto

	user, err := commonService.commonRepository.GetUserByUserId(userid)
	if err != nil {
		return result, err
	}
	if !user.IsAdmin {
		return result, errors.New(commonModel.NO_PERMISSION_DENIED)
	}

	if s3Dto.FileName == "" {
		return result, errors.New(commonModel.INVALID_PARAMS)
	}
	ext := filepath.Ext(s3Dto.FileName) // ".png"
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// 检查Content-Type是否为Image开头
	switch contentType[:5] {
	case "image":
		// 检查文件类型是否合法
		if !storageUtil.IsAllowedType(contentType, config.Config.Upload.AllowedTypes) {
			return result, errors.New(commonModel.FILE_TYPE_NOT_ALLOWED)
		}
	case "audio":
		// 检查文件类型是否合法
		if !storageUtil.IsAllowedType(contentType, config.Config.Upload.AllowedTypes) {
			return result, errors.New(commonModel.FILE_TYPE_NOT_ALLOWED)
		}
	default:
		return result, errors.New(commonModel.FILE_TYPE_NOT_ALLOWED)
	}

	result.FileName = s3Dto.FileName
	result.ContentType = contentType

	// 获取 S3 配置和客户端
	_, s3setting, err := commonService.GetS3Client()
	if err != nil {
		return result, err
	}

	// 检查是否开启了 S3
	if !s3setting.Enable {
		return result, errors.New(commonModel.S3_NOT_ENABLED)
	}

	// 生成 Object Key (包含 PathPrefix)
	prefix := strings.TrimSuffix(s3setting.PathPrefix, "/")
	objectKey := fmt.Sprintf("%s/%s_%d", prefix, s3Dto.FileName, time.Now().Unix())
	result.ObjectKey = objectKey

	// 生成预签名 URL
	presignURL, err := commonService.objStorage.PresignURL(context.Background(), objectKey, 24*time.Hour, method)
	if err != nil {
		return result, err
	}
	result.PresignURL = presignURL
	// 生成访问 URL
	fileURL, err := commonService.GetS3ObjectURL(s3setting, objectKey)
	if err != nil {
		return result, err
	}
	result.FileURL = fileURL

	// 保存到临时文件表
	now := time.Now().Unix()
	tempFile := commonModel.TempFile{
		FileName:       result.FileName,
		Storage:        string(commonModel.S3_FILE),
		FileType:       string(commonModel.ImageType),
		Bucket:         s3setting.BucketName,
		ObjectKey:      result.ObjectKey,
		Deleted:        false,
		CreatedAt:      now,
		LastAccessedAt: now,
	}
	commonService.txManager.Run(func(ctx context.Context) error {
		return commonService.commonRepository.SaveTempFile(ctx, tempFile)
	})

	return result, nil
}

// GetS3Client 获取 S3 客户端和配置信息
func (commonService *CommonService) GetS3Client() (storageUtil.ObjectStorage, settingModel.S3Setting, error) {
	// 检查是否配置了 S3
	var s3setting settingModel.S3Setting
	value, err := commonService.keyvalueRepository.GetKeyValue(commonModel.S3SettingKey)
	if err != nil || value == "" {
		return nil, s3setting, errors.New(commonModel.S3_NOT_CONFIGURED)
	}
	if err := jsonUtil.JSONUnmarshal([]byte(value.(string)), &s3setting); err != nil {
		return nil, s3setting, errors.New(commonModel.S3_CONFIG_ERROR)
	}
	s3setting.Endpoint = httpUtil.TrimURL(s3setting.Endpoint)

	// 初始化 S3 客户端
	commonService.objStorage, err = storageUtil.NewMinioStorage(
		s3setting.Endpoint,
		s3setting.AccessKey,
		s3setting.SecretKey,
		s3setting.BucketName,
		s3setting.UseSSL,
	)
	if err != nil {
		return nil, s3setting, errors.New(commonModel.S3_CONFIG_ERROR)
	}

	return commonService.objStorage, s3setting, nil
}

// GetS3ObjectURL 获取 S3 对象的 URL
func (CommonService *CommonService) GetS3ObjectURL(s3Setting settingModel.S3Setting, objectKey string) (string, error) {
	if s3Setting.Endpoint == "" || s3Setting.BucketName == "" || objectKey == "" {
		return "", errors.New(commonModel.S3_CONFIG_ERROR)
	}

	protocal := "http"
	if s3Setting.UseSSL {
		protocal = "https"
	}

	return fmt.Sprintf("%s://%s/%s/%s", protocal, s3Setting.Endpoint, s3Setting.BucketName, objectKey), nil
}

// CleanupTempFiles 清理过期的临时文件
func (commonService *CommonService) CleanupTempFiles() error {
	// 获取所有未删除的临时文件
	files, err := commonService.commonRepository.GetAllTempFiles()
	if err != nil {
		return err
	}

	// 当前时间戳
	now := time.Now().Unix()

	for _, file := range files {
		// 如果最后访问时间超过24小时，则删除
		if now-file.LastAccessedAt > 24*3600 {
			// 删除文件
			switch file.Storage {
			case string(commonModel.LOCAL_FILE):
				// TODO: 删除本地文件

			case string(commonModel.S3_FILE):
				// 获取 S3 客户端
				cli, _, err := commonService.GetS3Client()
				if err != nil {
					// 如果没有配置 S3，则无法删除,忽略
					continue
				}
				if file.ObjectKey == "" {
					// 如果没有传入 object_key，则无法删除,忽略
					continue
				}
				// 删除 S3 上的文件
				if err := cli.DeleteObject(context.Background(), file.ObjectKey); err != nil {
					// 记录日志，继续处理下一个文件
					fmt.Printf("删除S3临时文件失败: %s, 错误: %v\n", file.ObjectKey, err)
				}
			default:
				// 未知存储类型，忽略
			}

			// 从数据库中删除记录(开启事务)
			commonService.txManager.Run(func(ctx context.Context) error {
				return commonService.commonRepository.DeleteTempFilePermanently(ctx, file.ID)
			})
		}
	}

	return nil
}

func (commonService *CommonService) RefreshEchoImageURL(echo *echoModel.Echo) {
	_, s3setting, err := commonService.GetS3Client()
	if err != nil {
		return
	}

	// 用 channel 或 waitGroup 并发刷新 URL
	var wg sync.WaitGroup
	mu := sync.Mutex{}

	for i := range echo.Images {
		if echo.Images[i].ImageSource == echoModel.ImageSourceS3 && echo.Images[i].ObjectKey != "" {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				if newURL, err := commonService.GetS3ObjectURL(s3setting, echo.Images[i].ObjectKey); err == nil {
					mu.Lock()
					echo.Images[i].ImageURL = newURL
					mu.Unlock()
				}
			}(i)
		}
	}

	wg.Wait()

	// 所有 URL 都拿到了，再一次性更新 DB
	_ = commonService.txManager.Run(func(ctx context.Context) error {
		return commonService.echoRepository.UpdateEcho(ctx, echo)
	})
}
