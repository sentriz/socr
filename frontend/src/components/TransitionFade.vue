<template>
  <Transition
    enter-active-class="duration-500 ease-in-out"
    enter-to-class="opacity-100"
    enter-from-class="opacity-0"
    leave-active-class="duration-500 ease-in-out"
    leave-to-class="opacity-0"
    leave-from-class="opacity-100"
  >
    <slot />
  </Transition>
</template>
