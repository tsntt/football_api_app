import { useQuery } from "@tanstack/react-query"
import { apiClient } from "@/lib/api"
import { queryKeys } from "@/lib/query-client"

export function useChampionships() {
  return useQuery({
    queryKey: queryKeys.championships,
    queryFn: () => apiClient.getChampionships(),
    staleTime: 10 * 60 * 1000, // 10 minutes - championships don't change often
    gcTime: 30 * 60 * 1000, // 30 minutes
  })
}
