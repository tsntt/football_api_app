"use client"

import type React from "react"

import { useState } from "react"
import { useRouter } from "next/navigation"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { useAuthStore } from "@/lib/auth"
import { toast } from "sonner"

interface LoginFormProps {
  onToggleMode: () => void
  isRegister: boolean
}

export function LoginForm({ onToggleMode, isRegister }: LoginFormProps) {
  const [name, setname] = useState("")
  const [password, setPassword] = useState("")
  const [isLoading, setIsLoading] = useState(false)
  const { login, register } = useAuthStore()
  const router = useRouter()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!name || !password) {
      toast.error("Por favor, preencha todos os campos")
      return
    }

    setIsLoading(true)
    try {
      if (isRegister) {
        await register(name, password)
        toast.success("Conta criada com sucesso!")
      } else {
        await login(name, password)
        toast.success("Login realizado com sucesso!")
      }
      router.push("/dashboard")
    } catch (error) {
      toast.error(isRegister ? "Erro ao criar conta" : "Erro ao fazer login")
      console.error("Auth error:", error)
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <Card className="w-full max-w-md">
      <CardHeader className="text-center">
        <CardTitle className="text-2xl font-bold">{isRegister ? "Criar Conta" : "Entrar"}</CardTitle>
        <CardDescription>
          {isRegister ? "Crie sua conta para acessar o sistema" : "Entre com suas credenciais para continuar"}
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="name">Usuário</Label>
            <Input
              id="name"
              type="text"
              value={name}
              onChange={(e) => setname(e.target.value)}
              placeholder="Digite seu usuário"
              disabled={isLoading}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="password">Senha</Label>
            <Input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Digite sua senha"
              disabled={isLoading}
            />
          </div>
          <Button type="submit" className="w-full" disabled={isLoading}>
            {isLoading ? "Carregando..." : isRegister ? "Criar Conta" : "Entrar"}
          </Button>
        </form>
        <div className="mt-4 text-center">
          <Button variant="link" onClick={onToggleMode} disabled={isLoading} className="text-sm">
            {isRegister ? "Já tem uma conta? Faça login" : "Não tem conta? Registre-se"}
          </Button>
        </div>
      </CardContent>
    </Card>
  )
}
