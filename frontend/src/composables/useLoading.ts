import { ref } from 'vue'

export default <T extends Array<any>, U>(fn: (...args: T) => Promise<U>) => {
  const loading = ref(false)
  const load = async (...args: T): Promise<U> => {
    loading.value = true
    const resp = await fn(...args)
    loading.value = false
    return resp
  }

  return { loading, load }
}
