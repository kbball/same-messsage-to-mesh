import { render, screen, waitFor, fireEvent } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import SDRConfigTab from './SDRConfigTab'

describe('SDRConfigTab', () => {
  it('loads and shows device path', async () => {
    render(<SDRConfigTab />)
    await waitFor(() => {
      const input = screen.getByLabelText('Device Path')
      expect(input).toHaveValue('/dev/bus/usb')
    })
  })

  it('shows save button', async () => {
    render(<SDRConfigTab />)
    await waitFor(() => {
      expect(screen.getByText('Save')).toBeInTheDocument()
    })
  })

  it('shows RTL-SDR description text', async () => {
    render(<SDRConfigTab />)
    await waitFor(() => {
      expect(screen.getByText(/RTL-SDR device path/)).toBeInTheDocument()
    })
  })

  it('saves config and shows success message', async () => {
    render(<SDRConfigTab />)
    await waitFor(() => screen.getByText('Save'))
    fireEvent.click(screen.getByText('Save'))
    await waitFor(() => {
      expect(screen.getByText(/SDR config saved/)).toBeInTheDocument()
    })
  })

  it('updates device path field', async () => {
    render(<SDRConfigTab />)
    await waitFor(() => screen.getByLabelText('Device Path'))
    const input = screen.getByLabelText('Device Path')
    fireEvent.change(input, { target: { value: '/dev/rtl1' } })
    expect(input).toHaveValue('/dev/rtl1')
  })

  it('dismisses success alert when closed', async () => {
    render(<SDRConfigTab />)
    await waitFor(() => screen.getByText('Save'))
    fireEvent.click(screen.getByText('Save'))
    await waitFor(() => screen.getByText(/SDR config saved/))
    fireEvent.click(screen.getByTitle('Close'))
    await waitFor(() => {
      expect(screen.queryByText(/SDR config saved/)).not.toBeInTheDocument()
    })
  })
})
