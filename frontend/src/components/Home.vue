<template>
  <div class="container mx-auto min-h-screen p-6">
    <ul class="flex">
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
      <div class="hidden sm:block">
        <router-link class="nav" :to="{ name: 'logout' }">
          <i class="nav-icon fas fa-sign-out-alt"></i>
          <span>logout</span>
        </router-link>
      </div>
    </ul>
    <hr class="mt-3" />
    <router-view />
  </div>
</template>

<style scoped>
.nav       { @apply text-gray-500; }
.nav:hover { @apply text-blue-400; }
.nav-icon  { @apply mr-2; }

.router-link-active {
  @apply text-blue-500;
}
</style>
