<template>
  <Transition
    enter-active-class="ease-in-out duration-500"
    enter-to-class="opacity-100"
    enter-from-class="opacity-0"
    leave-active-class="ease-in-out duration-500"
    leave-to-class="opacity-0"
    leave-from-class="opacity-100"
  >
    <slot />
  </Transition>
</template>
