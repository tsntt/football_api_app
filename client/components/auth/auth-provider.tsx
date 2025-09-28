"use client"

import { useEffect, type ReactNode } from "react"
import { useAuthStore } from "@/lib/auth"

interface AuthProviderProps {
  children: ReactNode
}

export function AuthProvider({ children }: AuthProviderProps) {
  const { initializeAuth, isInitialized } = useAuthStore()

  useEffect(() => {
    if (!isInitialized) {
      initializeAuth()
    }
  }, [initializeAuth, isInitialized])

  // Show loading state while initializing auth
  if (!isInitialized) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
      </div>
    )
  }

  return <>{children}</>
}
