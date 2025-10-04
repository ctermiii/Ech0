<template>
  <div class="max-w-sm px-3 pb-4 py-2 mt-4 sm:mt-6 mb-10 mx-auto flex justify-center items-center">
    <div class="w-full sm:max-w-lg">
      <h1 class="text-5xl text-center font-bold text-gray-200 mt-2 mb-4">Ech0 Fediverse</h1>

      <div class="flex items-center justify-between">
        <!-- 返回上一页 -->
        <BaseButton @click="goBack" class="w-10 h-10 text-gray-600 rounded-md" title="返回首页">
          <Arrow class="w-7 h-7 rotate-180 mx-auto" />
        </BaseButton>

        <!-- Actor搜索框 && NotificationBox -->
        <div class="flex items-center gap-1">
          <!-- Actor 搜索框 -->
          <BaseInput
            title="搜索"
            type="text"
            v-model="searchTerm"
            placeholder="搜索 Actor..."
            class="w-50 sm:w-55 h-10"
            @keyup.enter="$event.target.blur()"
            @blur="handleSearch"
          />
          <!-- NotificationBox -->
          <BaseButton
            class="h-full w-full text-gray-600 rounded-md flex items-center justify-center"
            title="消息通知"
            :icon="InBox"
            @click="theToast.info('消息通知功能开发中，敬请期待！')"
          />
        </div>
      </div>

      <!-- 在搜索时板块显示搜索结果 -->
      <div v-if="shouldShowResults" class="mt-6 space-y-4">
        <p v-if="searchLoading" class="text-sm text-gray-400">正在召唤联邦宇宙的朋友们…</p>
        <TheActorCard v-else-if="searchResult" :actor="searchResult" />
      </div>
      <!-- 未搜索时显示已关注的 Actor 的动态 -->
      <div
        v-else
        class="mt-6 rounded-lg border border-dashed border-gray-700/60 px-4 py-8 text-center text-gray-500"
      >
        功能开发中，敬请期待！
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, computed, watch } from 'vue'
import { theToast } from '@/utils/toast'

import { useRouter } from 'vue-router'
import BaseButton from '@/components/common/BaseButton.vue'
import BaseInput from '@/components/common/BaseInput.vue'
import Arrow from '@/components/icons/arrow.vue'
import InBox from '@/components/icons/inbox.vue'
import { fetchSearchFediverseActor } from '@/service/api/fediverse'
import TheActorCard from '@/components/advanced/TheActorCard.vue'
import { useUserStore } from '@/stores/user'
import { storeToRefs } from 'pinia'

const userStore = useUserStore()
const { isLogin } = storeToRefs(userStore)
const router = useRouter()

// 返回首页
const goBack = () => {
  if (window.history.length > 2) {
    window.history.back()
  } else {
    router.push({ name: 'home' }) // 没有历史记录则跳首页
  }
}

//================================================
// 搜索相关
//================================================

const searchTerm = ref('')
const hasSearched = ref(false)
const searchLoading = ref(false)

const searchResult = ref<App.Api.Fediverse.Actor | null>(null)

// 是否显示搜索结果区域
const shouldShowResults = computed(() => hasSearched.value || searchLoading.value)

// 监听搜索词变化，清除状态
watch(
  () => searchTerm.value,
  (value) => {
    if (!value.trim()) {
      hasSearched.value = false
      searchLoading.value = false
      searchResult.value = null
    }
  },
)

// 处理搜索
const handleSearch = async (event?: KeyboardEvent | MouseEvent) => {
  if (!isLogin.value) {
    theToast.error('请先登录以使用联邦宇宙功能')
    return
  }

  if (event && 'target' in event) {
    const target = event.target as HTMLElement | null
    target?.blur()
  }

  // 去除前后空格
  const term = searchTerm.value.trim()
  if (!term) {
    hasSearched.value = false
    searchLoading.value = false
    searchResult.value = null
    return
  }

  // 重置状态
  hasSearched.value = true
  searchLoading.value = true
  searchResult.value = null

  // 调用搜索接口
  const response = await fetchSearchFediverseActor(term)
  searchLoading.value = false

  if (response.code === 1 && response.data) {
    searchResult.value = response.data
  } else {
    searchResult.value = null
    searchLoading.value = false
  }
}

onMounted(() => {
  theToast.info('欢迎来到联邦宇宙！🎉', { duration: 3000 })
})
</script>
