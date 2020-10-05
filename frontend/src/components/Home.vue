<template>
  <div class="container mx-auto bg-white h-screen p-5 mx-5">
    <ul class="flex my-3">
      <div class="flex-1 space-x-6">
        <router-link class="text-gray-500 hover:text-blue-300" :to="{ name: 'search' }">
          <i class="fas fa-search mr-2"></i>
          <span>search</span>
        </router-link>
        <router-link class="text-gray-500 hover:text-blue-300" :to="{ name: 'settings' }">
          <i class="fas fa-cogs mr-2"></i>
          <span>settings</span>
        </router-link>
      </div>
      <router-link class="text-gray-500 hover:text-blue-300" :to="{ name: 'logout' }">
        <i class="fas fa-sign-out-alt mr-2"></i>
        <span>logout</span>
      </router-link>
    </ul>
    <hr class="mt-0" />
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
