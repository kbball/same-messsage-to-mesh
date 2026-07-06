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
  OutlinedInput,
  Chip,
  Stack,
} from '@mui/material'
import { api } from '../../adapters/api'
import type { AlertFilter, FIPSCode, EventCode } from '../../domain/types'

export default function FiltersTab() {
  const [filter, setFilter] = useState<AlertFilter | null>(null)
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
      const [f, s, ec] = await Promise.all([api.getFilter(), api.listStates(), api.listEventCodes()])
      setFilter(f)
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

  const handleSave = async () => {
    setSaving(true)
    setMessage(null)
    try {
      await api.updateFilter({
        state_codes: selectedStates,
        fips_codes: selectedFIPS,
        event_codes: selectedEvents,
      })
      setMessage({ type: 'success', text: 'Filter saved.' })
    } catch (e) {
      setMessage({ type: 'error', text: (e as Error).message })
    } finally {
      setSaving(false)
    }
  }

  if (loading) return <CircularProgress sx={{ m: 4 }} />

  return (
    <Box>
      <Typography variant="h6" sx={{ mb: 2 }}>
        Alert Filters
      </Typography>
      <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
        Leave any dimension empty to receive all alerts for that category.
      </Typography>

      {message && (
        <Alert severity={message.type} sx={{ mb: 2 }} onClose={() => setMessage(null)}>
          {message.text}
        </Alert>
      )}

      <Stack spacing={3}>
        <Paper sx={{ p: 2 }}>
          <FormControl fullWidth>
            <InputLabel>States</InputLabel>
            <Select
              multiple
              value={selectedStates}
              onChange={(e) => setSelectedStates(e.target.value as string[])}
              input={<OutlinedInput label="States" />}
              renderValue={(selected) => (
                <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                  {(selected as string[]).map((v) => {
                    const s = states.find((x) => x.state_code === v)
                    return <Chip key={v} label={s ? s.state_name : v} size="small" />
                  })}
                </Box>
              )}
            >
              {states.length === 0 ? (
                <MenuItem disabled>No states loaded — refresh Reference Data first</MenuItem>
              ) : (
                states.map((s) => (
                  <MenuItem key={s.state_code} value={s.state_code}>
                    {s.state_name}
                  </MenuItem>
                ))
              )}
            </Select>
          </FormControl>
        </Paper>

        {selectedStates.length === 1 && counties.length > 0 && (
          <Paper sx={{ p: 2 }}>
            <FormControl fullWidth>
              <InputLabel>Counties (optional)</InputLabel>
              <Select
                multiple
                value={selectedFIPS}
                onChange={(e) => setSelectedFIPS(e.target.value as string[])}
                input={<OutlinedInput label="Counties (optional)" />}
                renderValue={(selected) => (
                  <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                    {(selected as string[]).map((v) => {
                      const c = counties.find((x) => x.state_code + x.county_code === v)
                      return <Chip key={v} label={c ? c.county_name : v} size="small" />
                    })}
                  </Box>
                )}
              >
                {counties.map((c) => (
                  <MenuItem key={c.state_code + c.county_code} value={c.state_code + c.county_code}>
                    {c.county_name}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Paper>
        )}

        <Paper sx={{ p: 2 }}>
          <FormControl fullWidth>
            <InputLabel>Event Types</InputLabel>
            <Select
              multiple
              value={selectedEvents}
              onChange={(e) => setSelectedEvents(e.target.value as string[])}
              input={<OutlinedInput label="Event Types" />}
              renderValue={(selected) => (
                <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                  {(selected as string[]).map((v) => {
                    const ec = eventCodes.find((x) => x.code === v)
                    return <Chip key={v} label={ec ? `${v} – ${ec.description}` : v} size="small" />
                  })}
                </Box>
              )}
            >
              {eventCodes.map((ec) => (
                <MenuItem key={ec.code} value={ec.code}>
                  {ec.code} – {ec.description}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </Paper>

        <Box>
          <Button variant="contained" onClick={handleSave} disabled={saving || !filter}>
            {saving ? 'Saving…' : 'Save Filters'}
          </Button>
        </Box>
      </Stack>
    </Box>
  )
}
