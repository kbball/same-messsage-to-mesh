import { useEffect, useState } from 'react'
import {
  Box,
  Typography,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Chip,
  CircularProgress,
} from '@mui/material'
import { api } from '../../adapters/api'
import { useStream } from '../../adapters/sse/useStream'
import type { SAMEAlert } from '../../domain/types'

export default function AlertsTab() {
  const [alerts, setAlerts] = useState<SAMEAlert[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const liveAlert = useStream<SAMEAlert>('alert')

  useEffect(() => {
    api
      .listAlerts(200)
      .then(setAlerts)
      .catch((e: Error) => setError(e.message))
      .finally(() => setLoading(false))
  }, [])

  useEffect(() => {
    if (liveAlert) {
      setAlerts((prev) => [liveAlert, ...prev])
    }
  }, [liveAlert])

  if (loading) return <CircularProgress sx={{ m: 4 }} />
  if (error) return <Typography color="error">{error}</Typography>

  return (
    <Box>
      <Typography variant="h6" sx={{ mb: 2 }}>
        Decoded Alerts
      </Typography>
      {alerts.length === 0 ? (
        <Typography color="text.secondary">No alerts received yet.</Typography>
      ) : (
        <TableContainer component={Paper}>
          <Table size="small">
            <TableHead>
              <TableRow>
                <TableCell>Received</TableCell>
                <TableCell>Event</TableCell>
                <TableCell>Originator</TableCell>
                <TableCell>FIPS Codes</TableCell>
                <TableCell>Call Sign</TableCell>
                <TableCell>Published</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {alerts.map((a) => (
                <TableRow key={a.id}>
                  <TableCell>{new Date(a.received_at).toLocaleString()}</TableCell>
                  <TableCell>
                    <Chip
                      label={a.event_code}
                      size="small"
                      color={
                        a.event_code === 'RWT' || a.event_code === 'RMT' ? 'default' : 'warning'
                      }
                    />
                  </TableCell>
                  <TableCell>{a.originator}</TableCell>
                  <TableCell>{a.fips_codes.join(', ')}</TableCell>
                  <TableCell>{a.call_sign}</TableCell>
                  <TableCell>
                    {a.published ? (
                      <Chip label="Yes" size="small" color="success" />
                    ) : (
                      <Chip label="No" size="small" />
                    )}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      )}
    </Box>
  )
}
