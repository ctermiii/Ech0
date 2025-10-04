<template>
  <div class="p-1 my-4">
    <!-- 扩展附加内容 -->
    <div class="mb-1">
      <h2 class="text-gray-500 font-bold mb-1">扩展附加内容</h2>
      <div class="flex flex-row items-center gap-2">
        <!-- 添加音乐 -->
        <BaseButton
          :icon="Music"
          class="w-7 h-7 rounded-md"
          title="添加音乐"
          @click="handleAddExtension(ExtensionType.MUSIC)"
        />
        <!-- 添加视频 -->
        <BaseButton
          :icon="Video"
          class="w-7 h-7 rounded-md"
          title="添加视频"
          @click="handleAddExtension(ExtensionType.VIDEO)"
        />
        <!-- 添加Github项目 -->
        <BaseButton
          :icon="Githubproj"
          class="w-7 h-7 rounded-md"
          title="添加Github项目"
          @click="handleAddExtension(ExtensionType.GITHUBPROJ)"
        />
        <!-- 添加网站链接 -->
        <BaseButton
          :icon="Weblink"
          class="w-7 h-7 rounded-md"
          title="添加网站链接"
          @click="handleAddExtension(ExtensionType.WEBSITE)"
        />
      </div>
    </div>

    <!-- 模式切换 -->
    <div class="mb-1">
      <h2 class="text-gray-500 font-bold mb-1">模式切换</h2>
      <div class="flex flex-row items-center gap-2">
        <!-- 打开Todo模式 -->
        <BaseButton :icon="Todo" @click="handleTodo" class="w-7 h-7 rounded-md" />
        <BaseButton
          :icon="Audio"
          class="w-7 h-7 rounded-md"
          title="音乐播放"
          @click="handlePlayMusic"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import Weblink from '@/components/icons/weblink.vue'
import Music from '@/components/icons/music.vue'
import Todo from '@/components/icons/todo.vue'
import Video from '@/components/icons/video.vue'
import Githubproj from '@/components/icons/githubproj.vue'
import BaseButton from '@/components/common/BaseButton.vue'
import Audio from '@/components/icons/audio.vue'

import { Mode, ExtensionType } from '@/enums/enums'
import { useTodoStore } from '@/stores/todo'
import { theToast } from '@/utils/toast'

const mode = defineModel<Mode>('mode', {
  required: true,
})
const theExtensionType = defineModel<ExtensionType | undefined>('extensionType', {
  required: true,
})
const theExtensionToAdd = defineModel<{
  extension_type: string
}>('extensionToAdd', {
  required: true,
})

const todoStore = useTodoStore()
const { setTodoMode } = todoStore

const handleAddExtension = (extensiontype: ExtensionType) => {
  mode.value = Mode.EXTEN
  theExtensionType.value = extensiontype
  theExtensionToAdd.value.extension_type = extensiontype
}

const handleTodo = () => {
  setTodoMode(true)
  mode.value = Mode.TODO
}

const handlePlayMusic = () => {
  theToast.info('音乐播放器模式维护中，敬请期待！')
  // mode.value = Mode.PlayMusic
}
</script>

<style scoped></style>
