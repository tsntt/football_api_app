export interface User {
  id: number
  name: string
  isAdmin?: boolean
}

export interface Championship {
  id: number
  name: string
  code: string
  type: string
  emblem: string
  currentSeason: {
    id: number
    startDate: string
    endDate: string
    currentMatchday: number
    winner: null | string
  }
}

export interface Team {
  id: number
  name: string
  shortName: string
  tla: string
  crest: string
}

export interface Match {
  id: number
  utcDate: string
  status: string
  matchday: number
  stage: string
  group: string
  lastUpdated: string
  homeTeam: Team
  awayTeam: Team
  score: {
    winner: string
    duration: string
    fullTime: {
      home: number | null
      away: number | null
    }
    halfTime: {
      home: number | null
      away: number | null
    }
  }
  competition: Championship
  season: {
    id: number
    startDate: string
    endDate: string
    currentMatchday: number
    winner: null | string
  }
}

export interface AuthResponse {
  token: string
}

export interface BroadcastResponse {
  message: string
  data: {
    match_id: number
    message: string
    notification_id: string
    targets_count: number
  }
}

export interface WebSocketUpdate {
  channel_id: number
  total_sent: number
  sent_count: number
  failed_count: number
  is_completed: boolean
  error_details: string[]
}
