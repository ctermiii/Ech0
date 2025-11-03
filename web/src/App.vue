<script setup lang="ts">
import { RouterView } from 'vue-router'
import { onMounted, ref } from 'vue'
import { watch } from 'vue'
import { useSettingStore } from '@/stores/setting'
import { storeToRefs } from 'pinia'
import { Toaster } from 'vue-sonner'
import 'vue-sonner/style.css'
import BaseDialog from './components/common/BaseDialog.vue'

import { useBaseDialog } from '@/composables/useBaseDialog'

const { register, title, description, handleConfirm } = useBaseDialog()
const dialogRef = ref()

const settingStore = useSettingStore()
const { SystemSetting } = storeToRefs(settingStore)

watch(
  () => SystemSetting.value.site_title,
  (title) => {
    if (title) document.title = title
  },
  { immediate: true },
)

const injectCustomContent = () => {
  // 注入自定义 CSS
  if (SystemSetting.value.custom_css && SystemSetting.value.custom_css.length > 0) {
    const styleTag = document.createElement('style')
    styleTag.textContent = SystemSetting.value.custom_css
    document.head.appendChild(styleTag)
  }

  // 注入自定义 JS
  if (SystemSetting.value.custom_js && SystemSetting.value.custom_js.length > 0) {
    const scriptTag = document.createElement('script')
    scriptTag.textContent = SystemSetting.value.custom_js
    document.body.appendChild(scriptTag)
  }
}

onMounted(() => {
  // 注入自定义CSS 和 JS
  watch(
    () => SystemSetting.value.custom_css || SystemSetting.value.custom_js,
    (newSetting) => {
      if (newSetting) {
        injectCustomContent()
      }
    },
    { immediate: true },
  )

  // 初始注入
  register(dialogRef.value) // 全局注册弹窗对话框
})
</script>

<template>
  <RouterView />
  <Toaster theme="light" position="top-right" :expand="false" richColors />
  <BaseDialog ref="dialogRef" :title="title" :description="description" @confirm="handleConfirm" />
</template>

<style scoped></style>
