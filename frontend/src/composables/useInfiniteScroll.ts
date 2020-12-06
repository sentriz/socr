import { ref, onMounted, onUnmounted } from "vue";

export default (onNewPage: () => Promise<void>) => {
  const scroller = ref<HTMLElement>();

  const handleScroll = async () => {
    if (!scroller.value) return;

    const scrollerBottom = scroller.value.getBoundingClientRect().bottom;
    const windowBottom = window.innerHeight;
    if (Math.floor(scrollerBottom) > windowBottom) return;

    await onNewPage();
  };

  onMounted(() => window.addEventListener("scroll", handleScroll));
  onUnmounted(() => window.removeEventListener("scroll", handleScroll));

  return scroller;
};