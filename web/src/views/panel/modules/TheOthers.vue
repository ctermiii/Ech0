<template>
  <div class="w-full px-2">
    <div class="rounded-md shadow-sm ring-1 ring-gray-200 ring-inset bg-white p-4 mb-2">
      <!-- 设置 -->
      <div>
        <div class="flex items-center justify-start mb-3">
          <h1 class="text-gray-600 font-bold text-lg">数据管理</h1>
        </div>

        <div class="flex flex-col gap-4">
          <!-- 备份数据 -->
          <div class="flex flex-start items-center gap-2">
            <p class="text-gray-400">创建快照:</p>
            <BaseButton
              :icon="CreateBackup"
              @click="handleBackup"
              class="rounded-lg !bg-gray-100 !text-gray-600 hover:!bg-gray-200"
              title="创建快照"
            />
          </div>
          <!-- 导出快照 -->
          <div class="flex flex-start items-center gap-2">
            <p class="text-gray-400">导出快照:</p>
            <BaseButton
              :icon="ExportBackup"
              @click="handleBackupExport"
              class="rounded-lg !bg-gray-100 !text-gray-600 hover:!bg-gray-200"
              title="导出快照"
            />
          </div>
          <!-- 恢复数据 -->
          <div class="flex flex-start items-center gap-2">
            <p class="text-gray-400">恢复快照:</p>
            <BaseButton
              :icon="RestoreBackup"
              @click="handleBackupRestore"
              class="rounded-lg !bg-gray-100 !text-gray-600 hover:!bg-gray-200"
              title="恢复快照"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import BaseButton from '@/components/common/BaseButton.vue'
import CreateBackup from '@/components/icons/createbackup.vue'
import ExportBackup from '@/components/icons/exportbackup.vue'
import RestoreBackup from '@/components/icons/restorebackup.vue'
import { fetchBackup, fetchExportBackup } from '@/service/api'
import { theToast } from '@/utils/toast'

const handleBackup = async () => {
  fetchBackup().then((res) => {
    if (res.code === 1) {
      theToast.success('备份成功')
    }
  })
}

const handleBackupExport = async () => {
  try {
    theToast.info('导出中...请稍等', {
      duration: 4000,
    })
    // 1. 获取文件数据（Blob 对象）
    const blob = await fetchExportBackup()

    // 2. 创建一个临时的对象 URL
    const url = window.URL.createObjectURL(blob)

    // 3. 创建一个隐藏的 <a> 元素
    const link = document.createElement('a')
    link.href = url
    link.download = 'backup-latest.zip' // 指定下载文件名

    // 4. 将元素添加到 DOM 并触发点击
    document.body.appendChild(link)
    link.click() // 这会触发浏览器的下载行为

    // 5. 清理：移除元素和释放内存
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
  } catch (error) {
    theToast.error('导出失败')
    console.error('导出备份失败:', error)
  }
}

const handleBackupRestore = async () => {
  theToast.info('功能开发中，请使用TUI/CLI模式执行恢复', {
    duration: 3000,
  })
}
</script>
