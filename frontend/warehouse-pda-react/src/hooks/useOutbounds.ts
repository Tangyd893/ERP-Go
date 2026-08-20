import { useQuery } from "@tanstack/react-query";
import { fetchOutbounds, fetchPickTasks } from "@/lib/warehouse-api";
import type { OutboundOrder, PickTask } from "@/stores/warehouse";

export const outboundQueryKey = ["warehouse", "outbounds"] as const;

export function useOutbounds() {
  return useQuery<OutboundOrder[]>({
    queryKey: outboundQueryKey,
    queryFn: fetchOutbounds,
  });
}

export function usePickTasks(outboundId: string) {
  return useQuery<PickTask[]>({
    queryKey: ["warehouse", "pick-tasks", outboundId] as const,
    queryFn: () => fetchPickTasks(outboundId),
    enabled: !!outboundId,
  });
}
