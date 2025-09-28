"use client"

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Progress } from "@/components/ui/progress"
import { Badge } from "@/components/ui/badge"
import { CheckCircle, XCircle, Clock, Users } from "lucide-react"
import type { WebSocketUpdate } from "@/lib/types"

interface BroadcastProgressProps {
  updates: WebSocketUpdate[]
}

export function BroadcastProgress({ updates }: BroadcastProgressProps) {
  if (updates.length === 0) return null

  return (
    <div className="space-y-4">
      <h3 className="text-lg font-semibold">Status das Notificações</h3>
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {updates.map((update, index) => {
          const progress = update.total_sent > 0 ? (update.sent_count / update.total_sent) * 100 : 0
          const hasErrors = update.failed_count > 0

          return (
            <Card key={`${update.channel_id}-${index}`} className="relative">
              <CardHeader className="pb-3">
                <div className="flex items-center justify-between">
                  <CardTitle className="text-sm">Canal #{update.channel_id}</CardTitle>
                  <div className="flex items-center gap-1">
                    {update.is_completed ? (
                      hasErrors ? (
                        <Badge variant="destructive" className="gap-1">
                          <XCircle className="h-3 w-3" />
                          Concluído com erros
                        </Badge>
                      ) : (
                        <Badge variant="default" className="gap-1">
                          <CheckCircle className="h-3 w-3" />
                          Concluído
                        </Badge>
                      )
                    ) : (
                      <Badge variant="secondary" className="gap-1">
                        <Clock className="h-3 w-3" />
                        Em andamento
                      </Badge>
                    )}
                  </div>
                </div>
              </CardHeader>
              <CardContent className="space-y-3">
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <Users className="h-4 w-4" />
                  {update.total_sent} usuários
                </div>

                <div className="space-y-2">
                  <div className="flex justify-between text-sm">
                    <span>Progresso</span>
                    <span>{Math.round(progress)}%</span>
                  </div>
                  <Progress value={progress} className="h-2" />
                </div>

                <div className="grid grid-cols-2 gap-4 text-sm">
                  <div className="text-center">
                    <div className="font-semibold text-green-600">{update.sent_count}</div>
                    <div className="text-muted-foreground">Enviadas</div>
                  </div>
                  <div className="text-center">
                    <div className="font-semibold text-red-600">{update.failed_count}</div>
                    <div className="text-muted-foreground">Falharam</div>
                  </div>
                </div>

                {update.error_details.length > 0 && (
                  <div className="mt-2 p-2 bg-destructive/10 rounded text-xs">
                    <div className="font-medium text-destructive mb-1">Erros:</div>
                    {update.error_details.slice(0, 3).map((error, i) => (
                      <div key={i} className="text-destructive/80">
                        {error}
                      </div>
                    ))}
                    {update.error_details.length > 3 && (
                      <div className="text-destructive/60">+{update.error_details.length - 3} mais erros...</div>
                    )}
                  </div>
                )}
              </CardContent>
            </Card>
          )
        })}
      </div>
    </div>
  )
}
