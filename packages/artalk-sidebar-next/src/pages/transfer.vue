<script setup lang="ts">
import { useNavStore } from '../stores/nav'
import { artalk } from '../global'
import type { SiteData } from 'artalk/types/artalk-data'

const nav = useNavStore()

onMounted(() => {
  nav.updateTabs({
    'import': '导入',
    'export': '导出',
  }, 'import')
})

</script>

<template>
  <div class="atk-log-wrap" style="display: none">
    <div class="atk-log-back-btn">返回</div>
    <div class="atk-log"></div>
  </div>
  <div class="atk-form">
    <div class="atk-label atk-data-file-label">Artrans 数据文件</div>
    <div class="atk-file-upload-group">
      <div class="atk-file-input-wrap atk-fade-in">
        <input type="file" name="AtkDataFile" accept=".artrans" />
        <span class="atk-desc">
          使用「<a href="https://artalk.js.org/guide/transfer.html" target="_blank">转换工具</a>」将评论数据转为 Artrans 格式
        </span>
      </div>
      <div class="atk-uploading-wrap atk-fade-in">
        <div class="atk-progress">
          <div class="atk-bar"></div>
        </div>
        <div class="atk-status">
          上传中
          <span class="atk-curt">0%</span>
          ...
          <span class="atk-abort">取消</span>
        </div>
      </div>
    </div>
    <div class="atk-label">目标站点名</div>
    <input
      type="text"
      name="AtkSiteName"
      placeholder="输入内容..."
      autocomplete="off"
    />
    <div class="atk-label">目标站点 URL</div>
    <input
      type="text"
      name="AtkSiteURL"
      placeholder="输入内容..."
      autocomplete="off"
    />
    <div class="atk-label">启动参数（可选）</div>
    <textarea name="AtkPayload"></textarea>
    <span class="atk-desc">
      参考「<a href="https://artalk.js.org/guide/transfer.html" target="_blank">文档 · 数据迁移</a>」
    </span>
    <button class="atk-btn" name="AtkSubmit">导入</button>
  </div>
</template>

<style scoped lang="scss">

</style>
