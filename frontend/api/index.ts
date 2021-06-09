import router from '../router'

export const urlMedia = '/api/media'
export const urlSearch = '/api/search'
export const urlStartImport = '/api/start_import'
export const urlAuthenticate = '/api/authenticate'
export const urlSocket = '/api/websocket'
export const urlAbout = '/api/about'
export const urlDirectories = '/api/directories'
export const urlImportStatus = '/api/import_status'
export const urlPing = '/api/ping'
export const urlUpload = '/api/upload'

const tokenKey = 'token'
export const tokenSet = (token: string) => localStorage.setItem(tokenKey, token)
export const tokenGet = () => localStorage.getItem(tokenKey) || undefined
export const tokenHas = () => !!localStorage.getItem(tokenKey)

export interface Error {
  error: string
}

export interface Success<T> {
  result: T
}

export type Reponse<T> = Promise<Success<T> | Error>

export const isError = <T>(r: Success<T> | Error): r is Error => (r as Error).error !== undefined

type ReqMethod = 'get' | 'post' | 'put'
const req = async <P, R>(method: ReqMethod, url: string, data?: P): Reponse<R> => {
  const token = tokenGet()

  let headers: HeadersInit = {}
  if (token) headers.authorization = `bearer ${token}`

  let body: BodyInit = ''
  if (data instanceof FormData) body = data
  else body = JSON.stringify(data)

  const response = await fetch(url, { method, body, headers })
  if (response?.status === 401) {
    router.push({ name: 'login' })
  }

  return await response.json()
}

export enum SortOrder {
  Asc = 'asc',
  Desc = 'desc',
}

export interface PayloadSort {
  field: string
  order: SortOrder
}

export interface PayloadSearch {
  body: string
  limit: number
  offset: number
  sort: PayloadSort
  directory?: string
  media?: MediaType
}

export const reqSearch = (data: PayloadSearch) => {
  return req<PayloadSearch, Search>('post', urlSearch, data)
}

export interface PayloadAuthenticate {
  username: string
  password: string
}

export const reqAuthenticate = (data: PayloadAuthenticate) => {
  return req<PayloadAuthenticate, Authenticate>('put', urlAuthenticate, data)
}

export const reqStartImport = () => {
  return req<{}, StartImport>('post', urlStartImport)
}

export const reqMedia = (id: string) => {
  return req<{}, Media>('get', `${urlMedia}/${id}`)
}

export const reqAbout = () => {
  return req<{}, About>('get', urlAbout)
}

export const reqDirectories = () => {
  return req<{}, Directory[]>('get', urlDirectories)
}

export const reqImportStatus = () => {
  return req<{}, ImportStatus>('get', urlImportStatus)
}

export const reqPing = () => {
  return req<{}, {}>('get', urlPing)
}

export const reqUpload = (data: FormData) => {
  return req<FormData, Upload>('post', urlUpload, data)
}

const socketGuesses: { [key: string]: string } = {
  'https:': 'wss:',
  'http:': 'ws:',
}

const socketProtocol = socketGuesses[window.location.protocol]
const socketHost = window.location.host

interface SocketParams {
  want_settings?: 0 | 1
  want_media_hash?: string
  token?: string
}

export const newSocketAuth = (params: SocketParams) => newSocket({ ...params, token: tokenGet() })
export const newSocket = (params: SocketParams) => {
  // @ts-ignore
  const paramsEnc = new URLSearchParams(params)
  return new WebSocket(`${socketProtocol}//${socketHost}${urlSocket}?${paramsEnc}`)
}

export interface Block {
  id: number
  media_id: number
  index: number
  min_x: number
  min_y: number
  max_x: number
  max_y: number
  body: string
}

export enum MediaType {
  Image = 'image',
  Video = 'video',
}

export interface Media {
  id: number
  type: MediaType
  mime: string
  hash: string
  timestamp: any
  dim_width: number
  dim_height: number
  dominant_colour: string
  blurhash: string
  blocks?: Block[]
  highlighted_blocks?: Block[]
  directories?: string[]
}

export interface Similarity {
  similarity: number
}

export interface Search {
  medias?: (Media & Similarity)[]
  took: number
}

export interface Authenticate {
  token: string
}

export interface StartImport {}

export type About = { [key: string]: number | string }

export interface Directory {
  directory_alias: string
  count: number
  is_uploads?: boolean
}

export interface ImportStatus {
  running: boolean
  errors: {
    error: string
    time: string
  }[]
  last_hash: string
  count_processed: number
  count_total: number
}

export interface Upload {
  id: string
}
