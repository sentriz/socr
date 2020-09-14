export const urlImage = "/api/image";
export const urlSearch = "/api/search";
export const urlSocket = "/api/ws";

export const doSearch = async (body) => {
  const response = await fetch(urlSearch, {
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
};
