import { render, screen, waitFor, fireEvent } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
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
})
