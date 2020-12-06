export const urlScreenshot = "/api/screenshot";
export const urlSearch = "/api/search";
export const urlStartImport = "/api/start_import";
export const urlAuthenticate = "/api/authenticate";
export const urlSocket = "/api/websocket";
export const urlAbout = "/api/about";
export const urlImportStatus = "/api/import_status";

const req = async (method: 'get' | 'post' | 'put', url: string, body?: object) => {
  const token = tokenGet();
  const response = await fetch(url, {
    method,
    body: JSON.stringify(body),
    headers: token
      ? { authorization: `bearer ${token}` }
      : {},
  });
  return await response.json();
};

export const reqSearch = (body: PayloadSearch): Promise<ResponseSearch<Screenshot>> =>
  req("post", urlSearch, body);

export const reqAuthenticate = (body: PayloadAuthenticate): Promise<ResponseAuthenticate> =>
  req("put", urlAuthenticate, body);

export const reqStartImport = (): Promise<ResponseStartImport> =>
  req("post", urlStartImport);

export const reqScreenshot = (id: string): Promise<ResponseSearch<Screenshot>> =>
  req("get", `${urlScreenshot}/${id}`);

export const reqAbout = (): Promise<ResponseAbout> =>
  req("get", urlAbout);

export const reqImportStatus = (): Promise<ResponseImportStatus> =>
  req("get", urlImportStatus);

const tokenKey = "token";
export const tokenSet = (token: string) => localStorage.setItem(tokenKey, token);
export const tokenGet = () => localStorage.getItem(tokenKey) || undefined;
export const tokenHas = () => !!localStorage.getItem(tokenKey);

const socketGuesses: { [key: string]: string } = {
  "https:": "wss:",
  "http:": "ws:",
};

const socketProtocol = socketGuesses[window.location.protocol];
const socketHost = window.location.host;

interface SocketParams {
  want_settings?: 0 | 1,
  want_screenshot_id?: string,
  token?: string,
}

export const newSocketAuth = (params: SocketParams) => newSocket({ ...params, token: tokenGet() });
export const newSocket = (params: SocketParams) => {
  // @ts-ignore
  const paramsEnc = new URLSearchParams(params);
  return new WebSocket(`${socketProtocol}//${socketHost}${urlSocket}?${paramsEnc}`);
};

export interface PayloadSearch {
  term: string
  query?: Query
  size: number
  from: number
  highlight?: {
    style: {}
    fields: Field[]
  }
  fields?: string[]
  facets?: {}
  explain?: boolean
  sort: string[]
  includeLocations?: boolean
  search_after?: {}
  search_before?: {}
}

export interface PayloadAuthenticate {
  username: string
  password: string
}

export interface Status {
  total: number
  failed: number
  successful: number
}

export interface Query {
  boost?: {}
  match_all?: {}
}

export interface Block {
  pos: number;
  start: number;
  end: number;
  array_positions: number[];
}

type Locations = { [key in Field]?: { [q: string]: Block[] } }
type Fragments  = { [key in Field]?: string[] }

export enum Field {
  TIMESTAMP = "timestamp",
  TAGS = "tags",
  BLOCKS_TEXT = "blocks.text",
  BLOCKS_POSITION = "blocks.position",
  SIZE_HEIGHT = "dimensions.height",
  SIZE_WIDTH = "dimensions.width",
  DOMINANT_COLOUR = "dominant_colour",
};

export type FieldSort = `-${Field}` | `${Field}`

export interface ScreenshotFields {
  [Field.BLOCKS_POSITION]: number[]
  [Field.BLOCKS_TEXT]: string | string[]
  [Field.SIZE_HEIGHT]: number
  [Field.SIZE_WIDTH]: number
  [Field.DOMINANT_COLOUR]: string
  [Field.TAGS]: string
  [Field.TIMESTAMP]: Date
}

export interface Hit<T, Sort> {
  index: string
  id: string
  score: number
  sort: Sort[]
  fields: T
  locations?: Locations
  fragments?: Fragments
}

export type Screenshot = Hit<ScreenshotFields, FieldSort>

export interface ResponseSearch<T> {
  status: Status
  request: PayloadSearch
  hits: T[]
  total_hits: number
  max_score: number
  took: number
  facets?: {}
}

export interface ResponseAuthenticate {
  token: string
}

export interface ResponseStartImport { }

export interface ResponseAbout {
  version: string
  screenshots_indexed: number
  api_key: string
  socket_clients: number
  import_path: string
  screenshots_path: string
}

export interface ResponseImportStatus {
  running: boolean
  errors: {
    error: string
    time: string
  }[]
  last_id: string
  count_processed: number
  count_total: number
}