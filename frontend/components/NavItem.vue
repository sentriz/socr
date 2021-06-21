<template>
  <router-link class="hover:text-blue-400 flex items-center h-5 space-x-2 text-gray-500" active-class="text-blue-500">
    <slot name="icon"></slot>
    <span><slot name="name"></slot></span>
  </router-link>
</template>
