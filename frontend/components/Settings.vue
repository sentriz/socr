<template>
  <div class="flex flex-wrap gap-3">
    <div class="flex-auto space-y-2">
      <h2>about</h2>
      <settings-about />
    </div>
    <div class="flex-auto space-y-2">
      <h2>medias</h2>
      <settings-directories />
    </div>
  </div>
</template>

<script setup lang="ts">
import SettingsAbout from './SettingsAbout.vue'
import SettingsDirectories from './SettingsDirectories.vue'
</script>
