export const urlImage = "/api/image";
export const urlSearch = "/api/search";
export const urlStartImport = "/api/start_import";
export const urlAuthenticate = "/api/authenticate";
export const urlSocket = "/api/websocket";
export const urlSocketPub = "/api/websocket_pub";

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

export const reqImage = async (id) => req(`${urlImage}/${id}`, {
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

export const newSocket = () => new WebSocket(
  `${socketProtocol}//${socketHost}${urlSocket}?token=${tokenGet()}`)
export const newSocketPub = () => new WebSocket(
  `${socketProtocol}//${socketHost}${urlSocketPub}`)
