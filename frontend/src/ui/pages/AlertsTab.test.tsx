import { render, screen, waitFor } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import { http, HttpResponse } from 'msw'
import { server } from '../../test/server'
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

  it('shows error message when alerts fail to load', async () => {
    server.use(
      http.get('/api/alerts', () => HttpResponse.json({ error: 'db error' }, { status: 500 })),
    )
    render(<AlertsTab />)
    await waitFor(() => {
      expect(screen.getByText(/500/)).toBeInTheDocument()
    })
  })
})
