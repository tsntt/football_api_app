"use client"

import { useState, useMemo } from "react"
import { Check, ChevronsUpDown } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from "@/components/ui/command"
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover"
import { cn } from "@/lib/utils"
import type { Team } from "@/lib/types"

interface TeamSelectorProps {
  teams: Team[]
  value: string
  onValueChange: (value: string) => void
  disabled?: boolean
}

export function TeamSelector({ teams, value, onValueChange, disabled = false }: TeamSelectorProps) {
  const [open, setOpen] = useState(false)

  const selectedTeam = teams.find((t) => t.shortName === value)

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          role="combobox"
          aria-expanded={open}
          className="w-full justify-between bg-transparent"
          disabled={disabled}
        >
          {selectedTeam ? (
            <div className="flex items-center gap-2">
              <img src={selectedTeam.crest || "/placeholder.svg"} alt={selectedTeam.name} className="w-4 h-4" />
              {selectedTeam.name}
            </div>
          ) : (
            "Selecione um time..."
          )}
          <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-full p-0">
        <Command>
          <CommandInput placeholder="Buscar time..." />
          <CommandList>
            <CommandEmpty>Nenhum time encontrado.</CommandEmpty>
            <CommandGroup>
              <CommandItem
                value=""
                onSelect={() => {
                  onValueChange("")
                  setOpen(false)
                }}
              >
                <Check className={cn("mr-2 h-4 w-4", value === "" ? "opacity-100" : "opacity-0")} />
                Todos os times
              </CommandItem>
              {teams.map((team) => (
                <CommandItem
                  key={team.shortName}
                  value={team.name}
                  onSelect={() => {
                    onValueChange(team.shortName === value ? "" : team.shortName)
                    setOpen(false)
                  }}
                >
                  <Check className={cn("mr-2 h-4 w-4", value === team.shortName ? "opacity-100" : "opacity-0")} />
                  <div className="flex items-center gap-2">
                    <img src={team.crest || "/placeholder.svg"} alt={team.name} className="w-4 h-4" />
                    {team.name}
                  </div>
                </CommandItem>
              ))}
            </CommandGroup>
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  )
}
