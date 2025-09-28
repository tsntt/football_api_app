"use client"

import { useMemo } from "react"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import type { Match } from "@/lib/types"

interface StageSelectorProps {
  matches: Match[]
  value: string
  onValueChange: (value: string) => void
  disabled?: boolean
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

export function StageSelector({ matches, value, onValueChange, disabled = false }: StageSelectorProps) {
  const stages = useMemo(() => {
    if (!matches || !Array.isArray(matches)) {
      return []
    }

    const stageSet = new Set<string>()
    matches.forEach((match) => {
      stageSet.add(match.stage)
    })
    return Array.from(stageSet).sort()
  }, [matches])

  return (
    <Select value={value} onValueChange={onValueChange} disabled={disabled}>
      <SelectTrigger className="w-full">
        <SelectValue placeholder="Selecione uma fase..." />
      </SelectTrigger>
      <SelectContent>
        <SelectItem value="all">Todas as fases</SelectItem>
        {stages.map((stage) => (
          <SelectItem key={stage} value={stage}>
            {stageLabels[stage] || stage}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  )
}
