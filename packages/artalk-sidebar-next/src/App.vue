<script setup lang="ts">

</script>

<template>
  <Header />
  <Tab />

  <h2>{{ $route.name }}</h2>
  <RouterLink to="/">Home</RouterLink>
  <RouterLink to="/about">About</RouterLink>
  <main>
    <router-view />
  </main>
</template>

<style scoped lang="scss">

</style>
