import { render, screen, waitFor, fireEvent } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import { http, HttpResponse } from 'msw'
import { server } from '../../test/server'
import FiltersTab from './FiltersTab'

describe('FiltersTab', () => {
  it('loads and shows section headings', async () => {
    render(<FiltersTab />)
    await waitFor(() => {
      expect(screen.getByText('States')).toBeInTheDocument()
      expect(screen.getByText('Event Types')).toBeInTheDocument()
    })
  })

  it('shows Alert Filters heading and instructions', async () => {
    render(<FiltersTab />)
    await waitFor(() => {
      expect(screen.getByText('Alert Filters')).toBeInTheDocument()
      expect(screen.getByText(/Leave any category empty/)).toBeInTheDocument()
    })
  })

  it('shows Save Filters button', async () => {
    render(<FiltersTab />)
    await waitFor(() => {
      expect(screen.getByText('Save Filters')).toBeInTheDocument()
    })
  })

  it('saves filters and shows success message', async () => {
    render(<FiltersTab />)
    await waitFor(() => screen.getByText('Save Filters'))
    fireEvent.click(screen.getByText('Save Filters'))
    await waitFor(() => {
      expect(screen.getByText('Filters saved.')).toBeInTheDocument()
    })
  })

  it('adds a state chip when a state option is clicked', async () => {
    render(<FiltersTab />)
    await waitFor(() => screen.getByText('States'))

    // Open the State select (first combobox on the page)
    const selects = screen.getAllByRole('combobox')
    fireEvent.mouseDown(selects[0])

    const georgiaOption = await screen.findByRole('option', { name: /Georgia/ })
    fireEvent.click(georgiaOption)

    await waitFor(() => {
      expect(screen.getByText('Georgia')).toBeInTheDocument()
    })
  })

  it('adds an event code chip when an event option is clicked', async () => {
    render(<FiltersTab />)
    await waitFor(() => screen.getByText('Event Types'))

    // Event Types is the second combobox
    const selects = screen.getAllByRole('combobox')
    fireEvent.mouseDown(selects[selects.length - 1])

    const torOption = await screen.findByRole('option', { name: /Required Weekly Test/ })
    fireEvent.click(torOption)

    await waitFor(() => {
      // Chip label contains the event code
      expect(screen.getByText(/RWT/)).toBeInTheDocument()
    })
  })

  it('removes an event code chip when delete is clicked', async () => {
    render(<FiltersTab />)
    await waitFor(() => screen.getByText('Event Types'))

    // Add an event first
    const selects = screen.getAllByRole('combobox')
    fireEvent.mouseDown(selects[selects.length - 1])
    await screen.findByRole('option', { name: /Required Weekly Test/ })
    fireEvent.click(screen.getByRole('option', { name: /Required Weekly Test/ }))
    await waitFor(() => screen.getByText(/RWT/))

    // Remove it via chip delete button
    const chipDeleteBtn = screen.getByTestId('CancelIcon')
    fireEvent.click(chipDeleteBtn)

    await waitFor(() => {
      expect(screen.queryByText(/RWT – Required Weekly Test/)).not.toBeInTheDocument()
    })
  })

  it('dismisses success message when alert is closed', async () => {
    render(<FiltersTab />)
    await waitFor(() => screen.getByText('Save Filters'))
    fireEvent.click(screen.getByText('Save Filters'))
    await waitFor(() => screen.getByText('Filters saved.'))
    fireEvent.click(screen.getByTitle('Close'))
    await waitFor(() => {
      expect(screen.queryByText('Filters saved.')).not.toBeInTheDocument()
    })
  })

  it('shows county section when exactly one state is selected', async () => {
    render(<FiltersTab />)
    await waitFor(() => screen.getByText('States'))

    const selects = screen.getAllByRole('combobox')
    fireEvent.mouseDown(selects[0])
    await screen.findByRole('option', { name: /Georgia/ })
    fireEvent.click(screen.getByRole('option', { name: /Georgia/ }))

    await waitFor(() => {
      expect(screen.getByText('Counties')).toBeInTheDocument()
    })
  })

  it('allows adding a county when counties are available', async () => {
    server.use(
      http.get('/api/reference/counties/:stateCode', () =>
        HttpResponse.json([
          {
            state_code: '13',
            county_code: '121',
            state_name: 'Georgia',
            county_name: 'Fulton County',
            updated_at: new Date().toISOString(),
          },
        ]),
      ),
    )
    render(<FiltersTab />)
    await waitFor(() => screen.getByText('States'))

    // Select a state to show the county section
    const selects = screen.getAllByRole('combobox')
    fireEvent.mouseDown(selects[0])
    await screen.findByRole('option', { name: /Georgia/ })
    fireEvent.click(screen.getByRole('option', { name: /Georgia/ }))

    await waitFor(() => screen.getByText('Counties'))

    // Open county select
    const updatedSelects = screen.getAllByRole('combobox')
    const countySelect = updatedSelects[1]
    fireEvent.mouseDown(countySelect)

    const fultonOption = await screen.findByRole('option', { name: /Fulton County/ })
    fireEvent.click(fultonOption)

    await waitFor(() => {
      expect(screen.getByText('Fulton County')).toBeInTheDocument()
    })
  })

  it('shows no states loaded message when states list is empty', async () => {
    server.use(http.get('/api/reference/states', () => HttpResponse.json([])))
    render(<FiltersTab />)
    await waitFor(() => screen.getByText('States'))

    const selects = screen.getAllByRole('combobox')
    fireEvent.mouseDown(selects[0])

    await waitFor(() => {
      expect(screen.getByText(/No states loaded/)).toBeInTheDocument()
    })
  })
})
