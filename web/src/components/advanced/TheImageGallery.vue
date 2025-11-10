<template>
  <!-- 瀑布流缩略图 -->
  <div
    v-if="images?.length"
    :class="[
      'imgwidth mx-auto grid gap-2 mb-4',
      images.length === 1 ? 'grid-cols-1 justify-items-center' : 'grid-cols-2',
    ]"
  >
    <button
      v-for="(src, idx) in images"
      :key="idx"
      class="bg-transparent border-0 p-0 cursor-pointer w-fit"
      :class="getColSpan(idx, images.length)"
      @click="openFancybox(idx)"
    >
      <img
        :src="props.baseUrl ? getHubImageUrl(src, props.baseUrl) : getImageUrl(src)"
        alt="`预览图片${idx + 1}`"
        loading="lazy"
        class="echoimg block max-w-full h-auto"
      />
    </button>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onBeforeUnmount } from 'vue'
import { getImageUrl, getHubImageUrl } from '@/utils/other'
import { Fancybox } from '@fancyapps/ui'
import '@fancyapps/ui/dist/fancybox/fancybox.css'

const props = defineProps<{
  images: App.Api.Ech0.Echo['images']
  baseUrl?: string
}>()

// const active = ref<number | null>(null)
const getColSpan = (idx: number, total: number) => {
  // 单张图片占满
  if (total === 1) return 'col-span-1 justify-self-center'
  // 第一张在奇数张时跨两列
  if (idx === 0 && total % 2 !== 0) return 'col-span-2'
  // 其他图片默认占一列
  return ''
}

/* ---------- Fancybox 相关 ---------- */
function openFancybox(startIndex: number) {
  // 构造 Fancybox 需要的 items 数组
  const items = props.images.map((src) => ({
    src: props.baseUrl ? getHubImageUrl(src, props.baseUrl!) : getImageUrl(src),
    type: 'image' as const,
    thumb: props.baseUrl ? getHubImageUrl(src, props.baseUrl!) : getImageUrl(src),
  }))

  Fancybox.show(items, {
    theme: 'auto',
    zoomEffect: true,
    fadeEffect: true,
    startIndex, // 从点击的那张开始
    backdropClick: 'close', // 点击背景不关闭
    dragToClose: true, // 启用拖拽关闭
    keyboard: {
      Escape: 'close', // 按下 Esc 键关闭
      ArrowRight: 'next', // 右箭头切换到下一张
      ArrowLeft: 'prev', // 左箭头切换到上一张
      Delete: 'close', // 删除键关闭
      Backspace: 'close', // 退格键关闭
      ArrowDown: 'next', // 下箭头切换到下一张
      ArrowUp: 'prev',
      PageUp: 'close',
      PageDown: 'close',
    },
  })
}

onMounted(() => {
  // window.addEventListener('keydown', onKeyDown)
  Fancybox.bind('[data-fancybox]', {})
})

onBeforeUnmount(() => {
  // window.removeEventListener('keydown', onKeyDown)
})
</script>

<style scoped>
.imgwidth {
  width: 88%;
}

.echoimg {
  border-radius: 8px;
  box-shadow:
    0 1px 2px rgba(0, 0, 0, 0.02),
    0 2px 4px rgba(0, 0, 0, 0.02),
    0 4px 8px rgba(0, 0, 0, 0.02),
    0 8px 16px rgba(0, 0, 0, 0.02);
}
</style>
