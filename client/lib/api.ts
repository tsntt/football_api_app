import type { Championship, Match, AuthResponse, BroadcastResponse } from "./types"

const API_BASE_URL = process.env.SERVER_URL || "http://localhost:4000"

class ApiClient {
  private getAuthHeaders() {
    const token = localStorage.getItem("auth_token")
    return token ? { Authorization: `Bearer ${token}` } : {}
  }

  async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`
    const config: RequestInit = {
      headers: {
        "Content-Type": "application/json",
        ...this.getAuthHeaders(),
        ...options.headers,
      },
      ...options,
    }

    console.log("API Request:", url, config)

    const response = await fetch(url, config)

    if (!response.ok) {
      const errorMessage = `API Error: ${response.status} ${response.statusText}`
      const error = new Error(errorMessage)

        // Add status code to error for retry logic
        ; (error as any).status = response.status

      throw error
    }

    return response.json()
  }

  // Auth endpoints
  async register(name: string, password: string) {
    return this.request<{ message: string }>("/api/v1/auth/register", {
      method: "POST",
      body: JSON.stringify({ name, password }),
    })
  }

  async login(name: string, password: string) {
    return this.request<AuthResponse>("/api/v1/auth/login", {
      method: "POST",
      body: JSON.stringify({ name, password }),
    })
  }

  async logout() {
    return this.request<{ message: string }>("/api/v1/auth/logout", {
      method: "POST",
    })
  }

  // Championship endpoints
  async getChampionships() {
    return this.request<Championship[]>("/api/v1/championship")
  }

  async getMatches(championshipId: number, team?: string, stage?: string) {
    const params = new URLSearchParams()
    if (team) params.append("team", team)
    if (stage) params.append("stage", stage)

    const query = params.toString() ? `?${params.toString()}` : ""
    return this.request<Match[]>(`/api/v1/championship/${championshipId}/matches${query}`)
  }

  // Admin endpoints
  async getAdminMatches() {
    return this.request<Match[]>("/api/v1/admin/")
  }

  async broadcastMatch(matchId: number) {
    return this.request<BroadcastResponse>(`/api/v1/admin/broadcast/${matchId}`, {
      method: "POST",
    })
  }
}

export const apiClient = new ApiClient()
