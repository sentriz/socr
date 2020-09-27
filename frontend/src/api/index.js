export const urlImage = "/api/image";
export const urlSearch = "/api/search";
export const urlSocket = "/api/ws";
export const urlStartImport = "/api/start_import";
export const urlAuthenticate = "/api/authenticate";

const req = async (url, options) => {
  const response = await fetch(url, {
    ...options,
    headers: { "authorization": `bearer ${tokenGet()}` }
  })
  return await response.json();
}

export const reqSearch = async (body) => req(urlSearch, {
  method: "POST",
  body: JSON.stringify(body),
})

export const reqAuthenticate = async (body) => req(urlAuthenticate, {
  method: "PUT",
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

const tokenKey = "token"
export const tokenSet = (token) => localStorage.setItem(tokenKey, token)
export const tokenGet = () => localStorage.getItem(tokenKey)
export const tokenHas = () => !!localStorage.getItem(tokenKey)
