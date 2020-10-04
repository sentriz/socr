<template>
  <div class="container mx-auto bg-white h-screen p-5 m-5">
    <ul class="flex my-3">
      <div class="flex-1 space-x-6">
        <router-link class="text-gray-500 hover:text-blue-300" :to="{ name: 'search' }">
          <i class="fas fa-search"></i>
          search
        </router-link>
        <router-link class="text-gray-500 hover:text-blue-300" :to="{ name: 'settings' }">
          <i class="fas fa-cogs"></i>
          settings
        </router-link>
      </div>
      <router-link class="text-gray-500 hover:text-blue-300" :to="{ name: 'logout' }">
        <i class="fas fa-sign-out-alt"></i>
        logout
      </router-link>
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
