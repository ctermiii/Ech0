<template>
  <div
    class="bg-white rounded-lg ring-1 ring-gray-200 ring-inset mx-auto shadow-sm hover:shadow-md"
  >
    <div class="mx-auto w-full px-3 py-4">
      <!-- The Title && Nav -->
      <TheTitleAndNav />

      <!-- The Editor -->
      <div class="rounded-lg p-2 sm:p-3 mb-1">
        <TheMdEditor
          v-model="echoToAdd.content"
          class="rounded-lg"
          v-if="currentMode === Mode.ECH0"
        />

        <!-- ImageMode : TheImageEditor -->
        <div v-if="currentMode === Mode.Image">
          <h2 class="text-gray-500 font-bold my-2">插入图片（支持直链、本地、S3存储）</h2>
          <div v-if="!ImageUploading" class="flex items-center justify-between mb-3">
            <div class="flex items-center gap-2">
              <span class="text-gray-500">选择添加方式：</span>
              <!-- 直链 -->
              <BaseButton
                :icon="Url"
                class="w-7 h-7 sm:w-7 sm:h-7 rounded-md"
                @click="imageToAdd.image_source = ImageSource.URL"
                title="插入图片链接"
              />
              <!-- 上传本地 -->
              <BaseButton
                :icon="Upload"
                class="w-7 h-7 sm:w-7 sm:h-7 rounded-md"
                @click="imageToAdd.image_source = ImageSource.LOCAL"
                title="上传本地图片"
              />
              <!-- S3 存储 -->
              <BaseButton
                v-if="S3Setting.enable"
                :icon="Bucket"
                class="w-7 h-7 sm:w-7 sm:h-7 rounded-md"
                @click="imageToAdd.image_source = ImageSource.S3"
                title="S3存储图片"
              />
            </div>
            <div>
              <BaseButton
                v-if="imageToAdd.image_url != ''"
                :icon="Addmore"
                class="w-7 h-7 sm:w-7 sm:h-7 rounded-md"
                @click="handleAddMoreImage"
                title="添加更多图片"
              />
            </div>
          </div>

          <div class="text-gray-300 text-sm mb-1">
            当前上传方式为
            <span class="font-bold">
              {{
                imageToAdd.image_source === ImageSource.URL
                  ? '直链'
                  : imageToAdd.image_source === ImageSource.LOCAL
                    ? '本地存储'
                    : 'S3存储'
              }}</span
            >
            {{ !ImageUploading ? '' : '，正在上传中...' }}
          </div>

          <div class="my-1">
            <!-- 图片上传本地 -->
            <!-- <input
              id="file-input"
              class="hidden"
              type="file"
              accept="image/*"
              ref="fileInput"
              @change="handleUploadImage"
            />
            <BaseButton
              v-if="imageToAdd.image_source === ImageSource.LOCAL"
              @click="handleTriggerUpload"
              class="rounded-md"
              title="上传本地图片"
            >
              <span class="text-gray-400">点击上传</span>
            </BaseButton> -->

            <TheUppy
              v-if="imageToAdd.image_source !== ImageSource.URL"
              @uppyUploaded="handleUppyUploaded"
              @uppy-set-image-source="handleSetImageSource"
              :TheImageSource="imageToAdd.image_source"
            />

            <!-- 图片直链 -->
            <BaseInput
              v-if="imageToAdd.image_source === ImageSource.URL"
              v-model="imageToAdd.image_url"
              class="rounded-lg h-auto w-full"
              placeholder="请输入图片链接..."
            />
          </div>
        </div>

        <!-- TodoMode : TheTodoModeEditor -->
        <TheTodoModeEditor :current-mode="currentMode" v-model:content="todoToAdd.content" />

        <!-- MusicMode : TheMusicModeEditor -->
        <TheMusicModeEditor @refresh-audio="handleRefreshAudio" :current-mode="currentMode" />

        <!-- The Mode Panel -->
        <TheModePanel
          v-if="currentMode === Mode.Panel"
          v-model:mode="currentMode"
          v-model:extension-type="currentExtensionType"
          v-model:extension-to-add="extensionToAdd"
        />

        <!-- ExtensionMode: TheExtensionEditor -->
        <TheExtensionEditor
          v-model:current-mode="currentMode"
          v-model:current-extension-type="currentExtensionType"
          v-model:extension-to-add="extensionToAdd"
          v-model:video-u-r-l="videoURL"
          v-model:website-to-add="websiteToAdd"
        />
      </div>

      <!-- Editor Buttons -->
      <TheEditorButtons
        :echo-to-add="echoToAdd"
        :current-mode="currentMode"
        @handle-addor-update="handleAddorUpdate"
        @handle-change-mode="handleChangeMode"
        @handle-add-image-mode="handleAddImageMode"
        @handle-exit-update-mode="handleExitUpdateMode"
        @handle-private="handlePrivate"
      />

      <!-- Editor Image -->
      <TheEditorImage
        :imagesToAdd="imagesToAdd"
        :current-mode="currentMode"
        @handleAddorUpdateEcho="handleAddorUpdateEcho"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import Upload from '@/components/icons/upload.vue'
import Url from '@/components/icons/url.vue'
import Addmore from '@/components/icons/addmore.vue'
import Bucket from '@/components/icons/bucket.vue'

import BaseButton from '@/components/common/BaseButton.vue'
import BaseInput from '@/components/common/BaseInput.vue'

import TheMdEditor from '@/components/advanced/TheMdEditor.vue'
import TheModePanel from './TheEditor/TheModePanel.vue'
import TheTitleAndNav from './TheEditor/TheTitleAndNav.vue'
import TheEditorImage from './TheEditor/TheEditorImage.vue'
import TheEditorButtons from './TheEditor/TheEditorButtons.vue'
import TheTodoModeEditor from './TheEditor/TheTodoModeEditor.vue'
import TheMusicModeEditor from './TheEditor/TheMusicModeEditor.vue'
import TheExtensionEditor from './TheEditor/TheExtensionEditor.vue'
import TheUppy from '@/components/advanced/TheUppy.vue'

import { theToast } from '@/utils/toast'
import { ref, watch } from 'vue'
import { fetchAddEcho, fetchAddTodo, fetchUpdateEcho } from '@/service/api'
import { useEchoStore } from '@/stores/echo'
import { useTodoStore } from '@/stores/todo'
import { useSettingStore } from '@/stores/setting'
import { storeToRefs } from 'pinia'
import { Mode, ExtensionType, ImageSource } from '@/enums/enums'
import { useEditorStore } from '@/stores/editor'

/* --------------- 与音乐播放相关 ---------------- */
const emit = defineEmits(['refreshAudio'])
const handleRefreshAudio = () => {
  emit('refreshAudio')
}
/* ----------------------------------------------- */

/* --------------- 与Pinia相关 ---------------- */
const echoStore = useEchoStore()
const todoStore = useTodoStore()
const settingStore = useSettingStore()
const editorStore = useEditorStore()

const { setTodoMode, getTodos } = todoStore

const { todoMode } = storeToRefs(todoStore)
const { echoToUpdate, isUpdateMode } = storeToRefs(echoStore)
const { S3Setting } = storeToRefs(settingStore)
const { ImageUploading } = storeToRefs(editorStore)

/* -------------------------------------------- */

/* --------------- 与模式和扩展类型相关 ---------------- */
const currentMode = ref<Mode>(Mode.ECH0)
const currentExtensionType = ref<ExtensionType>()

const handleChangeMode = () => {
  if (currentMode.value === Mode.ECH0) {
    currentMode.value = Mode.Panel
  } else if (
    currentMode.value === Mode.TODO ||
    currentMode.value === Mode.PlayMusic ||
    currentMode.value === Mode.EXTEN
  ) {
    currentMode.value = Mode.Panel
    setTodoMode(false)

    if (!echoToAdd.value.image_url || echoToAdd.value.image_url.length === 0) {
      echoToAdd.value.image_source = null
    }
  } else {
    currentMode.value = Mode.ECH0
  }
}
/* ----------------------------------------------------- */

/* --------------- 与各种编辑器输入相关的变量 ---------------- */
const isSubmitting = ref(false)
// 临时网站链接变量
const websiteToAdd = ref<{
  title: string
  site: string
}>({
  title: '',
  site: '',
})
// 临时的视频分享链接变量
const videoURL = ref<string>('')
// 临时扩展变量
const extensionToAdd = ref({
  extension: '',
  extension_type: '',
})
// 临时图片索引变量
const imageIndex = ref<number>(0)
// 临时图片来源变量
// const imageSourceMemory = ref<string>()
// 临时图片添加变量
const imageToAdd = ref<App.Api.Ech0.ImageToAdd>({
  image_url: '',
  image_source: ImageSource.LOCAL,
  object_key: '',
})
// 临时的多张图片数组变量
const imagesToAdd = ref<App.Api.Ech0.ImageToAdd[]>([])
// 最终的Echo添加变量
const echoToAdd = ref<App.Api.Ech0.EchoToAdd>({
  content: '',
  image_url: null,
  image_source: null,
  images: [],
  private: false,
  extension: null,
  extension_type: null,
})
// 临时的Todo添加变量
const todoToAdd = ref<App.Api.Todo.TodoToAdd>({
  content: '',
})
/* ----------------------------------------------------------- */

/* --------------- 与图片相关的各种函数 ---------------- */
const handleAddMoreImage = () => {
  imagesToAdd.value.push({
    image_url: String(imageToAdd.value.image_url),
    image_source: String(imageToAdd.value.image_source),
    object_key: imageToAdd.value.object_key ? String(imageToAdd.value.object_key) : '',
  })

  imageToAdd.value.image_url = ''
  imageToAdd.value.image_source = ''
  imageToAdd.value.object_key = ''
}

const handleAddImageMode = () => {
  if (imageToAdd.value.image_source === '') {
    imageToAdd.value.image_source = ImageSource.LOCAL
  }
  currentMode.value = Mode.Image
}
// const fileInput = ref<HTMLInputElement | null>(null)
// const handleTriggerUpload = () => {
//   imageSourceMemory.value = imageToAdd.value.image_source
//   if (fileInput.value) {
//     fileInput.value.click()
//   }
// }

// 旧版本单图片上传函数
// const handleUploadImage = async (event: Event) => {
//   const target = event.target as HTMLInputElement
//   const file = target.files?.[0]
//   if (!file) return

//   try {
//     const res = await theToast.promise(fetchUploadImage(file), {
//       loading: '图片上传中...',
//       success: '图片上传成功！',
//       error: '图片上传失败，请重试',
//     })

//     if (res.code === 1) {
//       imageToAdd.value.image_url = String(res.data)
//       imageToAdd.value.image_source = ImageSource.LOCAL
//       handleAddMoreImage()

//       if (isUpdateMode.value && echoToUpdate.value) {
//         handleAddorUpdateEcho(true)
//       }
//     }
//   } catch (err) {
//     console.error('上传图片出错:', err)
//   } finally {
//     // 重置 input，防止同图重复上传失效
//     target.value = ''
//   }
// }

const handleSetImageSource = (newSource: string) => {
  imageToAdd.value.image_source = newSource
}

const handleUppyUploaded = (files: App.Api.Ech0.ImageToAdd[]) => {
  files.forEach((file) => {
    imageToAdd.value.image_url = file.image_url
    imageToAdd.value.image_source = file.image_source
    console.log('上传成功的图片:', file.object_key)
    imageToAdd.value.object_key = file.object_key ? file.object_key : ''
    handleAddMoreImage()
  })

  if (isUpdateMode.value && echoToUpdate.value) {
    handleAddorUpdateEcho(true)
  }
}
/* ----------------------------------------------------- */

/* ------------------ 与Echo/Todo相关的各种函数 -------------- */
// 处理Echo的私密性切换
const handlePrivate = () => {
  echoToAdd.value.private = !echoToAdd.value.private
}

// 执行编辑器清空操作
const handleClear = () => {
  echoToAdd.value.content = ''
  echoToAdd.value.image_url = null
  echoToAdd.value.image_source = null
  echoToAdd.value.images = []
  echoToAdd.value.private = false
  echoToAdd.value.extension = null
  echoToAdd.value.extension_type = null
  extensionToAdd.value.extension = ''
  extensionToAdd.value.extension_type = ''
  videoURL.value = ''
  imagesToAdd.value = []
  imageToAdd.value.image_url = ''
  imageToAdd.value.image_source = ImageSource.LOCAL
  imageToAdd.value.object_key = ''
  imageIndex.value = 0
}

// 处理Echo的添加或更新
const handleAddorUpdateEcho = async (justSyncImages: boolean) => {
  if (isSubmitting.value) return // 防止重复提交
  isSubmitting.value = true

  try {
    echoToAdd.value.images = imagesToAdd.value // 将图片数组添加到Echo中

    // 检查是否有外部链接分享
    if (extensionToAdd.value.extension_type === ExtensionType.WEBSITE) {
      // 检查是否存在网站链接
      if (websiteToAdd.value.title.length > 0 && websiteToAdd.value.site.length === 0) {
        theToast.error('网站链接不能为空！')
        return
      }

      // 检查是否存在网站标题
      if (websiteToAdd.value.title.length === 0 && websiteToAdd.value.site.length > 0) {
        websiteToAdd.value.title = '外部链接'
      }

      // 检查网站标题和链接是否都存在
      if (websiteToAdd.value.title.length > 0 && websiteToAdd.value.site.length > 0) {
        // 将网站标题和链接添加到扩展中 (序列化为json)
        extensionToAdd.value.extension = JSON.stringify({
          title: websiteToAdd.value.title,
          site: websiteToAdd.value.site,
        })
      } else {
        extensionToAdd.value.extension = ''
        extensionToAdd.value.extension_type = ''
      }
    }

    // 检查最终的Extension模块是否有内容
    if (
      extensionToAdd.value.extension.length > 0 &&
      extensionToAdd.value.extension_type.length > 0
    ) {
      echoToAdd.value.extension = extensionToAdd.value.extension
      echoToAdd.value.extension_type = extensionToAdd.value.extension_type
    } else {
      echoToAdd.value.extension = null
      echoToAdd.value.extension_type = null
    }

    // 检查Echo是否为空
    if (
      !echoToAdd.value.content &&
      (!echoToAdd.value.images || echoToAdd.value.images.length === 0) &&
      !echoToAdd.value.extension &&
      !echoToAdd.value.extension_type
    ) {
      if (isUpdateMode.value) {
        theToast.error('待更新的Echo不能为空！')
        return
      } else {
        theToast.error('待添加的Echo不能为空！')
        return
      }
    }

    // if (!echoToAdd.value.image_url || echoToAdd.value.image_url.length === 0) {
    //   echoToAdd.value.image_source = null
    // }

    // === 更新模式 ===
    // 检查是否处于更新模式
    if (isUpdateMode.value) {
      // 处于更新模式，执行更新操作
      if (!echoToUpdate.value) {
        theToast.error('没有待更新的Echo！')
        return
      }

      // 回填 echoToUpdate
      echoToUpdate.value.content = echoToAdd.value.content
      echoToUpdate.value.private = echoToAdd.value.private
      echoToUpdate.value.images = echoToAdd.value.images
      echoToUpdate.value.extension = echoToAdd.value.extension
      echoToUpdate.value.extension_type = echoToAdd.value.extension_type

      // 更新Echo
      const res = await fetchUpdateEcho(echoToUpdate.value)
      if (res.code === 1 && !justSyncImages) {
        theToast.success('更新成功！')
        handleClear()
        echoStore.refreshEchos()
        isUpdateMode.value = false
        echoToUpdate.value = null
        currentMode.value = Mode.ECH0
      } else if (res.code === 1 && justSyncImages) {
        theToast.success('发现图片更改，已自动更新同步Echo！')
      } else {
        theToast.error('更新失败，请稍后再试！')
      }
      isSubmitting.value = false
      return
    }

    // === 添加模式 ===
    // 不是Echo更新模式，执行添加操作
    const res = await fetchAddEcho(echoToAdd.value)
    if (res.code === 1) {
      theToast.success('发布成功！')
      handleClear()
      echoStore.refreshEchos()
      currentMode.value = Mode.ECH0
    }
  } finally {
    isSubmitting.value = false
  }
}

// 处理Todo的添加
const handleAddTodo = () => {
  if (todoToAdd.value.content === '') {
    theToast.error('内容不能为空！')
    return
  }

  fetchAddTodo(todoToAdd.value).then((res) => {
    if (res.code === 1) {
      theToast.success('添加成功！')
      todoToAdd.value.content = ''
      getTodos()
    }
  })
}

// 处理不同模式下Echo或Todo的添加和更新操作
const handleAddorUpdate = () => {
  if (todoMode.value) {
    handleAddTodo()
  } else {
    handleAddorUpdateEcho(false)
  }
}

// 退出Echo更新模式
const handleExitUpdateMode = () => {
  isUpdateMode.value = false
  echoToUpdate.value = null
  handleClear()
  theToast.info('已退出更新模式！')
}
/* ----------------------------------------------------------- */

/* ------------------ 与Watch相关的各种函数 -------------- */
// 监听用户输入
watch(
  () => videoURL.value,
  (newVal) => {
    if (newVal.length > 0) {
      const bvRegex = /(BV[0-9A-Za-z]{10})/
      const ytRegex =
        /(?:https?:\/\/(?:www\.)?)?(?:youtu\.be\/|youtube\.com\/(?:(?:watch)?\?(?:.*&)?v(?:i)?=|(?:embed)\/))([\w-]+)/
      let match = newVal.match(bvRegex)
      if (match) {
        extensionToAdd.value.extension = match[0] //bilibili
      } else {
        match = newVal.match(ytRegex)
        if (match) {
          extensionToAdd.value.extension = match[1] //youtube
        } else {
          theToast.error('请输入正确的B站/YT分享链接！')
        }
      }
    }
  },
)

// 监听是否处于更新模式
watch(
  () => isUpdateMode.value,
  (newVal) => {
    if (newVal) {
      // 处于更新模式（将待更新的数据填充到编辑器）
      // 0. 清空编辑器
      handleClear()

      // 1. 填充本文
      echoToAdd.value.content = echoToUpdate.value?.content || ''
      echoToAdd.value.private = echoToUpdate.value?.private || false

      // 2. 填充图片
      if (echoToUpdate.value?.images && echoToUpdate.value.images.length > 0) {
        imagesToAdd.value = echoToUpdate.value.images.map((img) => ({
          image_url: img.image_url || '',
          image_source: img.image_source || '',
        }))
      } else {
        imagesToAdd.value = []
      }

      // 3. 填充扩展
      if (echoToUpdate.value?.extension && echoToUpdate.value.extension_type) {
        currentExtensionType.value = echoToUpdate.value.extension_type as ExtensionType
        extensionToAdd.value.extension = echoToUpdate.value.extension
        extensionToAdd.value.extension_type = echoToUpdate.value.extension_type
        // 根据扩展类型填充
        switch (echoToUpdate.value.extension_type) {
          case ExtensionType.MUSIC:
            break

          case ExtensionType.VIDEO:
            videoURL.value = echoToUpdate.value.extension // 直接使用extension填充B站链接
            break

          case ExtensionType.GITHUBPROJ:
            break

          case ExtensionType.WEBSITE:
            // 反序列化网站链接
            const websiteData = JSON.parse(echoToUpdate.value.extension) as {
              title?: string
              site?: string
            }
            websiteToAdd.value.title = websiteData.title || ''
            websiteToAdd.value.site = websiteData.site || ''
            break
        }
      }

      // 4. 回到页面顶部（触发BackToTop）
      window.scrollTo({ top: 0, behavior: 'smooth' })

      // 5. 弹出通知，提示可以编辑了
      theToast.info('已进入更新模式，请编辑内容后点击更新按钮！')
    } else {
      // 退出更新模式
      echoToUpdate.value = null
    }
  },
)
/* ------------------------------------------------------- */
</script>
