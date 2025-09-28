import { create } from "zustand"
import { persist } from "zustand/middleware"
import type { User } from "./types"
import { apiClient } from "./api"

interface AuthState {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  isInitialized: boolean
  login: (name: string, password:string) => Promise<void>
  register: (name: string, password: string) => Promise<void>
  logout: () => Promise<void>
  setUser: (user: User) => void
  initializeAuth: () => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      token: null,
      isAuthenticated: false,
      isInitialized: false,

      initializeAuth: () => {
        const { token, user } = get()
        if (token && user) {
          set({ isAuthenticated: true })
        }
        set({ isInitialized: true })
      },

      login: async (name: string, password: string) => {
        try {
          const response = await apiClient.login(name, password)
          localStorage.setItem("auth_token", response.token)

          // Mock user data - in real app, you'd get this from token or separate endpoint
          // TODO: use real user data
          const user: User = {
            id: 1,
            name,
            isAdmin: name === "admin", // Simple admin check
          }

          set({
            user,
            token: response.token,
            isAuthenticated: true,
          })
        } catch (error) {
          throw error
        }
      },

      register: async (name: string, password: string) => {
        try {
          await apiClient.register(name, password)
          // Auto-login after registration
          await get().login(name, password)
        } catch (error) {
          throw error
        }
      },

      logout: async () => {
        try {
          await apiClient.logout()
        } catch (error) {
          // Continue with logout even if API call fails
          console.error("Logout API error:", error)
        } finally {
          localStorage.removeItem("auth_token")
          set({
            user: null,
            token: null,
            isAuthenticated: false,
          })
        }
      },

      setUser: (user: User) => {
        set({ user })
      },
    }),
    {
      name: "auth-storage",
      partialize: (state) => ({
        user: state.user,
        token: state.token,
      }),
    },
  ),
)