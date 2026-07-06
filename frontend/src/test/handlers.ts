import { http, HttpResponse } from 'msw'
import type {
  SAMEAlert,
  AlertFilter,
  SDRDeviceConfig,
  EventCode,
  FIPSCode,
  MQTTConfig,
} from '../domain/types'

const mockFilter: AlertFilter = {
  state_codes: [],
  fips_codes: [],
  event_codes: [],
  updated_at: new Date().toISOString(),
}

const mockSDRConfig: SDRDeviceConfig = {
  device_path: '/dev/bus/usb',
  frequency: 162550000,
  updated_at: new Date().toISOString(),
}

const mockMQTTConfig: MQTTConfig = {
  enabled: false,
  host: 'localhost',
  port: 1883,
  publish_topic: 'same/alerts',
  updated_at: new Date().toISOString(),
}

const mockAlerts: SAMEAlert[] = [
  {
    id: 1,
    received_at: new Date().toISOString(),
    originator: 'WXR',
    event_code: 'RWT',
    fips_codes: ['013121'],
    purge_time: '0030',
    issue_time: '1820218',
    call_sign: 'KRAH/NWS',
    raw_message: 'ZCZC-WXR-RWT-013121+0030-1820218-KRAH/NWS-',
    published: false,
  },
]

const mockEventCodes: EventCode[] = [
  {
    code: 'RWT',
    description: 'Required Weekly Test',
    category: 'Test',
    is_warning: false,
    updated_at: new Date().toISOString(),
  },
  {
    code: 'TOR',
    description: 'Tornado Warning',
    category: 'Warning',
    is_warning: true,
    updated_at: new Date().toISOString(),
  },
]

const mockStates: FIPSCode[] = [
  {
    state_code: '13',
    county_code: '000',
    state_name: 'Georgia',
    county_name: '',
    updated_at: new Date().toISOString(),
  },
]

export const handlers = [
  http.get('/api/alerts', () => HttpResponse.json(mockAlerts)),
  http.get('/api/filters', () => HttpResponse.json(mockFilter)),
  http.put('/api/filters', () => HttpResponse.json(mockFilter)),
  http.get('/api/sdr-config', () => HttpResponse.json(mockSDRConfig)),
  http.put('/api/sdr-config', () => HttpResponse.json(mockSDRConfig)),
  http.get('/api/reference/states', () => HttpResponse.json(mockStates)),
  http.get('/api/reference/counties/:stateCode', () => HttpResponse.json([])),
  http.get('/api/reference/event-codes', () => HttpResponse.json(mockEventCodes)),
  http.post('/api/reference/fips/refresh', () => HttpResponse.json({ updated: 3235 })),
  http.post('/api/reference/event-codes/refresh', () => HttpResponse.json({ updated: 62 })),
  http.get('/api/reference/fips/count', () => HttpResponse.json({ count: 3235 })),
  http.get('/api/mqtt-config', () => HttpResponse.json(mockMQTTConfig)),
  http.put('/api/mqtt-config', () => HttpResponse.json({ ...mockMQTTConfig, enabled: true })),
  http.post('/api/mqtt-config/test', () => HttpResponse.json({ ok: true })),
]
