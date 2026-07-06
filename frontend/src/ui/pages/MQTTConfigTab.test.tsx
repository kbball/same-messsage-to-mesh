import { render, screen, waitFor, fireEvent } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import { http, HttpResponse } from 'msw'
import { server } from '../../test/server'
import MQTTConfigTab from './MQTTConfigTab'

const enabledConfig = {
  enabled: true,
  host: 'broker.local',
  port: 1883,
  publish_topic: 'same/alerts',
  updated_at: new Date().toISOString(),
}

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

  it('enables test message button when config is enabled', async () => {
    server.use(
      http.get('/api/mqtt-config', () => HttpResponse.json(enabledConfig)),
      http.put('/api/mqtt-config', () => HttpResponse.json(enabledConfig)),
    )
    render(<MQTTConfigTab />)
    await waitFor(() => {
      // After load, toggle should be on and button should be enabled
      const testBtn = screen.getByText('Send Test Message')
      expect(testBtn).not.toBeDisabled()
    })
  })

  it('sends a test message and shows success', async () => {
    server.use(
      http.get('/api/mqtt-config', () => HttpResponse.json(enabledConfig)),
      http.put('/api/mqtt-config', () => HttpResponse.json(enabledConfig)),
    )
    render(<MQTTConfigTab />)
    await waitFor(() => screen.getByText('Send Test Message'))
    const testBtn = screen.getByText('Send Test Message')
    await waitFor(() => expect(testBtn).not.toBeDisabled())
    fireEvent.click(testBtn)
    await waitFor(() => {
      expect(screen.getByText('Test message published successfully.')).toBeInTheDocument()
    })
  })

  it('shows warning when save returns a warning', async () => {
    server.use(
      http.put('/api/mqtt-config', () =>
        HttpResponse.json({ warning: 'config saved but MQTT reconnect failed' }),
      ),
    )
    render(<MQTTConfigTab />)
    await waitFor(() => screen.getByText('Save'))
    fireEvent.click(screen.getByText('Save'))
    await waitFor(() => {
      expect(screen.getByText(/MQTT reconnect failed/)).toBeInTheDocument()
    })
  })

  it('updates host field value', async () => {
    render(<MQTTConfigTab />)
    await waitFor(() => screen.getByLabelText('Broker Host'))
    const hostInput = screen.getByLabelText('Broker Host')
    fireEvent.change(hostInput, { target: { value: 'newbroker.local' } })
    expect(hostInput).toHaveValue('newbroker.local')
  })

  it('updates publish topic field value', async () => {
    render(<MQTTConfigTab />)
    await waitFor(() => screen.getByLabelText('Publish Topic'))
    const topicInput = screen.getByLabelText('Publish Topic')
    fireEvent.change(topicInput, { target: { value: 'alerts/same' } })
    expect(topicInput).toHaveValue('alerts/same')
  })

  it('dismisses success message when alert is closed', async () => {
    render(<MQTTConfigTab />)
    await waitFor(() => screen.getByText('Save'))
    fireEvent.click(screen.getByText('Save'))
    await waitFor(() => screen.getByText('MQTT config saved.'))
    fireEvent.click(screen.getByTitle('Close'))
    await waitFor(() => {
      expect(screen.queryByText('MQTT config saved.')).not.toBeInTheDocument()
    })
  })
})
