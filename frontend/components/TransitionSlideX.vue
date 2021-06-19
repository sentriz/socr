<template>
  <transition
    enter-active-class="transition duration-200 ease-in-out"
    enter-to-class="translate-x-0"
    enter-from-class="translate-x-full"
    leave-active-class="transition duration-200 ease-in-out"
    leave-to-class="translate-x-full"
    leave-from-class="translate-x-0"
  >
    <slot />
  </transition>
</template>
