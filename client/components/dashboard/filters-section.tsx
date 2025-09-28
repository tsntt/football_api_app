"use client"

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { ChampionshipSelector } from "./championship-selector"
import { TeamSelector } from "./team-selector"
import { StageSelector } from "./stage-selector"
import type { Championship, Match } from "@/lib/types"

interface FiltersSectionProps {
  championships: Championship[]
  matches: Match[]
  selectedChampionship: number | null
  selectedTeam: string
  selectedStage: string
  onChampionshipChange: (value: number | null) => void
  onTeamChange: (value: string) => void
  onStageChange: (value: string) => void
  isLoadingChampionships: boolean
  isLoadingMatches: boolean
}

export function FiltersSection({
  championships,
  matches,
  selectedChampionship,
  selectedTeam,
  selectedStage,
  onChampionshipChange,
  onTeamChange,
  onStageChange,
  isLoadingChampionships,
  isLoadingMatches,
}: FiltersSectionProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Filtros</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">Campeonato</label>
            <ChampionshipSelector
              championships={championships}
              value={selectedChampionship}
              onValueChange={onChampionshipChange}
              disabled={isLoadingChampionships}
            />
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium">Time</label>
            <TeamSelector
              matches={matches}
              value={selectedTeam}
              onValueChange={onTeamChange}
              disabled={isLoadingMatches || !selectedChampionship}
            />
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium">Fase</label>
            <StageSelector
              matches={matches}
              value={selectedStage}
              onValueChange={onStageChange}
              disabled={isLoadingMatches || !selectedChampionship}
            />
          </div>
        </div>
      </CardContent>
    </Card>
  )
}
