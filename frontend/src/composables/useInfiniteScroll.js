import { ref, onMounted, onUnmounted } from "vue";

export default (onNewPage) => {
  const isLoading = ref(false);
  const scroller = ref(null);

  const handleScroll = async () => {
    if (isLoading.value) return;
    if (!scroller.value) return;

    const scrollerBottom = scroller.value.getBoundingClientRect().bottom;
    const windowBottom = window.innerHeight;
    if (Math.floor(scrollerBottom) > windowBottom) return;

    isLoading.value = true;
    await onNewPage();
    isLoading.value = false;
  };

  onMounted(() => window.addEventListener("scroll", handleScroll));
  onUnmounted(() => window.removeEventListener("scroll", handleScroll));

  return {
    scroller,
    isLoading,
  };
};
