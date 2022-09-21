<route lang="json">
{
  "meta": {
    "title": "Home"
  }
}
</route>

<script setup lang="ts">
const router = useRouter()

onMounted(() => {
  router.replace('/comments')
})
</script>

<template>
  <div>
    Hello index page
  </div>
</template>
