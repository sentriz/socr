<template>
  <div class="container mx-auto px-6 py-3">
    <ul class="flex my-3">
      <li class="mr-6">
        <router-link
          class="text-blue-500 hover:text-blue-800"
          :to="{ name: 'search' }"
        >
          Search
        </router-link>
      </li>
      <li class="mr-6">
        <router-link
          class="text-blue-500 hover:text-blue-800"
          :to="{ name: 'settings' }"
        >
          Settings
        </router-link>
      </li>
    </ul>
    <hr class="mt-3 mb-6" />
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
  @apply font-bold;
}
</style>
