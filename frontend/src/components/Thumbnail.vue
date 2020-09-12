<template>
  <canvas
    ref="canvas"
    :height="screenshot.fields['dimensions.height']"
    :width="screenshot.fields['dimensions.width']"
    :style="{
      background: `url(${imageURL(screenshot.id)})`,
      backgroundSize: 'cover',
    }"
  />
</template>

<script>
const zipBlocks = (screenshot) => {
  const textField = "blocks.text";
  const positionField = "blocks.position";

  let flatText = screenshot.fields[textField];
  let flatPosition = screenshot.fields[positionField];
  if (!Array.isArray(flatText)) flatText = [flatText];
  if (!Array.isArray(flatPosition)) flatPosition = [flatPosition];
  // const locationText =
  //   screenshot.locations[Object.keys(screenshot.locations)[0]];
  // console.log(locationText);
  // const locationIndex = locationText[0].array_positions[0];

  return flatText.map((block, i) => {
    const [minX, minY, maxX, maxY] = flatPosition.slice(4 * i, 4 * i + 4);
    return {
      text: block,
      position: { minX, minY, maxX, maxY },
      // match: i === locationIndex,
    };
  });
};

import { ref, onMounted } from "vue";
import { imageURL } from "../api/";

export default {
  name: "Thumbnail",
  props: {
    screenshot: Object,
  },
  setup(props) {
    const canvas = ref(null);
    const blocks = zipBlocks(props.screenshot);

    onMounted(() => {
      var ctx = canvas.value.getContext("2d");
      for (const block of blocks) {
        ctx.fillStyle = "rgba(236, 201, 75, 0.25)";
        ctx.fillRect(
          block.position.minX,
          block.position.minY,
          block.position.maxX - block.position.minX,
          block.position.maxY - block.position.minY
        );
      }
    });
    return { canvas };
  },
  methods: {
    imageURL,
  },
};
</script>
