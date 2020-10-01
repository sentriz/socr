<template>
  <div class="container mx-auto bg-white h-screen p-5">
    <ul class="flex my-3">
      <li class="mr-6">
        <router-link class="text-gray-500 hover:text-blue-300" :to="{ name: 'search' }">
          <i class="fas fa-search"></i>
          Search
        </router-link>
      </li>
      <li class="mr-6">
        <router-link class="text-gray-500 hover:text-blue-300" :to="{ name: 'settings' }">
          <i class="fas fa-cogs"></i>
          Settings
        </router-link>
      </li>
    </ul>
    <hr class="mt-0 mb-6" />
    <router-view />
  </div>
</template>

<script setup>
export { default as Search } from "./Search.vue";
</script>

<style scoped>
.router-link-active {
  @apply text-blue-500;
}
</style>
