"use client"

import { useState } from "react"
import { LoginForm } from "@/components/auth/login-form"

export default function HomePage() {
  const [isRegister, setIsRegister] = useState(false)

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100 dark:from-gray-900 dark:to-gray-800 p-4">
      <div className="w-full max-w-md">
        <div className="text-center mb-8">
          <h1 className="text-4xl font-bold text-foreground mb-2">Football Manager</h1>
          <p className="text-muted-foreground">Sistema de gerenciamento de campeonatos</p>
        </div>
        <LoginForm onToggleMode={() => setIsRegister(!isRegister)} isRegister={isRegister} />
      </div>
    </div>
  )
}
