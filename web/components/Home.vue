<template>
  <div class="min-w-min container min-h-screen p-6 mx-auto">
    <ul class="flex">
      <div class="flex flex-1 space-x-6">
        <nav-item :to="{ name: 'search' }">
          <template v-slot:icon><search-icon class="h-5" /></template>
          <template v-slot:name>search</template>
        </nav-item>
        <nav-item :to="{ name: 'importer' }">
          <template v-slot:icon><document-search-icon class="h-5" /></template>
          <template v-slot:name>importer</template>
        </nav-item>
        <nav-item :to="{ name: 'settings' }">
          <template v-slot:icon><cog-icon class="h-5" /></template>
          <template v-slot:name>settings</template>
        </nav-item>
      </div>
      <div class="sm:block hidden">
        <nav-item :to="{ name: 'logout' }">
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
</script>
