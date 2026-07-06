import { useEffect, useState } from 'react'

interface SSEEnvelope<T> {
  type: string
  payload: T
}

export function useStream<T>(eventType: string): T | undefined {
  const [data, setData] = useState<T | undefined>(undefined)

  useEffect(() => {
    const es = new EventSource('/api/stream')

    es.onmessage = (e: MessageEvent) => {
      try {
        const envelope = JSON.parse(e.data as string) as SSEEnvelope<T>
        if (envelope.type === eventType) {
          setData(envelope.payload)
        }
      } catch {
        // ignore malformed events
      }
    }

    return () => es.close()
  }, [eventType])

  return data
}
