<route lang="json">
{
  "meta": {
    "title": "Artalk Sidebar"
  }
}
</route>

<script setup lang="ts">
import { bootParams } from '../global'

const router = useRouter()

onMounted(() => {
  if (bootParams.view) {
    const splitted = bootParams.view.split('|')
    if (splitted[0]) bootParams.view = splitted[0]
    if (splitted[1]) bootParams.viewParams = JSON.parse(splitted[1])
  }

  const LinkMap: {[key:string]:string} = {
    comments: '/comments',
    pages: '/pages',
    sites: '/sites',
    settings: '/settings'
  }
  router.replace(LinkMap[bootParams.view] || '/comments')
})
</script>

<template>
  <div></div>
</template>
