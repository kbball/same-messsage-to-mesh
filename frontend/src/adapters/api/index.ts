import { get, post, put } from './client'
import type { SAMEAlert, EventCode, FIPSCode, AlertFilter, SDRDeviceConfig } from '../../domain/types'

export const api = {
  // Alerts
  listAlerts: (limit = 100) => get<SAMEAlert[]>(`/api/alerts?limit=${limit}`),

  // Filters
  getFilter: () => get<AlertFilter>('/api/filters'),
  updateFilter: (filter: Partial<AlertFilter>) => put<AlertFilter>('/api/filters', filter),

  // SDR config
  getSDRConfig: () => get<SDRDeviceConfig>('/api/sdr-config'),
  updateSDRConfig: (cfg: Partial<SDRDeviceConfig>) => put<SDRDeviceConfig>('/api/sdr-config', cfg),

  // Reference data
  listStates: () => get<FIPSCode[]>('/api/reference/states'),
  listCounties: (stateCode: string) => get<FIPSCode[]>(`/api/reference/counties/${stateCode}`),
  listEventCodes: () => get<EventCode[]>('/api/reference/event-codes'),
  refreshFIPS: () => post<{ updated: number }>('/api/reference/fips/refresh'),
  refreshEventCodes: () => post<{ updated: number }>('/api/reference/event-codes/refresh'),
  fipsCount: () => get<{ count: number }>('/api/reference/fips/count'),
}
