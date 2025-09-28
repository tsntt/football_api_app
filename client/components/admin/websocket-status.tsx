"use client"

import { Badge } from "@/components/ui/badge"
import { Wifi, WifiOff, Loader2, AlertTriangle } from "lucide-react"

interface WebSocketStatusProps {
  status: "connecting" | "connected" | "disconnected" | "error"
  className?: string
}

export function WebSocketStatus({ status, className }: WebSocketStatusProps) {
  const statusConfig = {
    connecting: {
      icon: Loader2,
      label: "Conectando...",
      variant: "secondary" as const,
      className: "animate-spin",
    },
    connected: {
      icon: Wifi,
      label: "Conectado",
      variant: "default" as const,
      className: "",
    },
    disconnected: {
      icon: WifiOff,
      label: "Desconectado",
      variant: "secondary" as const,
      className: "",
    },
    error: {
      icon: AlertTriangle,
      label: "Erro de Conexão",
      variant: "destructive" as const,
      className: "",
    },
  }

  const config = statusConfig[status]
  const Icon = config.icon

  return (
    <Badge variant={config.variant} className={className}>
      <Icon className={`h-3 w-3 mr-1 ${config.className}`} />
      {config.label}
    </Badge>
  )
}
