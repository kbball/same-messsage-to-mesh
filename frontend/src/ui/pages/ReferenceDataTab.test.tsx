import { render, screen, waitFor, fireEvent } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import { http, HttpResponse } from 'msw'
import { server } from '../../test/server'
import ReferenceDataTab from './ReferenceDataTab'

describe('ReferenceDataTab', () => {
  it('shows FIPS count and event codes after load', async () => {
    render(<ReferenceDataTab />)
    await waitFor(() => {
      expect(screen.getByText(/3,235 counties loaded/)).toBeInTheDocument()
    })
    expect(screen.getByText('RWT')).toBeInTheDocument()
  })

  it('calls refresh FIPS and shows result', async () => {
    render(<ReferenceDataTab />)
    await waitFor(() => expect(screen.getAllByText('Refresh from NOAA').length).toBeGreaterThan(0))
    const buttons = screen.getAllByText('Refresh from NOAA')
    fireEvent.click(buttons[0])
    await waitFor(() => {
      expect(screen.getByText(/FIPS codes updated/)).toBeInTheDocument()
    })
  })

  it('calls refresh event codes and shows result', async () => {
    render(<ReferenceDataTab />)
    await waitFor(() => expect(screen.getAllByText('Refresh from NOAA').length).toBeGreaterThan(0))
    const buttons = screen.getAllByText('Refresh from NOAA')
    fireEvent.click(buttons[buttons.length - 1])
    await waitFor(() => {
      expect(screen.getByText(/Event codes updated/)).toBeInTheDocument()
    })
  })

  it('shows error message when refresh FIPS fails', async () => {
    server.use(
      http.post('/api/reference/fips/refresh', () =>
        HttpResponse.json({ error: 'network error' }, { status: 500 }),
      ),
    )
    render(<ReferenceDataTab />)
    await waitFor(() => screen.getAllByText('Refresh from NOAA'))
    fireEvent.click(screen.getAllByText('Refresh from NOAA')[0])
    await waitFor(() => {
      expect(screen.getByText(/FIPS refresh failed/)).toBeInTheDocument()
    })
  })

  it('dismisses message alert', async () => {
    render(<ReferenceDataTab />)
    await waitFor(() => screen.getAllByText('Refresh from NOAA'))
    fireEvent.click(screen.getAllByText('Refresh from NOAA')[0])
    await waitFor(() => screen.getByText(/FIPS codes updated/))
    fireEvent.click(screen.getByTitle('Close'))
    await waitFor(() => {
      expect(screen.queryByText(/FIPS codes updated/)).not.toBeInTheDocument()
    })
  })
})
