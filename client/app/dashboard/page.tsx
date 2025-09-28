"use client"

import { useState, useEffect } from "react"
import { useRouter } from "next/navigation"
import { useAuthStore } from "@/lib/auth"
import { useChampionships } from "@/hooks/use-championships"
import { useMatches } from "@/hooks/use-matches"
import { DashboardHeader } from "@/components/dashboard/dashboard-header"
import { FiltersSection } from "@/components/dashboard/filters-section"
import { MatchesList } from "@/components/matches/matches-list"

export default function DashboardPage() {
  const { isAuthenticated } = useAuthStore()
  const router = useRouter()

  const [selectedChampionship, setSelectedChampionship] = useState<number | null>(null)
  const [selectedTeam, setSelectedTeam] = useState("") // Now stores team ID as string
  const [selectedStage, setSelectedStage] = useState("")

  const { data: championships = [], isLoading: isLoadingChampionships } = useChampionships()

  const { data: allMatches = [], isLoading: isLoadingAllMatches } = useMatches(
    selectedChampionship,
    undefined, // No team filter for selectors
    undefined, // No stage filter for selectors
  )

  const { data: filteredMatches = [], isLoading: isLoadingFilteredMatches } = useMatches(
    selectedChampionship,
    selectedTeam || undefined, 
    selectedStage === "all" ? undefined : selectedStage || undefined,
  )

  useEffect(() => {
    if (!isAuthenticated) {
      router.push("/")
    }
  }, [isAuthenticated, router])

  if (!isAuthenticated) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <h2 className="text-xl font-semibold mb-2">Carregando...</h2>
          <p className="text-muted-foreground">Modo de teste ativado</p>
        </div>
      </div>
    )
  }

  const handleChampionshipChange = (championshipId: number | null) => {
    setSelectedChampionship(championshipId)
    setSelectedTeam("")
    setSelectedStage("")
  }

  return (
    <div className="min-h-screen bg-background">
      <DashboardHeader />

      <main className="container mx-auto px-4 py-6 space-y-6">
        <FiltersSection
          championships={championships}
          matches={allMatches} // Pass unfiltered matches to selectors
          selectedChampionship={selectedChampionship}
          selectedTeam={selectedTeam}
          selectedStage={selectedStage}
          onChampionshipChange={handleChampionshipChange}
          onTeamChange={setSelectedTeam}
          onStageChange={setSelectedStage}
          isLoadingChampionships={isLoadingChampionships}
          isLoadingMatches={isLoadingAllMatches}
        />

        <div className="space-y-4">
          {selectedChampionship ? (
            <MatchesList matches={filteredMatches} isLoading={isLoadingFilteredMatches} /> // Use filtered matches for display
          ) : (
            <div className="text-center py-12">
              <div className="text-6xl mb-4">🏆</div>
              <h3 className="text-lg font-semibold mb-2">Selecione um campeonato</h3>
              <p className="text-muted-foreground">
                Escolha um campeonato acima para visualizar as partidas disponíveis.
              </p>
            </div>
          )}
        </div>
      </main>
    </div>
  )
}
