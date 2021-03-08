import router from "../router"

export const urlScreenshot = "/api/screenshot";
export const urlSearch = "/api/search";
export const urlStartImport = "/api/start_import";
export const urlAuthenticate = "/api/authenticate";
export const urlSocket = "/api/websocket";
export const urlAbout = "/api/about";
export const urlImportStatus = "/api/import_status";
export const urlPing = "/api/ping";

const tokenKey = "token";
export const tokenSet = (token: string) => localStorage.setItem(tokenKey, token);
export const tokenGet = () => localStorage.getItem(tokenKey) || undefined;
export const tokenHas = () => !!localStorage.getItem(tokenKey);

export interface Error {
  error: string
}

export interface Success<T> {
  result: T
}

export type Reponse<T> = Promise<Success<T> | Error>

export const isError = <T>(r: Success<T> | Error): r is Error =>
  (r as Error).error !== undefined

type ReqMethod = 'get' | 'post' | 'put'
const req = async <P, R>(method: ReqMethod, url: string, body?: P): Reponse<R> => {
  const token = tokenGet();

  const response = await fetch(url, {
    method,
    body: JSON.stringify(body),
    headers: token
      ? { authorization: `bearer ${token}` }
      : {},
  });

  if (response?.status === 401) {
    router.push(({ name: "login" }))
  }

  const json = await response.json();
  return json
};

export interface PayloadSearch {
  term: string
  size: number
  from: number
}

export const reqSearch = (body: PayloadSearch) =>
  req<PayloadSearch, Search>("post", urlSearch, body);

export interface PayloadAuthenticate {
  username: string
  password: string
}

export const reqAuthenticate = (body: PayloadAuthenticate) =>
  req<PayloadAuthenticate, Authenticate>("put", urlAuthenticate, body);

export const reqStartImport = () =>
  req<{}, StartImport>("post", urlStartImport);

export const reqScreenshot = (id: string) =>
  req<{}, Screenshot>("get", `${urlScreenshot}/${id}`);

export const reqAbout = () =>
  req<{}, About>("get", urlAbout);

export const reqImportStatus = () =>
  req<{}, ImportStatus>("get", urlImportStatus);

export const reqPing = () =>
  req<{}, {}>("get", urlPing); 

const socketGuesses: { [key: string]: string } = {
  "https:": "wss:",
  "http:": "ws:",
};

const socketProtocol = socketGuesses[window.location.protocol];
const socketHost = window.location.host;

interface SocketParams {
  want_settings?: 0 | 1,
  want_screenshot_hash?: string,
  token?: string,
}

export const newSocketAuth = (params: SocketParams) => newSocket({ ...params, token: tokenGet() });
export const newSocket = (params: SocketParams) => {
  // @ts-ignore
  const paramsEnc = new URLSearchParams(params);
  return new WebSocket(`${socketProtocol}//${socketHost}${urlSocket}?${paramsEnc}`);
};

export interface Block {
  id: number
  screenshot_id: number
  index: number
  min_x: number
  min_y: number
  max_x: number
  max_y: number
  body: string
}

export interface Screenshot {
  id: number
  hash: string
  timestamp: any
  dim_width: number
  dim_height: number
  dominant_colour: string
  blurhash: string
  blocks: Block[]
}

export interface Similarity {
  similarity: number
}

export type Search = (Screenshot & Similarity)[]

export interface Authenticate {
  token: string
}

export interface StartImport { }

export interface About {
  version: string
  screenshots_indexed: number
  api_key: string
  socket_clients: number
  import_path: string
  screenshots_path: string
}

export interface ImportStatus {
  running: boolean
  errors: {
    error: string
    time: string
  }[]
  last_id: string
  count_processed: number
  count_total: number
}