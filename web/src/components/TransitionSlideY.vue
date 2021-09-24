<template>
  <transition
    enter-active-class="transition duration-100"
    enter-to-class="translate-y-0"
    enter-from-class="translate-y-full"
    leave-active-class="transition duration-100"
    leave-to-class="translate-y-full"
    leave-from-class="translate-y-0"
  >
    <slot />
  </transition>
</template>
