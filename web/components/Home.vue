<template>
  <div class="container mx-auto min-h-screen min-w-min p-6">
    <ul class="flex">
      <div class="flex flex-1 space-x-6">
        <nav-item :to="{ name: routes.SEARCH }">
          <template v-slot:icon><search-icon class="h-5" /></template>
          <template v-slot:name>search</template>
        </nav-item>
        <nav-item :to="{ name: routes.IMPORTER }">
          <template v-slot:icon><document-search-icon class="h-5" /></template>
          <template v-slot:name>importer</template>
        </nav-item>
        <nav-item :to="{ name: routes.SETTINGS }">
          <template v-slot:icon><cog-icon class="h-5" /></template>
          <template v-slot:name>settings</template>
        </nav-item>
      </div>
      <div class="hidden sm:block">
        <nav-item :to="{ name: routes.LOGOUT }">
          <template v-slot:icon><logout-icon class="h-5" /></template>
          <template v-slot:name>logout</template>
        </nav-item>
      </div>
    </ul>
    <hr class="mt-3" />
    <router-view />
    <teleport to="#overlays">
      <toast-overlay />
    </teleport>
  </div>
</template>

<script setup lang="ts">
import NavItem from './NavItem.vue'
import ToastOverlay from './ToastOverlay.vue'
import { SearchIcon, DocumentSearchIcon, CogIcon, LogoutIcon } from '@heroicons/vue/outline'
import { routes } from '~/router'
</script>
