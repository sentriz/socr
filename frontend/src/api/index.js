export const urlScreenshot = "/api/screenshot";
export const urlSearch = "/api/search";
export const urlStartImport = "/api/start_import";
export const urlAuthenticate = "/api/authenticate";
export const urlSocket = "/api/websocket";
export const urlAbout = "/api/about";
export const urlImportStatus = "/api/import_status";

const req = async (method, url, body, options = {}) => {
  const token = tokenGet();
  const response = await fetch(url, {
    method,
    body: JSON.stringify(body),
    ...options,
    headers: token
      ? {
          authorization: `bearer ${token}`,
        }
      : {},
  });
  return await response.json();
};

export const reqSearch = (body) => req("POST", urlSearch, body);
export const reqAuthenticate = (body) => req("PUT", urlAuthenticate, body);
export const reqStartImport = () => req("POST", urlStartImport);
export const reqScreenshot = (id) => req("GET", `${urlScreenshot}/${id}`);
export const reqAbout = () => req("GET", urlAbout);
export const reqImportStatus = () => req("GET", urlImportStatus);

export const fields = {
  TIMESTAMP: "timestamp",
  TAGS: "tags",
  BLOCKS_TEXT: "blocks.text",
  BLOCKS_POSITION: "blocks.position",
  SIZE_HEIGHT: "dimensions.height",
  SIZE_WIDTH: "dimensions.width",
};

const tokenKey = "token";
export const tokenSet = (token) => localStorage.setItem(tokenKey, token);
export const tokenGet = () => localStorage.getItem(tokenKey);
export const tokenHas = () => !!localStorage.getItem(tokenKey);

const socketGuesses = {
  "https:": "wss:",
  "http:": "ws:",
};

const socketProtocol = socketGuesses[window.location.protocol];
const socketHost = window.location.host;

export const newSocketAuth = (params) => newSocket({ ...params, token: tokenGet() });
export const newSocket = (params) => {
  const paramsEnc = new URLSearchParams(params);
  return new WebSocket(`${socketProtocol}//${socketHost}${urlSocket}?${paramsEnc}`);
};
