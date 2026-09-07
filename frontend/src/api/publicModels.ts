import { apiClient } from './client'

export interface PublicModelDTO {
  name: string
  provider: string
  billing_mode: string
  currency?: string
  input_price?: number | null
  output_price?: number | null
  cache_read_price?: number | null
  cache_write_price?: number | null
  per_request_price?: number | null
  unit: string
  endpoint_types: string[]
  tags: string[]
  description: string
}

export async function listPublicModels(): Promise<PublicModelDTO[]> {
  const response = await apiClient.get<PublicModelDTO[]>('/public/models')
  return response.data
}
