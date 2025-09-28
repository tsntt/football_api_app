import { useQuery } from "@tanstack/react-query"
import { apiClient } from "@/lib/api"
import { queryKeys } from "@/lib/query-client"

export function useAdminMatches() {
  return useQuery({
    queryKey: queryKeys.admin.matches(),
    queryFn: () => apiClient.getAdminMatches(),
    staleTime: 5 * 60 * 1000, // 5 minute - admin data needs to be fresh
    gcTime: 10 * 60 * 1000, // 10 minutes
    refetchInterval: 10 * 60 *1000, // Refetch every 10 minutes updates
    refetchIntervalInBackground: false, // Continue refetching when tab is not active
  })
}
