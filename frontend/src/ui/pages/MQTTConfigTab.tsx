import { useEffect, useState } from 'react'
import {
  Box,
  Typography,
  Paper,
  TextField,
  Button,
  CircularProgress,
  Alert,
  Stack,
  FormControlLabel,
  Switch,
} from '@mui/material'
import SendIcon from '@mui/icons-material/Send'
import { api } from '../../adapters/api'
import type { MQTTConfig } from '../../domain/types'

export default function MQTTConfigTab() {
  const [config, setConfig] = useState<MQTTConfig | null>(null)
  const [enabled, setEnabled] = useState(false)
  const [host, setHost] = useState('localhost')
  const [port, setPort] = useState(1883)
  const [publishTopic, setPublishTopic] = useState('same/alerts')
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [testing, setTesting] = useState(false)
  const [message, setMessage] = useState<{ type: 'success' | 'error' | 'warning'; text: string } | null>(null)

  useEffect(() => {
    api
      .getMQTTConfig()
      .then((cfg) => {
        setConfig(cfg)
        setEnabled(cfg.enabled)
        setHost(cfg.host)
        setPort(cfg.port)
        setPublishTopic(cfg.publish_topic)
      })
      .catch((e: Error) => setMessage({ type: 'error', text: e.message }))
      .finally(() => setLoading(false))
  }, [])

  const handleSave = async () => {
    setSaving(true)
    setMessage(null)
    try {
      const result = await api.updateMQTTConfig({ enabled, host, port, publish_topic: publishTopic })
      if ('warning' in result) {
        setMessage({ type: 'warning', text: (result as { warning: string }).warning })
      } else {
        setConfig(result as MQTTConfig)
        setMessage({ type: 'success', text: 'MQTT config saved.' })
      }
    } catch (e) {
      setMessage({ type: 'error', text: (e as Error).message })
    } finally {
      setSaving(false)
    }
  }

  const handleTest = async () => {
    setTesting(true)
    setMessage(null)
    try {
      await api.testMQTTPublish()
      setMessage({ type: 'success', text: 'Test message published successfully.' })
    } catch (e) {
      setMessage({ type: 'error', text: (e as Error).message })
    } finally {
      setTesting(false)
    }
  }

  if (loading) return <CircularProgress sx={{ m: 4 }} />

  return (
    <Box>
      <Typography variant="h6" sx={{ mb: 2 }}>
        MQTT Configuration
      </Typography>
      <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
        Configure the MQTT broker for publishing decoded alerts to the mesh network.
        Save changes to reconnect with the new settings.
      </Typography>

      {message && (
        <Alert severity={message.type} sx={{ mb: 2 }} onClose={() => setMessage(null)}>
          {message.text}
        </Alert>
      )}

      <Paper sx={{ p: 3, maxWidth: 480 }}>
        <Stack spacing={3}>
          <FormControlLabel
            control={
              <Switch
                checked={enabled}
                onChange={(e) => setEnabled(e.target.checked)}
                color="primary"
              />
            }
            label="Enable MQTT publishing"
          />

          <TextField
            label="Broker Host"
            value={host}
            onChange={(e) => setHost(e.target.value)}
            helperText="Hostname or IP address of the MQTT broker"
            fullWidth
          />

          <TextField
            label="Broker Port"
            type="number"
            value={port}
            onChange={(e) => setPort(Number(e.target.value))}
            slotProps={{ htmlInput: { min: 1, max: 65535 } }}
            fullWidth
          />

          <TextField
            label="Publish Topic"
            value={publishTopic}
            onChange={(e) => setPublishTopic(e.target.value)}
            helperText="MQTT topic alerts will be published to"
            fullWidth
          />

          {config && (
            <Typography variant="caption" color="text.secondary">
              Last updated: {new Date(config.updated_at).toLocaleString()}
            </Typography>
          )}

          <Stack direction="row" spacing={2}>
            <Button variant="contained" onClick={handleSave} disabled={saving || !host}>
              {saving ? 'Saving…' : 'Save'}
            </Button>
            <Button
              variant="outlined"
              startIcon={<SendIcon />}
              onClick={handleTest}
              disabled={testing || !enabled || !config?.enabled}
            >
              {testing ? 'Sending…' : 'Send Test Message'}
            </Button>
          </Stack>
        </Stack>
      </Paper>
    </Box>
  )
}
