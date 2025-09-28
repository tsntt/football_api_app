"use client"

import { Button } from "@/components/ui/button"
import { ArrowLeft, LogOut } from "lucide-react"
import { useAuthStore } from "@/lib/auth"
import { useRouter } from "next/navigation"
import { toast } from "sonner"
import { WebSocketStatus } from "./websocket-status"

interface AdminHeaderProps {
  wsStatus?: "connecting" | "connected" | "disconnected" | "error"
}

export function AdminHeader({ wsStatus = "disconnected" }: AdminHeaderProps) {
  const { user, logout } = useAuthStore()
  const router = useRouter()

  const handleLogout = async () => {
    try {
      await logout()
      toast.success("Logout realizado com sucesso!")
      router.push("/")
    } catch (error) {
      toast.error("Erro ao fazer logout")
    }
  }

  const handleBackToDashboard = () => {
    router.push("/dashboard")
  }

  return (
    <header className="border-b bg-card">
      <div className="container mx-auto px-4 py-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4">
            <Button variant="outline" size="sm" onClick={handleBackToDashboard} className="gap-2 bg-transparent">
              <ArrowLeft className="h-4 w-4" />
              Voltar
            </Button>
            <div>
              <h1 className="text-2xl font-bold text-foreground">Painel de Notificações</h1>
              <p className="text-sm text-muted-foreground">Administrador: {user?.name}</p>
            </div>
          </div>
          <div className="flex items-center gap-3">
            <WebSocketStatus status={wsStatus} />
            <Button variant="outline" size="sm" onClick={handleLogout} className="gap-2 bg-transparent">
              <LogOut className="h-4 w-4" />
              Sair
            </Button>
          </div>
        </div>
      </div>
    </header>
  )
}
