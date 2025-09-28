"use client"

import { Card, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Calendar, Clock } from "lucide-react"
import type { Match } from "@/lib/types"
import { format } from "date-fns"
import { ptBR } from "date-fns/locale"

interface MatchCardProps {
  match: Match
}

const statusLabels: Record<string, { label: string; variant: "default" | "secondary" | "destructive" | "outline" }> = {
  SCHEDULED: { label: "Agendado", variant: "outline" },
  LIVE: { label: "Ao Vivo", variant: "destructive" },
  IN_PLAY: { label: "Em Andamento", variant: "destructive" },
  PAUSED: { label: "Pausado", variant: "secondary" },
  FINISHED: { label: "Finalizado", variant: "default" },
  POSTPONED: { label: "Adiado", variant: "secondary" },
  SUSPENDED: { label: "Suspenso", variant: "secondary" },
  CANCELLED: { label: "Cancelado", variant: "secondary" },
}

const stageLabels: Record<string, string> = {
  REGULAR_SEASON: "Temporada Regular",
  GROUP_STAGE: "Fase de Grupos",
  ROUND_OF_16: "Oitavas de Final",
  QUARTER_FINALS: "Quartas de Final",
  SEMI_FINALS: "Semifinais",
  FINAL: "Final",
  THIRD_PLACE: "Terceiro Lugar",
}

export function MatchCard({ match }: MatchCardProps) {
  const matchDate = new Date(match.utcDate)
  const statusInfo = statusLabels[match.status] || { label: match.status, variant: "outline" as const }

  const homeScore = match.score.fullTime.home
  const awayScore = match.score.fullTime.away
  const hasScore = homeScore !== null && awayScore !== null

  return (
    <Card className="hover:shadow-md transition-shadow">
      <CardContent className="p-4">
        <div className="flex items-center justify-between mb-3">
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <Calendar className="h-4 w-4" />
            {format(matchDate, "dd/MM/yyyy", { locale: ptBR })}
            <Clock className="h-4 w-4 ml-2" />
            {format(matchDate, "HH:mm", { locale: ptBR })}
          </div>
          <Badge variant={statusInfo.variant}>{statusInfo.label}</Badge>
        </div>

        <div className="flex items-center justify-between">
          {/* Home Team */}
          <div className="flex items-center gap-3 flex-1">
            <img
              src={match.homeTeam.crest || "/placeholder.svg?height=32&width=32&query=team logo"}
              alt={match.homeTeam.name}
              className="w-8 h-8 object-contain"
            />
            <div className="min-w-0 flex-1">
              <p className="font-medium text-sm truncate">{match.homeTeam.name}</p>
              <p className="text-xs text-muted-foreground truncate">{match.homeTeam.shortName}</p>
            </div>
          </div>

          {/* Score */}
          <div className="flex items-center gap-4 px-4">
            {hasScore ? (
              <div className="text-center">
                <div className="text-2xl font-bold">
                  {homeScore} - {awayScore}
                </div>
                {match.score.halfTime.home !== null && match.score.halfTime.away !== null && (
                  <div className="text-xs text-muted-foreground">
                    ({match.score.halfTime.home} - {match.score.halfTime.away})
                  </div>
                )}
              </div>
            ) : (
              <div className="text-2xl font-bold text-muted-foreground">VS</div>
            )}
          </div>

          {/* Away Team */}
          <div className="flex items-center gap-3 flex-1 justify-end">
            <div className="min-w-0 flex-1 text-right">
              <p className="font-medium text-sm truncate">{match.awayTeam.name}</p>
              <p className="text-xs text-muted-foreground truncate">{match.awayTeam.shortName}</p>
            </div>
            <img
              src={match.awayTeam.crest || "/placeholder.svg?height=32&width=32&query=team logo"}
              alt={match.awayTeam.name}
              className="w-8 h-8 object-contain"
            />
          </div>
        </div>

        {/* Match Details */}
        <div className="mt-3 pt-3 border-t flex items-center justify-between text-xs text-muted-foreground">
          <span>{stageLabels[match.stage] || match.stage}</span>
          <span>Rodada {match.matchday}</span>
        </div>
      </CardContent>
    </Card>
  )
}
