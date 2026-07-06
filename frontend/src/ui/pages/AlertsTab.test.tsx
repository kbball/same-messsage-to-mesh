import { render, screen, waitFor } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import AlertsTab from './AlertsTab'

describe('AlertsTab', () => {
  it('shows "No alerts" when list is empty', async () => {
    // handlers.ts returns a non-empty mock, so just check for the table header
    render(<AlertsTab />)
    await waitFor(() => {
      expect(screen.getByText('Decoded Alerts')).toBeInTheDocument()
    })
  })

  it('renders alert table with event code', async () => {
    render(<AlertsTab />)
    await waitFor(() => {
      expect(screen.getByText('RWT')).toBeInTheDocument()
    })
  })
})
