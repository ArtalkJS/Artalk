<script setup lang="ts">

</script>

<template>
  <h2>{{ $route.meta.title }}</h2>
  <RouterLink to="/">Home</RouterLink>
  <RouterLink to="/about">About</RouterLink>
  <main>
    <router-view />
  </main>
</template>

<style scoped>
</style>
