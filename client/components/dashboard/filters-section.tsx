"use client"

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { ChampionshipSelector } from "./championship-selector"
import { TeamSelector } from "./team-selector"
import { StageSelector } from "./stage-selector"
import { useMemo } from "react"
import type { Championship, Match, Team } from "@/lib/types"

import { Button } from "@/components/ui/button"
import { toast } from "sonner"
import { apiClient } from "@/lib/api"
import { useMutation } from "@tanstack/react-query"

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
  const { mutate: subscribe } = useMutation({
    mutationFn: ({ teamId, teamName }: { teamId: number; teamName: string }) =>
      apiClient.subscribeToTeam(teamId, teamName),
    onSuccess: () => {
      toast.success("Subscribed!", {
        description: `You are now subscribed to ${selectedTeam}`,
      })
    },
    onError: (error: any) => {
      toast.error("Error", {
        description: error.message || "Something went wrong",
      })
    },
  })

  const handleSubscribe = () => {
    if (!selectedTeam) {
      toast.error("No team selected", {
        description: "Please select a team to subscribe",
      })
      return
    }

    const team = matches
      .flatMap((match) => [match.homeTeam, match.awayTeam])
      .filter(Boolean)
      .find((team) => team.shortName === selectedTeam)

    if (!team) {
      toast.error("Team not found", {
        description: "Could not find the selected team",
      })
      return
    }

    subscribe({ teamId: team.id, teamName: team.name })
  }

  const filteredTeams = useMemo(() => {
    if (!matches) return []
    if (!selectedStage || selectedStage === "all") {
      const teamMap = new Map<number, Team>()
      matches.forEach((match) => {
        if (match?.homeTeam?.id && match?.awayTeam?.id) {
          teamMap.set(match.homeTeam.id, match.homeTeam)
          teamMap.set(match.awayTeam.id, match.awayTeam)
        }
      })
      return Array.from(teamMap.values()).sort((a, b) =>
        a.name.localeCompare(b.name),
      )
    }

    const teamMap = new Map<number, Team>()
    matches
      .filter((match) => match.stage === selectedStage)
      .forEach((match) => {
        if (match?.homeTeam?.id && match?.awayTeam?.id) {
          teamMap.set(match.homeTeam.id, match.homeTeam)
          teamMap.set(match.awayTeam.id, match.awayTeam)
        }
      })
    return Array.from(teamMap.values()).sort((a, b) =>
      a.name.localeCompare(b.name),
    )
  }, [matches, selectedStage])

  const filteredStages = useMemo(() => {
    if (!matches) return []
    if (!selectedTeam) {
      const stageSet = new Set<string>()
      matches.forEach((match) => {
        if (match.stage) {
          stageSet.add(match.stage)
        }
      })
      return Array.from(stageSet)
    }

    const stageSet = new Set<string>()
    matches
      .filter(
        (match) =>
          match.homeTeam?.shortName === selectedTeam ||
          match.awayTeam?.shortName === selectedTeam,
      )
      .forEach((match) => {
        if (match.stage) {
          stageSet.add(match.stage)
        }
      })
    return Array.from(stageSet)
  }, [matches, selectedTeam])

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
              teams={filteredTeams}
              value={selectedTeam}
              onValueChange={onTeamChange}
              disabled={isLoadingMatches || !selectedChampionship}
            />
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium">Fase</label>
            <StageSelector
              stages={filteredStages}
              value={selectedStage}
              onValueChange={onStageChange}
              disabled={isLoadingMatches || !selectedChampionship}
            />
          </div>
        </div>
        <div className="mt-4 w-full flex flex-row justify-end">
           {selectedTeam && <Button
            className="cursor-pointer"
            onClick={handleSubscribe}
            disabled={!selectedTeam || isLoadingMatches}
          >
            Subscribe to {selectedTeam}
          </Button>}
        </div>
      </CardContent>
    </Card>
  )
}
