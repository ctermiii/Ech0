<template>
  <div v-if="images?.length" class="image-gallery-container">
    <!-- 瀑布流布局 -->
    <div
      v-if="layout === ImageLayout.WATERFALL || !layout"
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
          :src="baseUrl ? getHubImageUrl(src, baseUrl) : getImageUrl(src)"
          :alt="`预览图片${idx + 1}`"
          loading="lazy"
          class="echoimg block max-w-full h-auto"
        />
      </button>
    </div>

    <!-- 九宫格布局 -->
    <div v-if="layout === ImageLayout.GRID" class="imgwidth mx-auto mb-4">
      <div class="grid grid-cols-3 gap-2">
        <button
          v-for="(src, idx) in displayedImages"
          :key="idx"
          class="bg-transparent border-0 p-0 cursor-pointer overflow-hidden aspect-square relative"
          @click="openFancybox(idx)"
        >
          <img
            :src="baseUrl ? getHubImageUrl(src, baseUrl) : getImageUrl(src)"
            :alt="`预览图片${idx + 1}`"
            loading="lazy"
            class="echoimg w-full h-full object-cover"
          />

          <div v-if="extraCount > 0 && idx === 8" class="more-overlay" aria-hidden="true">
            +{{ extraCount }}
          </div>
        </button>
      </div>
    </div>

    <!-- 单图轮播布局 -->
    <div v-if="layout === ImageLayout.CAROUSEL" class="imgwidth mx-auto mb-4">
      <div class="carousel-container">
        <button
          v-if="images[carouselIndex]"
          class="carousel-slide bg-transparent border-0 p-0 cursor-pointer w-full overflow-hidden"
          @click="openFancybox(carouselIndex)"
        >
          <img
            :src="
              baseUrl
                ? getHubImageUrl(images[carouselIndex]!, baseUrl)
                : getImageUrl(images[carouselIndex]!)
            "
            :alt="`预览图片${carouselIndex + 1}`"
            loading="lazy"
            class="echoimg w-full h-auto"
          />
        </button>

        <div v-if="images.length > 1" class="carousel-controls">
          <button
            class="carousel-btn carousel-prev"
            @click="prevCarousel"
            :disabled="carouselIndex === 0"
          >
            ←
          </button>
          <div class="carousel-indicator">{{ carouselIndex + 1 }} / {{ images.length }}</div>
          <button
            class="carousel-btn carousel-next"
            @click="nextCarousel"
            :disabled="carouselIndex === images.length - 1"
          >
            →
          </button>
        </div>

        <div v-if="images.length > 1" class="carousel-dots">
          <button
            v-for="(_, idx) in images"
            :key="idx"
            class="carousel-dot"
            :class="{ active: idx === carouselIndex }"
            @click="carouselIndex = idx"
          />
        </div>
      </div>
    </div>

    <!-- 水平轮播布局 -->
    <div v-if="layout === ImageLayout.HORIZONTAL" class="imgwidth mx-auto mb-4">
      <div class="horizontal-scroll-container">
        <div class="horizontal-scroll-wrapper">
          <button
            v-for="(src, idx) in images"
            :key="idx"
            class="horizontal-item bg-transparent border-0 p-0 cursor-pointer shrink-0"
            @click="openFancybox(idx)"
          >
            <img
              :src="baseUrl ? getHubImageUrl(src, baseUrl) : getImageUrl(src)"
              :alt="`预览图片${idx + 1}`"
              loading="lazy"
              class="echoimg h-full w-auto object-contain"
            />
          </button>
        </div>
      </div>
      <div class="scroll-hint">← 左右滑动查看更多 →</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onBeforeUnmount, ref, computed } from 'vue'
import { getImageUrl, getHubImageUrl } from '@/utils/other'
import { Fancybox } from '@fancyapps/ui'
import '@fancyapps/ui/dist/fancybox/fancybox.css'
import { ImageLayout } from '@/enums/enums'

const props = defineProps<{
  images: App.Api.Ech0.Echo['images']
  baseUrl?: string
  layout?: ImageLayout | string | undefined
}>()

const baseUrl = props.baseUrl

// 布局状态（来自 props.layout）
const layout = props.layout || ImageLayout.WATERFALL

// 轮播索引
const carouselIndex = ref(0)

// 只显示前 9 张（用于九宫格），第 9 张显示 "+N" 覆盖层
const displayedImages = computed(() => (props.images ? props.images.slice(0, 9) : []))
const extraCount = computed(() =>
  props.images ? (props.images.length > 9 ? props.images.length - 9 : 0) : 0,
)

// 瀑布流布局：获取列跨度
const getColSpan = (idx: number, total: number) => {
  if (total === 1) return 'col-span-1 justify-self-center'
  if (idx === 0 && total % 2 !== 0) return 'col-span-2'
  return ''
}

// 轮播导航
const prevCarousel = () => {
  if (carouselIndex.value > 0) carouselIndex.value--
}
const nextCarousel = () => {
  if (carouselIndex.value < (props.images ? props.images.length - 1 : 0)) carouselIndex.value++
}

function openFancybox(startIndex: number) {
  const items = (props.images || []).map((src) => ({
    src: baseUrl ? getHubImageUrl(src, baseUrl) : getImageUrl(src),
    type: 'image' as const,
    thumb: baseUrl ? getHubImageUrl(src, baseUrl) : getImageUrl(src),
  }))

  Fancybox.show(items, {
    theme: 'auto',
    zoomEffect: true,
    fadeEffect: true,
    startIndex,
    backdropClick: 'close',
    dragToClose: true,
    keyboard: {
      Escape: 'close',
      ArrowRight: 'next',
      ArrowLeft: 'prev',
      Delete: 'close',
      Backspace: 'close',
      ArrowDown: 'next',
      ArrowUp: 'prev',
      PageUp: 'close',
      PageDown: 'close',
    },
  })
}

onMounted(() => {
  Fancybox.bind('[data-fancybox]', {})
})

onBeforeUnmount(() => {})
</script>

<style scoped>
.image-gallery-container {
  position: relative;
}

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
  transition:
    transform 0.3s ease,
    box-shadow 0.3s ease;
}

button:hover .echoimg {
  transform: scale(1.03);
  box-shadow:
    0 2px 4px rgba(0, 0, 0, 0.04),
    0 4px 8px rgba(0, 0, 0, 0.04),
    0 8px 16px rgba(0, 0, 0, 0.04),
    0 16px 32px rgba(0, 0, 0, 0.04);
}

/* carousel, horizontal, grid styles (copied/adapted from provided template) */
.carousel-container {
  position: relative;
  width: 100%;
}
.carousel-slide {
  position: relative;
  width: 100%;
  display: block;
}
.carousel-controls {
  position: absolute;
  bottom: 12px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  align-items: center;
  gap: 16px;
  background: rgba(0, 0, 0, 0.4);
  padding: 8px 16px;
  border-radius: 20px;
  z-index: 10;
}
.carousel-btn {
  background-color: rgba(255, 255, 255, 0.8);
  border: none;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  cursor: pointer;
  font-size: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
  color: #333;
}
.carousel-btn:hover:not(:disabled) {
  background-color: white;
  transform: scale(1.1);
}
.carousel-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
.carousel-indicator {
  color: white;
  font-size: 12px;
  min-width: 50px;
  text-align: center;
}
.carousel-dots {
  position: absolute;
  bottom: 60px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 8px;
  z-index: 10;
}
.carousel-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: rgba(255, 255, 255, 0.5);
  border: none;
  cursor: pointer;
  transition: all 0.3s ease;
  padding: 0;
}
.carousel-dot:hover {
  background-color: rgba(255, 255, 255, 0.8);
}
.carousel-dot.active {
  background-color: white;
  transform: scale(1.2);
}

.horizontal-scroll-container {
  position: relative;
  width: 100%;
  overflow-x: auto;
  overflow-y: hidden;
  scroll-behavior: smooth;
  -webkit-overflow-scrolling: touch;
  scrollbar-width: thin;
  scrollbar-color: rgba(0, 0, 0, 0.1) transparent;
}
.horizontal-scroll-container::-webkit-scrollbar {
  height: 4px;
}
.horizontal-scroll-wrapper {
  display: flex;
  gap: 8px;
  padding: 4px 0;
  align-items: center;
}
.horizontal-item {
  flex-shrink: 0;
  height: 200px;
  width: auto;
  min-width: 160px;
  overflow: hidden;
}
.scroll-hint {
  text-align: center;
  font-size: 12px;
  color: #999;
  margin-top: 8px;
  animation: hint-pulse 2s infinite;
}
@keyframes hint-pulse {
  0%,
  100% {
    opacity: 0.5;
  }
  50% {
    opacity: 1;
  }
}

.more-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.45);
  color: #fff;
  font-size: 20px;
  font-weight: 600;
  border-radius: 8px;
}
</style>
