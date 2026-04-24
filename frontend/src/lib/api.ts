import axios from 'axios'
import type {
  Mock,
  CreateMockRequest,
  UpdateMockRequest,
  TestMockRequest,
  TestMockResponse,
  ListResponse,
  HealthResponse,
} from '@/types/api'

const apiBaseURL = '/api/v1'

const api = axios.create({
  baseURL: apiBaseURL,
  headers: {
    'Content-Type': 'application/json',
  },
})

export const mocksApi = {
  list: async (): Promise<Mock[]> => {
    const { data } = await api.get<ListResponse<Mock>>('/mocks')
    return data.data
  },

  get: async (id: string): Promise<Mock> => {
    const { data } = await api.get<Mock>(`/mocks/${id}`)
    return data
  },

  create: async (mock: CreateMockRequest): Promise<Mock> => {
    const { data } = await api.post<Mock>('/mocks', mock)
    return data
  },

  update: async (id: string, mock: UpdateMockRequest): Promise<Mock> => {
    const { data } = await api.put<Mock>(`/mocks/${id}`, mock)
    return data
  },

  delete: async (id: string): Promise<void> => {
    await api.delete(`/mocks/${id}`)
  },

  activate: async (id: string): Promise<Mock> => {
    const { data } = await api.post<Mock>(`/mocks/${id}/activate`)
    return data
  },

  deactivate: async (id: string): Promise<Mock> => {
    const { data } = await api.post<Mock>(`/mocks/${id}/deactivate`)
    return data
  },

  test: async (id: string, req: TestMockRequest): Promise<TestMockResponse> => {
    const { data } = await api.post<TestMockResponse>(`/mocks/${id}/test`, req)
    return data
  },
}

export const healthApi = {
  check: async (): Promise<HealthResponse> => {
    const { data } = await api.get<HealthResponse>('/health')
    return data
  },
}

export default api