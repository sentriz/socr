<template>
  <h2>importer</h2>
  <settings-importer />
  <hr />
  <div class="flex flex-wrap gap-3">
    <div class="flex-auto">
      <h2>about</h2>
      <settings-about />
    </div>
    <div class="flex-auto">
      <h2>screenshots</h2>
      <settings-directories />
    </div>
  </div>
</template>

<script setup lang="ts">
import SettingsImporter from './SettingsImporter.vue'
import SettingsAbout from './SettingsAbout.vue'
import SettingsDirectories from './SettingsDirectories.vue'
</script>
