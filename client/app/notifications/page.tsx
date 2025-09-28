"use client"

import { useEffect, useState } from "react"
import { useRouter } from "next/navigation"
import { useAuthStore } from "@/lib/auth"
import { useAdminMatches } from "@/hooks/use-admin-matches"
import { useWebSocket } from "@/hooks/use-websocket"
import { AdminHeader } from "@/components/admin/admin-header"
import { AdminMatchesList } from "@/components/admin/admin-matches-list"
import { BroadcastProgress } from "@/components/admin/broadcast-progress"
import { Alert, AlertDescription } from "@/components/ui/alert"
import { AlertTriangle } from "lucide-react"
import { toast } from "sonner"
import type { WebSocketUpdate } from "@/lib/types"

const WS_URL = process.env.NEXT_PUBLIC_WS_URL || "ws://localhost:4000/api/v1/admin/"

export default function NotificacoesPage() {
  const { isAuthenticated, user } = useAuthStore()
  const router = useRouter()
  const { data: matches = [], isLoading, error } = useAdminMatches()
  const [broadcastUpdates, setBroadcastUpdates] = useState<WebSocketUpdate[]>([])

  const handleWebSocketMessage = (data: WebSocketUpdate) => {
    setBroadcastUpdates((prev) => {
      // Update existing or add new
      const existingIndex = prev.findIndex((update) => update.channel_id === data.channel_id)
      if (existingIndex >= 0) {
        const newUpdates = [...prev]
        newUpdates[existingIndex] = data
        return newUpdates
      } else {
        return [...prev, data]
      }
    })

    // Show toast notifications for completed broadcasts
    if (data.is_completed) {
      if (data.failed_count > 0) {
        toast.error(`Notificação concluída com ${data.failed_count} falhas`)
      } else {
        toast.success(`Notificação enviada para ${data.sent_count} usuários`)
      }
    }
  }

  const { connectionStatus } = useWebSocket({
    url: WS_URL,
    enabled: isAuthenticated && user?.isAdmin,
    onMessage: handleWebSocketMessage,
    onConnect: () => {
      toast.success("Conectado ao sistema de notificações em tempo real")
    },
    onDisconnect: () => {
      toast.warning("Desconectado do sistema de notificações")
    },
    onError: () => {
      toast.error("Erro na conexão WebSocket")
    },
  })

  useEffect(() => {
    if (!isAuthenticated) {
      router.push("/")
      return
    }

    if (!user?.isAdmin) {
      router.push("/dashboard")
      return
    }
  }, [isAuthenticated, user, router])

  // Clear completed broadcasts after 30 seconds
  useEffect(() => {
    const interval = setInterval(() => {
      setBroadcastUpdates((prev) => prev.filter((update) => !update.is_completed || Date.now() - Date.now() < 30000))
    }, 30000)

    return () => clearInterval(interval)
  }, [])

  if (!isAuthenticated || !user?.isAdmin) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <h2 className="text-xl font-semibold mb-2">Acesso Negado</h2>
          <p className="text-muted-foreground">Você precisa ser um administrador para acessar esta página.</p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background">
      <AdminHeader wsStatus={connectionStatus} />

      <main className="container mx-auto px-4 py-6 space-y-6">
        {error && (
          <Alert variant="destructive">
            <AlertTriangle className="h-4 w-4" />
            <AlertDescription>Erro ao carregar as partidas. Tente novamente em alguns instantes.</AlertDescription>
          </Alert>
        )}

        <BroadcastProgress updates={broadcastUpdates} />

        <AdminMatchesList matches={matches} isLoading={isLoading} />
      </main>
    </div>
  )
}
