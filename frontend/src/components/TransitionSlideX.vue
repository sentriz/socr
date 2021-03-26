<template>
  <Transition
    enter-active-class="transform transition ease-in-out duration-200"
    enter-to-class="translate-x-0"
    enter-from-class="translate-x-full"
    leave-active-class="transform transition ease-in-out duration-200"
    leave-to-class="translate-x-full"
    leave-from-class="translate-x-0"
  >
    <slot />
  </Transition>
</template>
