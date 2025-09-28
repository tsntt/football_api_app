"use client"

import { useMemo } from "react"
import { MatchCard } from "./match-card"
import type { Match } from "@/lib/types"

interface MatchesListProps {
  matches: Match[]
  isLoading?: boolean
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

export function MatchesList({ matches, isLoading = false }: MatchesListProps) {
  const groupedMatches = useMemo(() => {
    const groups: Record<string, Match[]> = {}

    console.log("Matches:", matches)

    if (Array.isArray(matches)) {
      matches.forEach((match) => {
        if (match) {
          const stage = match.stage
          if (!groups[stage]) {
            groups[stage] = []
          }
          groups[stage].push(match)
        }
      })
    }

    // Sort matches within each group by date
    Object.keys(groups).forEach((stage) => {
      groups[stage].sort((a, b) => new Date(a.utcDate).getTime() - new Date(b.utcDate).getTime())
    })

    return groups
  }, [matches])

  if (isLoading) {
    return (
      <div className="space-y-6">
        {[...Array(2)].map((_, groupIndex) => (
          <div key={groupIndex} className="space-y-4">
            <div className="h-6 bg-muted rounded w-48 animate-pulse" />
            <div className="space-y-3">
              {[...Array(3)].map((_, cardIndex) => (
                <div key={cardIndex} className="h-32 bg-muted rounded animate-pulse" />
              ))}
            </div>
          </div>
        ))}
      </div>
    )
  }

  if (!matches || matches.length === 0) {
    return (
      <div className="text-center py-12">
        <div className="text-6xl mb-4">⚽</div>
        <h3 className="text-lg font-semibold mb-2">Nenhuma partida encontrada</h3>
        <p className="text-muted-foreground">
          Tente ajustar os filtros para encontrar as partidas que você está procurando.
        </p>
      </div>
    )
  }

  const stageOrder = [
    "GROUP_STAGE",
    "ROUND_OF_16",
    "QUARTER_FINALS",
    "SEMI_FINALS",
    "THIRD_PLACE",
    "FINAL",
    "REGULAR_SEASON",
  ]

  const sortedStages = Object.keys(groupedMatches).sort((a, b) => {
    const aIndex = stageOrder.indexOf(a)
    const bIndex = stageOrder.indexOf(b)

    if (aIndex === -1 && bIndex === -1) return a.localeCompare(b)
    if (aIndex === -1) return 1
    if (bIndex === -1) return -1

    return aIndex - bIndex
  })

  return (
    <div className="space-y-8">
      {sortedStages.map((stage) => (
        <div key={stage} className="space-y-4">
          <div className="flex items-center gap-2">
            <h2 className="text-xl font-semibold text-foreground">{stageLabels[stage] || stage}</h2>
            <span className="text-sm text-muted-foreground bg-muted px-2 py-1 rounded">
              {groupedMatches[stage].length} partida{groupedMatches[stage].length !== 1 ? "s" : ""}
            </span>
          </div>
          <div className="grid gap-4">
            {groupedMatches[stage].map((match) => (
              <MatchCard key={match.id} match={match} />
            ))}
          </div>
        </div>
      ))}
    </div>
  )
}
