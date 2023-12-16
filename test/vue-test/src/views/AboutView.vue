<script lang="ts" setup>
import Artalk from '@/components/Artalk.vue'
</script>

<template>
  <div class="about">
    <h1>This is an about page</h1>
    <Artalk />
  </div>
</template>

<style>
@media (min-width: 1024px) {
  .about {
    min-height: 100vh;
    display: flex;
    align-items: center;
  }
}
</style>
