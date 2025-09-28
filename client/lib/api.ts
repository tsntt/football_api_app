import type { Championship, Match, AuthResponse, BroadcastResponse } from "./types"
import { toast } from "sonner"

const API_BASE_URL = process.env.SERVER_URL || "http://localhost:4000/api/v1"

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
        ...(this.getAuthHeaders() as Record<string, string>),
        ...(options.headers as Record<string, string>),
      } as Record<string, string>,
      ...options,
    }

    console.log("API Request:", url, config)

    const response = await fetch(url, config)

    if (!response.ok) {
      const errorMessage = `API Error: ${response.status} ${response.statusText}`
      const error = new Error(errorMessage);

        // Add status code to error for retry logic
        (error as any).status = response.status

        if (endpoint.endsWith('/fans') && response.status === 429) {
          const resp = await response.json()
          toast.warning(resp.message)
        }

      throw error
    }

    return response.json()
  }

  // Auth endpoints
  async register(name: string, password: string) {
    return this.request<{ message: string }>("/auth/register", {
      method: "POST",
      body: JSON.stringify({ name, password }),
    })
  }

  async login(name: string, password: string) {
    return this.request<AuthResponse>("/auth/login", {
      method: "POST",
      body: JSON.stringify({ name, password }),
    })
  }

  async logout() {
    return this.request<{ message: string }>("/auth/logout", {
      method: "POST",
    })
  }

  // Championship endpoints
  async getChampionships() {
    return this.request<Championship[]>("/championship")
  }

  async getMatches(championshipId: number, team?: string, stage?: string) {
    const params = new URLSearchParams()
    if (team) params.append("team", team)
    if (stage) params.append("stage", stage)

    const query = params.toString() ? `?${params.toString()}` : ""
    return this.request<Match[]>(`/championship/${championshipId}/matches${query}`)
  }

    async subscribeToTeam(teamId: number, teamName: string) {
      return this.request<{ message: string }>("/fans", {
        method: "POST",
        body: JSON.stringify({ team_id: teamId, team_name: teamName }),
      })
    }
  // Admin endpoints
  async getAdminMatches() {
    return this.request<Match[]>("/admin/")
  }

  async broadcastMatch(matchId: number) {
    return this.request<BroadcastResponse>(`/admin/broadcast/${matchId}`, {
      method: "POST",
    })
  }
}

export const apiClient = new ApiClient()
