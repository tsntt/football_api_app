"use client"

import { Button } from "@/components/ui/button"
import { LogOut, Settings } from "lucide-react"
import { useAuthStore } from "@/lib/auth"
import { useRouter } from "next/navigation"
import { toast } from "sonner"

export function DashboardHeader() {
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

  const handleAdminPanel = () => {
    router.push("/notifications")
  }

  return (
    <header className="border-b bg-card">
      <div className="container mx-auto px-4 py-4">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-foreground">Football Manager</h1>
            <p className="text-sm text-muted-foreground">Bem-vindo, {user?.name}</p>
          </div>
          <div className="flex items-center gap-2">
            {user?.isAdmin && (
              <Button variant="outline" size="sm" onClick={handleAdminPanel} className="gap-2 bg-transparent">
                <Settings className="h-4 w-4" />
                Enviar Notificações
              </Button>
            )}
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
