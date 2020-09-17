<template>
  <div class="container mx-auto p-6 py-3">
    <ul class="flex my-3">
      <li class="mr-6">
        <router-link
          class="text-gray-500 hover:text-blue-300"
          :to="{ name: 'search' }"
        >
          <i class="fas fa-search"></i>
          Search
        </router-link>
      </li>
      <li class="mr-6">
        <router-link
          class="text-gray-500 hover:text-blue-300"
          :to="{ name: 'settings' }"
        >
          <i class="fas fa-cogs"></i>
          Settings
        </router-link>
      </li>
    </ul>
    <hr class="mt-0 mb-6" />
    <router-view />
  </div>
</template>

<script>
import Search from "./Search.vue";

export default {
  name: "App",
};
</script>

<style scoped>
.router-link-active {
  @apply text-blue-500;
}
</style>
