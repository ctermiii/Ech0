import { request } from '@/service/request'

// Search Actor By Actor ID
export function fetchSearchFediverseActor(actor: string) {
  return request<App.Api.Fediverse.Actor>({
    url: `/search/actor?actor=${encodeURIComponent(actor)}`,
    method: 'GET',
  })
}

// Get Follow Status (获取关注状态)
export function fetchGetFollowStatus(targetActor: string) {
  return request<string>({
    url: `/follow/status?actor=${encodeURIComponent(targetActor)}`,
    method: 'GET',
  })
}

// Follow (发起关注请求)
export function fetchFollowFediverseActor(payload: App.Api.Fediverse.FollowActionRequest) {
  return request<App.Api.Fediverse.FollowResponse>({
    url: `/follow`,
    method: 'POST',
    data: payload,
  })
}

// // Unfollow (取消关注请求)
// export function fetchUnfollowFediverseActor(payload: App.Api.Fediverse.FollowActionRequest) {
//   return request<App.Api.Fediverse.UnfollowResponse>({
//     url: `/unfollow`,
//     method: 'POST',
//     data: payload,
//   })
// }

// // Post Like (点赞请求)
// export function fetchLikeFediverseObject(payload: App.Api.Fediverse.LikeActionRequest) {
//   return request<App.Api.Fediverse.LikeResponse>({
//     url: `/like`,
//     method: 'POST',
//     data: payload,
//   })
// }

// // Post Undo Like (取消点赞请求)
// export function fetchUndoLikeFediverseObject(payload: App.Api.Fediverse.LikeActionRequest) {
//   return request<App.Api.Fediverse.UndoLikeResponse>({
//     url: `/undo-like`,
//     method: 'POST',
//     data: payload,
//   })
// }
