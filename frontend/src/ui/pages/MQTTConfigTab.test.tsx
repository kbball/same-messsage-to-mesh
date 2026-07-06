import { render, screen, waitFor, fireEvent } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import MQTTConfigTab from './MQTTConfigTab'

describe('MQTTConfigTab', () => {
  it('loads and shows broker host', async () => {
    render(<MQTTConfigTab />)
    await waitFor(() => {
      const input = screen.getByLabelText('Broker Host')
      expect(input).toHaveValue('localhost')
    })
  })

  it('shows save button', async () => {
    render(<MQTTConfigTab />)
    await waitFor(() => {
      expect(screen.getByText('Save')).toBeInTheDocument()
    })
  })

  it('send test message button is disabled when not enabled', async () => {
    render(<MQTTConfigTab />)
    await waitFor(() => {
      const testBtn = screen.getByText('Send Test Message')
      expect(testBtn).toBeDisabled()
    })
  })

  it('saves config and shows success message', async () => {
    render(<MQTTConfigTab />)
    await waitFor(() => screen.getByText('Save'))
    fireEvent.click(screen.getByText('Save'))
    await waitFor(() => {
      expect(screen.getByText('MQTT config saved.')).toBeInTheDocument()
    })
  })
})
