import { useQuery } from "@tanstack/react-query"
import { apiClient } from "@/lib/api"
import { queryKeys } from "@/lib/query-client"

export function useMatches(championshipId: number | null, teamId?: string, stage?: string) {
  return useQuery({
    queryKey: queryKeys.matches(championshipId!, teamId, stage),
    queryFn: () => apiClient.getMatches(championshipId!, teamId, stage),
    enabled: !!championshipId,
    staleTime: 2 * 60 * 1000, // 2 minutes - match data changes more frequently
    gcTime: 5 * 60 * 1000, // 5 minutes
  })
}