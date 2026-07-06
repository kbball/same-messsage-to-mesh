import '@testing-library/jest-dom'
import { server } from './server'
import { beforeAll, afterEach, afterAll, vi } from 'vitest'

// jsdom doesn't implement EventSource; stub it so components using useStream don't throw.
class MockEventSource {
  onmessage: ((e: MessageEvent) => void) | null = null
  onerror: ((e: Event) => void) | null = null
  close() {}
}
vi.stubGlobal('EventSource', MockEventSource)

beforeAll(() => server.listen({ onUnhandledRequest: 'warn' }))
afterEach(() => server.resetHandlers())
afterAll(() => server.close())
