import { ref } from "vue";

export default () => {
  const loading = ref(false)
  const load = async (load, ...args) => {
    loading.value = true
    const resp = await load(...args)
    loading.value = false
    return resp
  }

  return { load, loading }
}
