import { useEffect, useState } from 'react'
import {
  Box,
  Typography,
  Paper,
  TextField,
  Button,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  CircularProgress,
  Alert,
  Stack,
} from '@mui/material'
import { api } from '../../adapters/api'
import { NOAA_FREQUENCIES } from '../../domain/types'
import type { SDRDeviceConfig } from '../../domain/types'

export default function SDRConfigTab() {
  const [config, setConfig] = useState<SDRDeviceConfig | null>(null)
  const [devicePath, setDevicePath] = useState('')
  const [frequency, setFrequency] = useState(162550000)
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null)

  useEffect(() => {
    api
      .getSDRConfig()
      .then((cfg) => {
        setConfig(cfg)
        setDevicePath(cfg.device_path)
        setFrequency(cfg.frequency)
      })
      .catch((e: Error) => setMessage({ type: 'error', text: e.message }))
      .finally(() => setLoading(false))
  }, [])

  const handleSave = async () => {
    setSaving(true)
    setMessage(null)
    try {
      const updated = await api.updateSDRConfig({ device_path: devicePath, frequency })
      setConfig(updated)
      setMessage({
        type: 'success',
        text: 'SDR config saved. Restart the server to apply changes.',
      })
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
        SDR Configuration
      </Typography>
      <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
        Configure the RTL-SDR device path and NOAA Weather Radio frequency. Changes take effect
        after the server restarts.
      </Typography>

      {message && (
        <Alert severity={message.type} sx={{ mb: 2 }} onClose={() => setMessage(null)}>
          {message.text}
        </Alert>
      )}

      <Paper sx={{ p: 3, maxWidth: 480 }}>
        <Stack spacing={3}>
          <TextField
            label="Device Path"
            value={devicePath}
            onChange={(e) => setDevicePath(e.target.value)}
            helperText="Linux device path, e.g. /dev/bus/usb or /dev/rtl0"
            fullWidth
          />

          <FormControl fullWidth>
            <InputLabel>NOAA Weather Radio Frequency</InputLabel>
            <Select
              value={frequency}
              onChange={(e) => setFrequency(Number(e.target.value))}
              label="NOAA Weather Radio Frequency"
            >
              {NOAA_FREQUENCIES.map((f) => (
                <MenuItem key={f.value} value={f.value}>
                  {f.label}
                </MenuItem>
              ))}
            </Select>
          </FormControl>

          {config && (
            <Typography variant="caption" color="text.secondary">
              Last updated: {new Date(config.updated_at).toLocaleString()}
            </Typography>
          )}

          <Button variant="contained" onClick={handleSave} disabled={saving || !devicePath}>
            {saving ? 'Saving…' : 'Save'}
          </Button>
        </Stack>
      </Paper>
    </Box>
  )
}
