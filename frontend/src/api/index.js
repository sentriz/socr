export const urlScreenshot = "/api/screenshot";
export const urlSearch = "/api/search";
export const urlStartImport = "/api/start_import";
export const urlAuthenticate = "/api/authenticate";
export const urlSocket = "/api/websocket";
export const urlAbout = "/api/about";

const req = async (url, options) => {
  const token = tokenGet()
  const response = await fetch(url, {
    ...options,
    headers: token ? {
      "authorization": `bearer ${token}`
    } : {}
  })
  return await response.json();
}


export const reqSearchParams = (match) => ({
  size: 40,
  fields: [
    fields.BLOCKS_TEXT,
    fields.BLOCKS_POSITION,
    fields.SIZE_HEIGHT,
    fields.SIZE_WIDTH,
  ],
  highlight: {
    fields: [fields.BLOCKS_TEXT],
  },
  query: {
    match,
    fuzziness: 1,
    field: fields.BLOCKS_TEXT,
    prefix_length: 0,
  },
});

export const reqSearch = async (body) => req(urlSearch, {
  method: "POST",
  body: JSON.stringify(body)
})

export const reqAuthenticate = async (body) => req(urlAuthenticate, {
  method: "PUT",
  body: JSON.stringify(body)
})

export const reqStartImport = async () => req(urlStartImport, {
  method: "POST"
})

export const reqScreenshot = async (id) => req(`${urlScreenshot}/${id}`, {
  method: "GET"
})

export const reqAbout = async () => req(urlAbout, {
  method: "GET"
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

const socketGuesses = {
  "https:": "wss:",
  "http:": "ws:",
}

const socketProtocol = socketGuesses[window.location.protocol]
const socketHost = window.location.host

export const newSocketAuth = (params) => newSocket({ ...params, token: tokenGet() })
export const newSocket = (params) => {
  const paramsEnc = new URLSearchParams(params)
  return new WebSocket(`${socketProtocol}//${socketHost}${urlSocket}?${paramsEnc}`)
}
