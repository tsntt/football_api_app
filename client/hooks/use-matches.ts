import { useQuery } from "@tanstack/react-query"
import { apiClient } from "@/lib/api"
import { queryKeys } from "@/lib/query-client"

export function useMatches(championshipId: number | null, team?: string, stage?: string) {
  return useQuery({
    queryKey: queryKeys.matches(championshipId!, team, stage),
    queryFn: () => apiClient.getMatches(championshipId!, team, stage),
    enabled: !!championshipId,
    staleTime: 5 * 60 * 1000, // 5 minutes - match data changes more frequently
    gcTime: 10 * 60 * 1000, // 10 minutes
  })
}