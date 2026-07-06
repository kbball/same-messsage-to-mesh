import { useEffect, useState, useCallback } from 'react'
import {
  Box,
  Typography,
  Paper,
  Button,
  CircularProgress,
  Alert,
  Stack,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
} from '@mui/material'
import RefreshIcon from '@mui/icons-material/Refresh'
import { api } from '../../adapters/api'
import type { EventCode } from '../../domain/types'

export default function ReferenceDataTab() {
  const [fipsCount, setFIPSCount] = useState<number | null>(null)
  const [eventCodes, setEventCodes] = useState<EventCode[]>([])
  const [loading, setLoading] = useState(true)
  const [refreshingFIPS, setRefreshingFIPS] = useState(false)
  const [refreshingEC, setRefreshingEC] = useState(false)
  const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null)

  const loadData = useCallback(async () => {
    try {
      const [countResult, codes] = await Promise.all([api.fipsCount(), api.listEventCodes()])
      setFIPSCount(countResult.count)
      setEventCodes(codes)
    } catch (e) {
      setMessage({ type: 'error', text: (e as Error).message })
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    loadData()
  }, [loadData])

  const handleRefreshFIPS = async () => {
    setRefreshingFIPS(true)
    setMessage(null)
    try {
      const result = await api.refreshFIPS()
      setMessage({ type: 'success', text: `FIPS codes updated: ${result.updated} records` })
      const countResult = await api.fipsCount()
      setFIPSCount(countResult.count)
    } catch (e) {
      setMessage({ type: 'error', text: `FIPS refresh failed: ${(e as Error).message}` })
    } finally {
      setRefreshingFIPS(false)
    }
  }

  const handleRefreshEventCodes = async () => {
    setRefreshingEC(true)
    setMessage(null)
    try {
      const result = await api.refreshEventCodes()
      setMessage({ type: 'success', text: `Event codes updated: ${result.updated} records` })
      const codes = await api.listEventCodes()
      setEventCodes(codes)
    } catch (e) {
      setMessage({ type: 'error', text: `Event code refresh failed: ${(e as Error).message}` })
    } finally {
      setRefreshingEC(false)
    }
  }

  if (loading) return <CircularProgress sx={{ m: 4 }} />

  return (
    <Box>
      <Typography variant="h6" sx={{ mb: 2 }}>
        Reference Data
      </Typography>

      {message && (
        <Alert severity={message.type} sx={{ mb: 2 }} onClose={() => setMessage(null)}>
          {message.text}
        </Alert>
      )}

      <Stack spacing={3}>
        <Paper sx={{ p: 3 }}>
          <Stack direction="row" sx={{ alignItems: 'center' }} spacing={2}>
            <Box sx={{ flexGrow: 1 }}>
              <Typography variant="subtitle1">FIPS County Codes</Typography>
              <Typography variant="body2" color="text.secondary">
                {fipsCount === null ? '…' : `${fipsCount.toLocaleString()} counties loaded`}
              </Typography>
              <Typography variant="caption" color="text.secondary">
                Source: NOAA NCEI ERDDAP
              </Typography>
            </Box>
            <Button
              variant="outlined"
              startIcon={refreshingFIPS ? <CircularProgress size={16} /> : <RefreshIcon />}
              onClick={handleRefreshFIPS}
              disabled={refreshingFIPS}
            >
              {refreshingFIPS ? 'Refreshing…' : 'Refresh from NOAA'}
            </Button>
          </Stack>
        </Paper>

        <Paper sx={{ p: 3 }}>
          <Stack direction="row" sx={{ alignItems: 'center', mb: 2 }} spacing={2}>
            <Box sx={{ flexGrow: 1 }}>
              <Typography variant="subtitle1">SAME Event Codes</Typography>
              <Typography variant="body2" color="text.secondary">
                {eventCodes.length} codes loaded
              </Typography>
              <Typography variant="caption" color="text.secondary">
                Source: NOAA NWS weather.gov/nwr/eventcodes
              </Typography>
            </Box>
            <Button
              variant="outlined"
              startIcon={refreshingEC ? <CircularProgress size={16} /> : <RefreshIcon />}
              onClick={handleRefreshEventCodes}
              disabled={refreshingEC}
            >
              {refreshingEC ? 'Refreshing…' : 'Refresh from NOAA'}
            </Button>
          </Stack>

          {eventCodes.length > 0 && (
            <TableContainer>
              <Table size="small">
                <TableHead>
                  <TableRow>
                    <TableCell>Code</TableCell>
                    <TableCell>Description</TableCell>
                    <TableCell>Category</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {eventCodes.map((ec) => (
                    <TableRow key={ec.code}>
                      <TableCell>
                        <strong>{ec.code}</strong>
                      </TableCell>
                      <TableCell>{ec.description}</TableCell>
                      <TableCell>{ec.category}</TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </TableContainer>
          )}
        </Paper>
      </Stack>
    </Box>
  )
}
