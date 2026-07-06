import { useEffect, useState, useCallback } from 'react'
import {
  Box,
  Typography,
  Paper,
  Button,
  CircularProgress,
  Alert,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Chip,
  Stack,
  Divider,
} from '@mui/material'
import { api } from '../../adapters/api'
import type { FIPSCode, EventCode } from '../../domain/types'

interface FilterSectionProps {
  label: string
  selected: string[]
  labelFor: (v: string) => string
  onRemove: (v: string) => void
  children: React.ReactNode
}

function FilterSection({ label, selected, labelFor, onRemove, children }: FilterSectionProps) {
  return (
    <Paper sx={{ p: 2 }}>
      <Typography variant="subtitle2" sx={{ mb: 1 }}>
        {label}
      </Typography>
      {selected.length > 0 && (
        <>
          <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5, mb: 1.5 }}>
            {selected.map((v) => (
              <Chip key={v} label={labelFor(v)} size="small" onDelete={() => onRemove(v)} />
            ))}
          </Box>
          <Divider sx={{ mb: 1.5 }} />
        </>
      )}
      {children}
    </Paper>
  )
}

export default function FiltersTab() {
  const [states, setStates] = useState<FIPSCode[]>([])
  const [counties, setCounties] = useState<FIPSCode[]>([])
  const [eventCodes, setEventCodes] = useState<EventCode[]>([])
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null)

  const [selectedStates, setSelectedStates] = useState<string[]>([])
  const [selectedFIPS, setSelectedFIPS] = useState<string[]>([])
  const [selectedEvents, setSelectedEvents] = useState<string[]>([])

  const loadData = useCallback(async () => {
    try {
      const [f, s, ec] = await Promise.all([
        api.getFilter(),
        api.listStates(),
        api.listEventCodes(),
      ])
      setStates(s)
      setEventCodes(ec)
      setSelectedStates(f.state_codes ?? [])
      setSelectedFIPS(f.fips_codes ?? [])
      setSelectedEvents(f.event_codes ?? [])
    } catch (e) {
      setMessage({ type: 'error', text: (e as Error).message })
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    loadData()
  }, [loadData])

  useEffect(() => {
    if (selectedStates.length === 1) {
      api
        .listCounties(selectedStates[0])
        .then(setCounties)
        .catch(() => setCounties([]))
    } else {
      setCounties([])
      setSelectedFIPS([])
    }
  }, [selectedStates])

  const addState = (code: string) => {
    if (!selectedStates.includes(code)) setSelectedStates((prev) => [...prev, code])
  }
  const removeState = (code: string) => {
    setSelectedStates((prev) => prev.filter((v) => v !== code))
  }

  const addFIPS = (code: string) => {
    if (!selectedFIPS.includes(code)) setSelectedFIPS((prev) => [...prev, code])
  }
  const removeFIPS = (code: string) => {
    setSelectedFIPS((prev) => prev.filter((v) => v !== code))
  }

  const addEvent = (code: string) => {
    if (!selectedEvents.includes(code)) setSelectedEvents((prev) => [...prev, code])
  }
  const removeEvent = (code: string) => {
    setSelectedEvents((prev) => prev.filter((v) => v !== code))
  }

  const handleSave = async () => {
    setSaving(true)
    setMessage(null)
    try {
      await api.updateFilter({
        state_codes: selectedStates,
        fips_codes: selectedFIPS,
        event_codes: selectedEvents,
      })
      setMessage({ type: 'success', text: 'Filters saved.' })
    } catch (e) {
      setMessage({ type: 'error', text: (e as Error).message })
    } finally {
      setSaving(false)
    }
  }

  if (loading) return <CircularProgress sx={{ m: 4 }} />

  const stateName = (code: string) => states.find((s) => s.state_code === code)?.state_name ?? code
  const countyName = (fips: string) => {
    const c = counties.find((c) => c.state_code + c.county_code === fips)
    return c ? c.county_name : fips
  }
  const eventLabel = (code: string) => {
    const ec = eventCodes.find((e) => e.code === code)
    return ec ? `${code} – ${ec.description}` : code
  }

  const availableStates = states.filter((s) => !selectedStates.includes(s.state_code))
  const availableCounties = counties.filter(
    (c) => !selectedFIPS.includes(c.state_code + c.county_code),
  )
  const availableEvents = eventCodes.filter((e) => !selectedEvents.includes(e.code))

  return (
    <Box>
      <Typography variant="h6" sx={{ mb: 1 }}>
        Alert Filters
      </Typography>
      <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
        Leave any category empty to receive all alerts for that dimension.
      </Typography>

      {message && (
        <Alert severity={message.type} sx={{ mb: 2 }} onClose={() => setMessage(null)}>
          {message.text}
        </Alert>
      )}

      <Stack spacing={2}>
        <FilterSection
          label="States"
          selected={selectedStates}
          labelFor={stateName}
          onRemove={removeState}
        >
          <FormControl fullWidth size="small">
            <InputLabel>Add a state</InputLabel>
            <Select
              value=""
              label="Add a state"
              onChange={(e) => addState(e.target.value as string)}
              displayEmpty
            >
              {availableStates.length === 0 ? (
                <MenuItem disabled>
                  {states.length === 0
                    ? 'No states loaded — refresh Reference Data first'
                    : 'All states selected'}
                </MenuItem>
              ) : (
                availableStates.map((s) => (
                  <MenuItem key={s.state_code} value={s.state_code}>
                    {s.state_name}
                  </MenuItem>
                ))
              )}
            </Select>
          </FormControl>
        </FilterSection>

        {selectedStates.length === 1 && (
          <FilterSection
            label="Counties"
            selected={selectedFIPS}
            labelFor={countyName}
            onRemove={removeFIPS}
          >
            <FormControl fullWidth size="small">
              <InputLabel>Add a county</InputLabel>
              <Select
                value=""
                label="Add a county"
                onChange={(e) => addFIPS(e.target.value as string)}
              >
                {availableCounties.length === 0 ? (
                  <MenuItem disabled>
                    {counties.length === 0 ? 'Loading counties…' : 'All counties selected'}
                  </MenuItem>
                ) : (
                  availableCounties.map((c) => (
                    <MenuItem
                      key={c.state_code + c.county_code}
                      value={c.state_code + c.county_code}
                    >
                      {c.county_name}
                    </MenuItem>
                  ))
                )}
              </Select>
            </FormControl>
          </FilterSection>
        )}

        <FilterSection
          label="Event Types"
          selected={selectedEvents}
          labelFor={eventLabel}
          onRemove={removeEvent}
        >
          <FormControl fullWidth size="small">
            <InputLabel>Add an event type</InputLabel>
            <Select
              value=""
              label="Add an event type"
              onChange={(e) => addEvent(e.target.value as string)}
            >
              {availableEvents.length === 0 ? (
                <MenuItem disabled>All event types selected</MenuItem>
              ) : (
                availableEvents.map((ec) => (
                  <MenuItem key={ec.code} value={ec.code}>
                    {ec.code} – {ec.description}
                  </MenuItem>
                ))
              )}
            </Select>
          </FormControl>
        </FilterSection>

        <Box>
          <Button variant="contained" onClick={handleSave} disabled={saving}>
            {saving ? 'Saving…' : 'Save Filters'}
          </Button>
        </Box>
      </Stack>
    </Box>
  )
}
