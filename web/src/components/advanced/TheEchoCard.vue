<template>
  <div class="w-full">
    <!-- 日期时间 && 操作按钮 -->
    <div class="flex justify-between items-center">
      <!-- 日期时间 -->
      <div @click="handleExpandEcho(echo.id)" class="flex justify-start items-center h-auto">
        <!-- 小点 -->
        <div class="w-2 h-2 rounded-full bg-orange-600 mr-2"></div>
        <!-- 具体日期时间 -->
        <div
          class="flex justify-start text-sm text-orange-500 hover:underline hover:decoration-offset-3 hover:decoration-1"
        >
          {{ formatDate(props.echo.created_at) }}
        </div>
      </div>

      <!-- 操作按钮 -->
      <div ref="menuRef" class="relative flex items-center justify-center gap-1 h-auto">
        <!-- 更多操作 -->
        <div
          v-if="!showMenu"
          @click.stop="toggleMenu"
          class="w-7 h-7 flex items-center justify-center bg-white ring-1 ring-gray-200 ring-inset rounded-full shadow-sm hover:shadow-md transition"
        >
          <!-- 默认图标，展开后隐藏 -->
          <More class="w-5 h-5" />
        </div>

        <!-- 展开后的按钮组 -->
        <div
          v-if="showMenu"
          class="flex items-center gap-4 bg-white rounded-full px-2 py-1 shadow-sm hover:shadow-md ring-1 ring-gray-200 ring-inset"
        >
          <!-- 是否隐私 -->
          <span v-if="props.echo.private" title="私密状态">
            <Lock />
          </span>

          <!-- 删除 -->
          <button
            v-if="userStore.isLogin"
            @click="handleDeleteEcho(props.echo.id)"
            title="删除"
            class="transform transition-transform duration-200 hover:scale-120"
          >
            <Roll />
          </button>

          <!-- 更新 -->
          <button
            v-if="userStore.isLogin"
            @click="handleUpdateEcho()"
            title="更新"
            class="transform transition-transform duration-200 hover:scale-120"
          >
            <EditEcho />
          </button>

          <!-- 展开内容 -->

          <button
            @click="handleExpandEcho(echo.id)"
            title="展开Echo"
            class="transform transition-transform duration-200 hover:scale-120"
          >
            <Expand />
          </button>

          <!-- 点赞 -->
          <div class="flex items-center justify-end" title="点赞">
            <div class="flex items-center gap-1">
              <!-- 点赞按钮   -->
              <button
                @click="handleLikeEcho(props.echo.id)"
                title="点赞"
                class="transform transition-transform duration-200 hover:scale-120"
              >
                <GrayLike class="w-4 h-4 transition-colors duration-200 hover:text-red-500" />
              </button>

              <!-- 点赞数量   -->
              <span class="text-sm text-gray-400">
                <!-- 如果点赞数不超过99，则显示数字，否则显示99+ -->
                {{ props.echo.fav_count > 99 ? '99+' : props.echo.fav_count }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 图片 && 内容 -->
    <div class="border-l-2 border-[#0000000d] ml-1">
      <div class="px-4 py-3">
        <TheImageGallery :images="props.echo.images" />

        <!-- 内容 -->
        <div class="mx-auto w-11/12 pl-1">
          <MdPreview
            :id="previewOptions.proviewId"
            :modelValue="props.echo.content"
            :theme="previewOptions.theme"
            :show-code-row-number="previewOptions.showCodeRowNumber"
            :preview-theme="previewOptions.previewTheme"
            :code-theme="previewOptions.codeTheme"
            :code-style-reverse="previewOptions.codeStyleReverse"
            :no-img-zoom-in="previewOptions.noImgZoomIn"
            :code-foldable="previewOptions.codeFoldable"
            :auto-fold-threshold="previewOptions.autoFoldThreshold"
          />
        </div>

        <!-- 扩展内容 -->
        <div v-if="props.echo.extension" class="my-2">
          <div v-if="props.echo.extension_type === ExtensionType.MUSIC">
            <TheAPlayerCard :echo="props.echo" />
          </div>
          <div v-if="props.echo.extension_type === ExtensionType.VIDEO">
            <TheVideoCard :videoId="props.echo.extension" class="px-2 mx-auto hover:shadow-md" />
          </div>
          <TheGithubCard
            v-if="props.echo.extension_type === ExtensionType.GITHUBPROJ"
            :GithubURL="props.echo.extension"
            class="px-2 mx-auto hover:shadow-md"
          />
          <TheWebsiteCard
            v-if="props.echo.extension_type === ExtensionType.WEBSITE"
            :website="props.echo.extension"
            class="px-2 mx-auto hover:shadow-md"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { MdPreview } from 'md-editor-v3'
import { onMounted, ref, onBeforeUnmount } from 'vue'
import { fetchDeleteEcho, fetchLikeEcho } from '@/service/api'
import { theToast } from '@/utils/toast'
import { useUserStore } from '@/stores/user'
import TheGithubCard from './TheGithubCard.vue'
import TheVideoCard from './TheVideoCard.vue'
import TheImageGallery from './TheImageGallery.vue'
import 'md-editor-v3/lib/preview.css'
import Roll from '../icons/roll.vue'
import Lock from '../icons/lock.vue'
import More from '../icons/more.vue'
import Expand from '../icons/expand.vue'
import GrayLike from '../icons/graylike.vue'
import EditEcho from '../icons/editecho.vue'
import TheAPlayerCard from './TheAPlayerCard.vue'
import TheWebsiteCard from './TheWebsiteCard.vue'
import { useEchoStore } from '@/stores/echo'
import { localStg } from '@/utils/storage'
import { useRouter } from 'vue-router'
import { ExtensionType } from '@/enums/enums'
import { formatDate } from '@/utils/other'

const emit = defineEmits(['refresh', 'updateLikeCount'])

type Echo = App.Api.Ech0.Echo

const props = defineProps<{
  echo: Echo
}>()

const userStore = useUserStore()
const previewOptions = {
  proviewId: 'preview-only',
  theme: 'light' as 'light' | 'dark',
  showCodeRowNumber: false,
  previewTheme: 'github',
  codeTheme: 'atom',
  codeStyleReverse: true,
  noImgZoomIn: false,
  codeFoldable: true,
  autoFoldThreshold: 15,
}

const echoStore = useEchoStore()
const router = useRouter()

const handleDeleteEcho = (echoId: number) => {
  // 浏览器alert弹窗确认删除
  if (confirm('确定要删除吗？')) {
    fetchDeleteEcho(echoId).then(() => {
      theToast.success('删除成功！')
      // 触发父组件的刷新事件emit
      emit('refresh')
    })
  }
}

const handleUpdateEcho = () => {
  if (echoStore.isUpdateMode) {
    // 如果已经在更新模式，返回顶部并提示用户先退出更新模式
    window.scrollTo({ top: 0, behavior: 'smooth' })
    theToast.warning('请先退出更新模式！')
    return
  }

  echoStore.echoToUpdate = props.echo
  echoStore.isUpdateMode = true
}

const LIKE_LIST_KEY = 'likedEchoIds'
const likedEchoIds: number[] = localStg.getItem(LIKE_LIST_KEY) || []
const hasLikedEcho = (echoId: number): boolean => {
  return likedEchoIds.includes(echoId)
}
const handleLikeEcho = (echoId: number) => {
  // 检查LocalStorage中是否已经点赞过
  if (hasLikedEcho(echoId)) {
    theToast.success('你已经点赞过了,感谢你的喜欢！')
    return
  }

  fetchLikeEcho(echoId).then((res) => {
    if (res.code === 1) {
      likedEchoIds.push(echoId)
      localStg.setItem(LIKE_LIST_KEY, likedEchoIds)
      // 发送更新事件
      emit('updateLikeCount', echoId)
      theToast.success('点赞成功！')
    }
  })
}

const handleExpandEcho = (echoId: number) => {
  // 跳转到Echo详情
  router.push({
    name: 'echo',
    params: { echoId: echoId },
  })
}

const showMenu = ref(false)
const menuRef = ref<HTMLElement | null>(null)

const toggleMenu = () => {
  showMenu.value = !showMenu.value
}

const handleClickOutside = (event: MouseEvent) => {
  if (menuRef.value && !menuRef.value.contains(event.target as Node)) {
    showMenu.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped lang="css">
#preview-only {
  background-color: inherit;
}

:deep(.md-editor) {
  /* font-family: var(--font-sans); */
  font-family: 'LXGW WenKai Screen';
}

:deep(.md-editor div.github-theme) {
  line-height: 1.6;
  color: #000;
}

:deep(ul li) {
  list-style-type: disc;
}

:deep(ul li li) {
  list-style-type: circle;
}

:deep(ul li li li) {
  list-style-type: square;
}

:deep(ol li) {
  list-style-type: decimal;
}

:deep(p) {
  white-space: normal; /* 允许正常换行 */
  overflow-wrap: break-word; /* 单词太长时自动换行 */
  word-break: normal; /* 保持单词整体性，不随便拆开 */
}
</style>
