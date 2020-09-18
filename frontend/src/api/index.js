export const urlImage = "/api/image";
export const urlSearch = "/api/search";
export const urlSocket = "/api/ws";
export const urlStartImport = "/api/start_import";

const req = async (url, options) => {
  const response = await fetch(url, options)
  return await response.json();
}

export const reqSearch = async (body) => req(urlSearch, {
  method: "POST",
  body: JSON.stringify(body),
})

export const reqStartImport = async () => req(urlStartImport, {
  method: "POST",
})

export const fields = {
  BLOCKS_TEXT: "blocks.text",
  BLOCKS_POSITION: "blocks.position",
  SIZE_HEIGHT: "dimensions.height",
  SIZE_WIDTH: "dimensions.width",
};
