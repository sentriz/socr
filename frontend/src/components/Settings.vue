<template>
  <h2>importer</h2>
  <SettingsImporter />
  <hr />
  <h2>about</h2>
  <SettingsAbout />
</template>

<script setup="props">
import SettingsImporter from "./SettingsImporter.vue";
import SettingsAbout from "./SettingsAbout.vue";
export default {
  components: { SettingsImporter, SettingsAbout },
  props: {},
};
</script>
