<template>
  <router-link :to="`${$attrs.to}`" class="flex h-5 items-center space-x-2 text-gray-500 hover:text-blue-400" active-class="text-blue-500">
    <slot name="icon"></slot>
    <span>
      <slot name="name"></slot>
    </span>
  </router-link>
</template>
