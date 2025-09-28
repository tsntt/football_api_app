import { QueryClient } from "@tanstack/react-query"

export function createQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: {
        // Stale time - how long data is considered fresh
        staleTime: 5 * 60 * 1000, // 5 minutes
        // Garbage collection time - how long unused data stays in cache
        gcTime: 10 * 60 * 1000, // 10 minutes (formerly cacheTime)
        // Retry configuration
        retry: (failureCount, error) => {
          // Don't retry on 4xx errors (client errors)
          if (error instanceof Error && error.message.includes("4")) {
            return false
          }
          // Retry up to 3 times for other errors
          return failureCount < 3
        },
        retryDelay: (attemptIndex) => Math.min(1000 * 2 ** attemptIndex, 30000),
        // Refetch configuration
        refetchOnWindowFocus: true,
        refetchOnReconnect: true,
        refetchOnMount: true,
      },
      mutations: {
        // Retry mutations once on network errors
        retry: (failureCount, error) => {
          if (error instanceof Error && error.message.includes("NetworkError")) {
            return failureCount < 1
          }
          return false
        },
      },
    },
  })
}

// Query keys factory for consistent key management
export const queryKeys = {
  // Authentication
  auth: ["auth"] as const,

  // Championships
  championships: ["championships"] as const,

  // Matches
  matches: (championshipId: number, team?: string, stage?: string) =>
    (["matches", championshipId, team, stage].filter(Boolean) as (string | number)[]),

  // Admin
  admin: {
    all: ["admin"] as const,
    matches: () => [...queryKeys.admin.all, "matches"] as const,
  },
} as const
