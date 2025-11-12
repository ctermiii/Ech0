<template>
  <div>
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
          @click="editorStore.handleAddMoreImage"
          title="添加更多图片"
        />
      </div>
    </div>

    <!-- 布局方式选择 -->
    <div class="mb-3 flex items-center gap-2">
      <span class="text-gray-500">布局方式：</span>
      <BaseSelect
        v-model="echoToAdd.layout"
        :options="layoutOptions"
        class="w-32 h-7"
        placeholder="请选择布局方式"
      />
    </div>

    <!-- 当前上传方式与状态 -->
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

      <!-- 图片上传 -->
      <TheUppy
        v-if="imageToAdd.image_source !== ImageSource.URL"
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
</template>

<script setup lang="ts">
import { useEditorStore } from '@/stores/editor'
import { useSettingStore } from '@/stores/setting'
import { storeToRefs } from 'pinia'
import { ref } from 'vue'
import { ImageSource } from '@/enums/enums'
import Url from '@/components/icons/url.vue'
import Upload from '@/components/icons/upload.vue'
import Bucket from '@/components/icons/bucket.vue'
import Addmore from '@/components/icons/addmore.vue'
import BaseButton from '@/components/common/BaseButton.vue'
import BaseSelect from '@/components/common/BaseSelect.vue'
import BaseInput from '@/components/common/BaseInput.vue'
import TheUppy from '@/components/advanced/TheUppy.vue'

const editorStore = useEditorStore()
const { imageToAdd, ImageUploading, echoToAdd } = storeToRefs(editorStore)
const settingStore = useSettingStore()
const { S3Setting } = storeToRefs(settingStore)

// 布局选择（绑定到 editorStore.echoToAdd.layout）
const layoutOptions = [
  { label: '瀑布流', value: 'waterfall' },
  { label: '九宫格', value: 'grid' },
  { label: '单图轮播', value: 'carousel' },
  { label: '水平轮播', value: 'horizontal' },
]
</script>

<style scoped></style>
