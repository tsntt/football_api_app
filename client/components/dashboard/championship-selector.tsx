"use client"

import { useState } from "react"
import { Check, ChevronsUpDown } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from "@/components/ui/command"
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover"
import { cn } from "@/lib/utils"
import type { Championship } from "@/lib/types"

interface ChampionshipSelectorProps {
  championships: Championship[]
  value: number | null
  onValueChange: (value: number | null) => void
  disabled?: boolean
}

export function ChampionshipSelector({
  championships,
  value,
  onValueChange,
  disabled = false,
}: ChampionshipSelectorProps) {
  const [open, setOpen] = useState(false)

  const selectedChampionship = championships.find((c) => c.id === value)

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
          {selectedChampionship ? (
            <div className="flex items-center gap-2">
              <img
                src={selectedChampionship.emblem || "/placeholder.svg"}
                alt={selectedChampionship.name}
                className="w-4 h-4"
              />
              {selectedChampionship.name}
            </div>
          ) : (
            "Selecione um campeonato..."
          )}
          <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-full p-0">
        <Command>
          <CommandInput placeholder="Buscar campeonato..." />
          <CommandList>
            <CommandEmpty>Nenhum campeonato encontrado.</CommandEmpty>
            <CommandGroup>
              {championships.map((championship) => (
                <CommandItem
                  key={championship.id}
                  value={championship.name}
                  onSelect={() => {
                    onValueChange(championship.id === value ? null : championship.id)
                    setOpen(false)
                  }}
                >
                  <Check className={cn("mr-2 h-4 w-4", value === championship.id ? "opacity-100" : "opacity-0")} />
                  <div className="flex items-center gap-2">
                    <img src={championship.emblem || "/placeholder.svg"} alt={championship.name} className="w-4 h-4" />
                    {championship.name}
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
