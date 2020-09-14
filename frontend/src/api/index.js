export const imageURL = (id) => {
  return `/api/image/${id}`
}

export const doSearch = async (body) => {
  const response = await fetch(`/api/search`, {
    method: "POST",
    body: JSON.stringify(body),
  });
  return await response.json();
};

export const fields = {
  BLOCKS_TEXT: "blocks.text",
  BLOCKS_POSITION: "blocks.position",
  DIMENSIONS_HEIGHT: "dimensions.height",
  DIMENSIONS_WIDTH: "dimensions.width",
}
