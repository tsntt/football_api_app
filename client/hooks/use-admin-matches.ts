import { useQuery } from "@tanstack/react-query"
import { apiClient } from "@/lib/api"
import { queryKeys } from "@/lib/query-client"

export function useAdminMatches() {
  return useQuery({
    queryKey: queryKeys.admin.matches(),
    queryFn: () => apiClient.getAdminMatches(),
    staleTime: 1 * 60 * 1000, // 1 minute - admin data needs to be fresh
    gcTime: 2 * 60 * 1000, // 2 minutes
    refetchInterval: 30 * 1000, // Refetch every 30 seconds for real-time updates
    refetchIntervalInBackground: true, // Continue refetching when tab is not active
  })
}
