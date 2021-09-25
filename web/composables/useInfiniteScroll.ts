import { ref } from 'vue'
import { useIntersectionObserver } from '@vueuse/core'

export default (onNewPage: () => Promise<void>) => {
  const target = ref(null)

  useIntersectionObserver(target, async ([{ isIntersecting }], _) => {
    if (!isIntersecting) return
    await onNewPage()
  })

  return target
}
