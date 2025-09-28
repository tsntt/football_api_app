"use client"

import { useEffect, useRef, useState } from "react"
import type { WebSocketUpdate } from "@/lib/types"

interface UseWebSocketOptions {
  url: string
  enabled?: boolean
  onMessage?: (data: WebSocketUpdate) => void
  onError?: (error: Event) => void
  onConnect?: () => void
  onDisconnect?: () => void
}

export function useWebSocket({
  url,
  enabled = true,
  onMessage,
  onError,
  onConnect,
  onDisconnect,
}: UseWebSocketOptions) {
  const [isConnected, setIsConnected] = useState(false)
  const [connectionStatus, setConnectionStatus] = useState<"connecting" | "connected" | "disconnected" | "error">(
    "disconnected",
  )
  const wsRef = useRef<WebSocket | null>(null)
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null)
  const reconnectAttempts = useRef(0)
  const maxReconnectAttempts = 5

  const connect = () => {
    if (!enabled || wsRef.current?.readyState === WebSocket.OPEN) return

    try {
      setConnectionStatus("connecting")
      wsRef.current = new WebSocket(url)

      wsRef.current.onopen = () => {
        setIsConnected(true)
        setConnectionStatus("connected")
        reconnectAttempts.current = 0
        onConnect?.()
      }

      wsRef.current.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data) as WebSocketUpdate
          onMessage?.(data)
        } catch (error) {
          console.error("Failed to parse WebSocket message:", error)
        }
      }

      wsRef.current.onclose = () => {
        setIsConnected(false)
        setConnectionStatus("disconnected")
        onDisconnect?.()

        // Attempt to reconnect
        if (enabled && reconnectAttempts.current < maxReconnectAttempts) {
          reconnectAttempts.current++
          const delay = Math.min(1000 * Math.pow(2, reconnectAttempts.current), 30000) // Exponential backoff, max 30s

          reconnectTimeoutRef.current = setTimeout(() => {
            connect()
          }, delay)
        }
      }

      wsRef.current.onerror = (error) => {
        setConnectionStatus("error")
        onError?.(error)
      }
    } catch (error) {
      setConnectionStatus("error")
      console.error("WebSocket connection error:", error)
    }
  }

  const disconnect = () => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current)
      reconnectTimeoutRef.current = null
    }

    if (wsRef.current) {
      wsRef.current.close()
      wsRef.current = null
    }

    setIsConnected(false)
    setConnectionStatus("disconnected")
  }

  useEffect(() => {
    if (enabled) {
      connect()
    } else {
      disconnect()
    }

    return () => {
      disconnect()
    }
  }, [enabled, url])

  return {
    isConnected,
    connectionStatus,
    connect,
    disconnect,
  }
}
