export interface Mock {
  id: string
  name: string
  path: string
  method: string
  description?: string
  isActive: boolean
  responseConfig: ResponseConfig
  createdAt: string
  updatedAt: string
}

export interface ResponseConfig {
  type: "static" | "dynamic" | "faker"
  statusCode: number
  headers?: Record<string, string>
  body: unknown
  delayMs?: number
  faker?: FakerConfig
}

export interface FakerConfig {
  fields: FakerField[]
}

export interface FakerField {
  name: string
  type: string
  format?: string
}

export interface CreateMockRequest {
  name: string
  path: string
  method: string
  description?: string
  responseConfig: ResponseConfig
}

export interface UpdateMockRequest {
  name?: string
  path?: string
  method?: string
  description?: string
  isActive?: boolean
  responseConfig?: ResponseConfig
}

export interface TestMockRequest {
  input: Record<string, unknown>
}

export interface TestMockResponse {
  output: string
  statusCode: number
  headers?: Record<string, string>
  tookMs: number
}

export interface MockTest {
  id: string
  mockId: string
  input: string
  output: string
  statusCode: number
  tookMs: number
  createdAt: string
}

export interface ListResponse<T> {
  data: T[]
  total: number
  page: number
  pageSize: number
  totalPages: number
}

export interface ErrorResponse {
  error: string
  code?: string
  details?: string
}

export interface HealthResponse {
  status: string
}