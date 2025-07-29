<script setup lang="ts">
import { RouterView } from 'vue-router'
import { onMounted } from 'vue'
import { watch } from 'vue'
import { useSettingStore } from '@/stores/settting'
import { storeToRefs } from 'pinia'
import { Toaster } from 'vue-sonner'
import 'vue-sonner/style.css'

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
  watch(
    () => SystemSetting.value.custom_css || SystemSetting.value.custom_js,
    (newSetting) => {
      if (newSetting) {
        injectCustomContent()
      }
    },
    { immediate: true },
  )
})
</script>

<template>
  <RouterView />
  <Toaster theme="light" position="top-right" :expand="false" richColors />
</template>

<style scoped></style>
