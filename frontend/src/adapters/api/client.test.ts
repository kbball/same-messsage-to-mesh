import { describe, it, expect, vi, beforeEach } from 'vitest'
import { get, post, put } from './client'

describe('API client', () => {
  beforeEach(() => {
    vi.resetAllMocks()
  })

  it('get makes a GET request and returns JSON', async () => {
    const mockData = { count: 5 }
    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve(mockData),
    } as Response)

    const result = await get<typeof mockData>('/api/test')
    expect(result).toEqual(mockData)
    expect(fetch).toHaveBeenCalledWith('/api/test', {
      method: 'GET',
      headers: undefined,
      body: undefined,
    })
  })

  it('post sends JSON body', async () => {
    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ updated: 1 }),
    } as Response)

    await post('/api/test', { key: 'val' })
    const call = (fetch as ReturnType<typeof vi.fn>).mock.calls[0]
    expect(call[1].method).toBe('POST')
    expect(call[1].body).toBe('{"key":"val"}')
  })

  it('put sends JSON body', async () => {
    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({}),
    } as Response)

    await put('/api/test', { x: 1 })
    const call = (fetch as ReturnType<typeof vi.fn>).mock.calls[0]
    expect(call[1].method).toBe('PUT')
  })

  it('throws on non-ok response', async () => {
    global.fetch = vi.fn().mockResolvedValue({
      ok: false,
      status: 500,
      statusText: 'Internal Server Error',
      text: () => Promise.resolve('server error'),
    } as unknown as Response)

    await expect(get('/api/bad')).rejects.toThrow('500')
  })
})
