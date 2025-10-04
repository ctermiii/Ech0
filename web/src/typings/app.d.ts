declare namespace App {
  /**
   * Namespace Api
   */
  namespace Api {
    type Response<T> = {
      code: number
      msg: string
      data: T
    }

    namespace Auth {
      type LoginParams = {
        username: string
        password: string
      }

      type LoginResponse = {
        token: string
      }

      type SignupParams = {
        username: string
        password: string
      }
    }

    namespace User {
      type User = {
        id: number
        username: string
        password?: string
        is_admin: boolean
        avatar?: string
      }

      type UserInfo = {
        username: string
        password: string
        is_admin: boolean
        avatar: string
      }

      type UserStatus = {
        user_id: number
        username: string
        is_admin: boolean
      }
    }

    namespace Ech0 {
      type ParamsByPagination = {
        page: number
        pageSize: number
        search?: string
      }

      type Echo = {
        id: number
        content: string
        username: string
        image_url: string
        image_source: string
        images: Image[]
        private: boolean
        user_id: number
        extension: string
        extension_type: string
        fav_count: number
        created_at: string
      }

      type Image = {
        id: number
        message_id: number
        image_url: string
        image_source: string
        object_key?: string // 对象存储的Key (如果是本地存储则为空)
      }

      type ImageToAdd = {
        image_url: string
        image_source: string
        object_key?: string // 对象存储的Key (如果是本地存储则为空)
      }

      type EchoToAdd = {
        content: string
        image_url?: string | null // 新版web已经废弃该字段
        image_source?: string | null // 新版web已经废弃该字段
        images?: ImageToAdd[] | null
        extension?: string | null
        extension_type?: string | null
        private: boolean
      }

      type EchoToUpdate = {
        id: number
        content: string
        username: string
        image_url: string
        image_source: string
        images?: ImageToAdd[] | null
        private: boolean
        user_id: number
        extension?: string | null
        extension_type?: string | null
        created_at: string
      }

      type PaginationResult = {
        items: Echo[]
        total: number
      }

      type Status = {
        sys_admin_id: number
        username: string
        logo: string
        users: App.Api.User.UserStatus[]
        total_echos: number
      }

      type HeatMap = {
        date: string
        count: number
      }[]

      type ImageToDelete = {
        url: string
        source: string
        object_key?: string // 对象存储的 Key, 用于删除 S3/R2 上的图片
      }

      type GithubCardData = {
        name: string
        stargazers_count: number
        forks_count: number
        description: string
        owner: {
          avatar_url: string
        }
      }

      type HelloEch0 = {
        hello: string
        version: string
        github: string
      }

      type PresignResult = {
        file_name: string
        content_type: string
        object_key: string
        presign_url: string
        file_url: string
      }
    }

    namespace Setting {
      type SystemSetting = {
        site_title: string
        server_name: string
        server_url: string
        allow_register: boolean
        ICP_number: string
        meting_api: string
        custom_css: string
        custom_js: string
      }

      type CommentSetting = {
        enable_comment: boolean
        provider: string // 评论提供者
        comment_api: string // 评论 API 地址
      }

      type S3Setting = {
        enable: boolean
        provider: string
        endpoint: string
        access_key: string
        secret_key: string
        bucket_name: string
        region: string
        use_ssl: boolean
        cdn_url: string
        path_prefix: string
        public_read: boolean
      }
    }

    namespace Connect {
      type Connect = {
        server_name: string
        server_url: string
        logo: string
        total_echos: number
        today_echos: number
        sys_username: string
      }

      type Connected = {
        id: number
        connect_url: string
      }
    }

    namespace Fediverse {
      type Actor = Record<string, unknown>

      type FollowActionRequest = {
        targetActor: string
      }

      type FollowResponse = {
        activityId: string
      }

      // type LikeActionRequest = {
      //   targetActor: string
      //   object: string
      //   objectType?: string
      // }

      // type UnfollowResponse = {
      //   activityId: string
      //   followActivityId?: string
      // }

      // type LikeResponse = {
      //   activityId: string
      // }

      // type UndoLikeResponse = {
      //   activityId: string
      //   likeActivityId?: string
      // }
    }

    namespace Todo {
      type Todo = {
        id: number
        content: string
        user_id: number
        username: string
        status: number
        created_at: string
      }

      type TodoToAdd = {
        content: string
      }
    }
  }
}
