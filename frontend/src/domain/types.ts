export interface SAMEAlert {
  id: number
  received_at: string
  originator: string
  event_code: string
  fips_codes: string[]
  purge_time: string
  issue_time: string
  call_sign: string
  raw_message: string
  published: boolean
}

export interface EventCode {
  code: string
  description: string
  category: string
  is_warning: boolean
  updated_at: string
}

export interface FIPSCode {
  state_code: string
  county_code: string
  state_name: string
  county_name: string
  updated_at: string
}

export interface AlertFilter {
  state_codes: string[]
  fips_codes: string[]
  event_codes: string[]
  updated_at: string
}

export interface SDRDeviceConfig {
  device_path: string
  frequency: number
  updated_at: string
}

export const NOAA_FREQUENCIES = [
  { label: '162.400 MHz', value: 162400000 },
  { label: '162.425 MHz', value: 162425000 },
  { label: '162.450 MHz', value: 162450000 },
  { label: '162.475 MHz', value: 162475000 },
  { label: '162.500 MHz', value: 162500000 },
  { label: '162.525 MHz', value: 162525000 },
  { label: '162.550 MHz (default)', value: 162550000 },
] as const
