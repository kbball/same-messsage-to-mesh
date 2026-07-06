import { render, screen, waitFor } from '@testing-library/react'
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
})
