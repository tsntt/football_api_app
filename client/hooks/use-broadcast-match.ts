import { useMutation, useQueryClient } from "@tanstack/react-query"
import { apiClient } from "@/lib/api"
import { queryKeys } from "@/lib/query-client"
import { toast } from "sonner"

export function useBroadcastMatch() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (matchId: number) => apiClient.broadcastMatch(matchId),
    onSuccess: (data) => {
      toast.success(data.message)
      queryClient.invalidateQueries({ queryKey: queryKeys.admin.matches() })
    },
    onError: (error) => {
      toast.error("Erro ao enviar notificação")
      console.error("Broadcast error:", error)
    },
    // Optimistic updates could be added here if needed
    onMutate: async (matchId) => {
      // Cancel any outgoing refetches
      await queryClient.cancelQueries({ queryKey: queryKeys.admin.matches() })

      // Optionally show immediate feedback
      toast.loading("Enviando notificação...", { id: `broadcast-${matchId}` })
    },
    onSettled: (data, error, matchId) => {
      // Dismiss loading toast
      toast.dismiss(`broadcast-${matchId}`)
    },
  })
}
