"use client"

import { useMemo } from "react"
import { AdminMatchCard } from "./admin-match-card"
import type { Match } from "@/lib/types"

interface AdminMatchesListProps {
  matches: Match[]
  isLoading?: boolean
}

const statusOrder = ["LIVE", "IN_PLAY", "SCHEDULED", "FINISHED", "PAUSED", "POSTPONED", "SUSPENDED", "CANCELLED"]

export function AdminMatchesList({ matches, isLoading = false }: AdminMatchesListProps) {
  const sortedMatches = useMemo(() => {
    return [...matches].sort((a, b) => {
      // First sort by status priority
      const aStatusIndex = statusOrder.indexOf(a.status)
      const bStatusIndex = statusOrder.indexOf(b.status)

      if (aStatusIndex !== bStatusIndex) {
        return (
          (aStatusIndex === -1 ? statusOrder.length : aStatusIndex) -
          (bStatusIndex === -1 ? statusOrder.length : bStatusIndex)
        )
      }

      // Then sort by date (most recent first for live/finished, earliest first for scheduled)
      const aDate = new Date(a.utcDate).getTime()
      const bDate = new Date(b.utcDate).getTime()

      if (a.status === "LIVE" || a.status === "IN_PLAY" || a.status === "FINISHED") {
        return bDate - aDate // Most recent first
      } else {
        return aDate - bDate // Earliest first
      }
    })
  }, [matches])

  if (isLoading) {
    return (
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {[...Array(6)].map((_, index) => (
          <div key={index} className="h-48 bg-muted rounded animate-pulse" />
        ))}
      </div>
    )
  }

  if (matches.length === 0) {
    return (
      <div className="text-center py-12">
        <div className="text-6xl mb-4">📢</div>
        <h3 className="text-lg font-semibold mb-2">Nenhuma partida com inscritos</h3>
        <p className="text-muted-foreground">Não há partidas com usuários inscritos para notificar no momento.</p>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-xl font-semibold">Partidas com Inscritos ({matches.length})</h2>
        <div className="text-sm text-muted-foreground">Atualizado automaticamente a cada 30 segundos</div>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {sortedMatches.map((match) => (
          <AdminMatchCard key={match.id} match={match} />
        ))}
      </div>
    </div>
  )
}
