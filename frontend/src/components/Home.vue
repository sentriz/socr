<template>
  <div class="container mx-auto h-screen p-5 mx-5">
    <ul class="flex my-3">
      <div class="flex-1 space-x-6">
        <router-link class="nav" :to="{ name: 'search' }">
          <i class="nav-icon fas fa-search"></i>
          <span>search</span>
        </router-link>
        <router-link class="nav" :to="{ name: 'settings' }">
          <i class="nav-icon fas fa-cogs"></i>
          <span>settings</span>
        </router-link>
      </div>
      <router-link class="nav" :to="{ name: 'logout' }">
        <i class="nav-icon fas fa-sign-out-alt"></i>
        <span>logout</span>
      </router-link>
    </ul>
    <hr class="mt-0" />
    <router-view />
  </div>
</template>

<script setup="props">
import Search from "./Search.vue";
export default {
  components: { Search },
  props: {},
};
</script>

<style scoped>
.router-link-active {
  @apply text-blue-500;
}
</style>
