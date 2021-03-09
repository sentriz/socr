<template>
  <h2>importer</h2>
  <SettingsImporter />
  <hr />
  <div class="flex gap-3">
    <div>
      <h2>about</h2>
      <SettingsAbout />
    </div>
    <div>
      <h2>directories</h2>
      <SettingsDirectories />
    </div>
  </div>
</template>

<script setup lang="ts">
import SettingsImporter from './SettingsImporter.vue'
import SettingsAbout from './SettingsAbout.vue'
import SettingsDirectories from './SettingsDirectories.vue'
</script>
