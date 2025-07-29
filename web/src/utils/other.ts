import { getApiUrl } from '@/service/request/shared'
import { ImageSource } from '@/enums/enums'
// 获取图片链接
export const getImageUrl = (image: App.Api.Ech0.Image) => {
  if (image.image_source === ImageSource.LOCAL) {
    return getApiUrl() + String(image.image_url)
  } else if (image.image_source === ImageSource.URL) {
    return String(image.image_url)
  } else if (image.image_source === ImageSource.S3) {
    return String(image.image_url)
  } else {
    // 未知的图片来源，按照本地图片处理
    return getApiUrl() + String(image.image_url)
  }
}

// 获取待添加图片链接
export const getImageToAddUrl = (image: App.Api.Ech0.ImageToAdd) => {
  if (image.image_source === ImageSource.LOCAL) {
    return getApiUrl() + String(image.image_url)
  } else if (image.image_source === ImageSource.URL) {
    return String(image.image_url)
  } else if (image.image_source === ImageSource.S3) {
    return String(image.image_url)
  } else {
    // 未知的图片来源，按照本地图片处理
    return getApiUrl() + String(image.image_url)
  }
}

export const formatDate = (dateString: string) => {
  // 当天则显示（时：分）
  // 非当天但是三内天则显示几天前
  // 超过三天则显示（时：分 年月日）
  const date = new Date(dateString)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const diffInDays = Math.floor(diff / (1000 * 60 * 60 * 24))
  const diffInHours = Math.floor(diff / (1000 * 60 * 60))
  const diffInMinutes = Math.floor(diff / (1000 * 60))

  const diffInSeconds = Math.floor(diff / 1000)
  if (diffInSeconds < 60) {
    return '刚刚'
  } else if (diffInMinutes < 60) {
    return `${diffInMinutes}分钟前`
  } else if (diffInHours < 24) {
    return `${diffInHours}小时前`
  } else if (diffInDays < 3) {
    return `${diffInDays}天前`
  } else {
    return date.toLocaleString() // 返回完整的日期和时间
  }
}

// 解析音乐链接（网易云、QQ音乐）
export const parseMusicURL = (url: string) => {
  url = url.trim()

  const neteaseMatch = url.match(/music\.163\.com\/(#\/)?(song|playlist|album)(\?id=|\/)(\d+)/)
  if (neteaseMatch) {
    return {
      server: 'netease',
      type: neteaseMatch[1], // song, playlist, album
      id: neteaseMatch[2],
    }
  }

  // QQ音乐 新格式支持，songDetail 路径，id一般是字母数字混合
  const qqNewSongMatch = url.match(/y\.qq\.com\/n\/ryqq\/songDetail\/([a-zA-Z0-9]+)/)
  if (qqNewSongMatch) {
    return {
      server: 'tencent',
      type: 'song',
      id: qqNewSongMatch[1],
    }
  }

  // 解析失败
  return null
}
