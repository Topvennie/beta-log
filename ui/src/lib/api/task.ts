import { useInfiniteQuery, useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Task, TaskHistoryFilter, TaskStatus } from "../types/task";
import { convertTaskHistories, convertTasks } from "../types/task";
import { apiGet, apiPost } from "./query";

const ENDPOINT = "task";
const PAGE_LIMIT = 100;
const REFETCH_RUNNING_MS = 5 * 1000;
const REFETCH_IDLE_MS = 60 * 1000;

export function useTaskGetAll() {
  const queryClient = useQueryClient();

  return useQuery({
    queryKey: ["task"],
    queryFn: async () => (await apiGet(ENDPOINT, convertTasks)).data,
    refetchInterval: (query) => {
      const data = query.state.data;
      return data?.some((task: Task) => task.status === TaskStatus.Running)
        ? REFETCH_RUNNING_MS
        : REFETCH_IDLE_MS;
    },
    structuralSharing(oldData, newData) {
      if (oldData && JSON.stringify(oldData) === JSON.stringify(newData)) {
        return oldData;
      }

      queryClient.invalidateQueries({ queryKey: ["task_history"] });

      return newData;
    },
    throwOnError: true,
  });
}

export function useTaskGetHistory(filter?: TaskHistoryFilter) {
  const { data, isLoading, fetchNextPage, isFetchingNextPage, hasNextPage, error, refetch, isFetching } = useInfiniteQuery({
    queryKey: ["task_history", JSON.stringify(filter)],
    queryFn: async ({ pageParam = 1 }) => {
      const queryParams = new URLSearchParams({
        page: pageParam.toString(),
        limit: PAGE_LIMIT.toString(),
      });

      if (filter?.uid !== undefined) {
        queryParams.append("uid", filter.uid);
      }
      if (filter?.result !== undefined) {
        queryParams.append("result", filter.result.toString())
      }
      if (filter?.recurring !== undefined) {
        queryParams.append("recurring", String(filter.recurring))
      }

      const url = `${ENDPOINT}/history?${queryParams.toString()}`;
      return (await apiGet(url, convertTaskHistories)).data;
    },
    initialPageParam: 1,
    getNextPageParam: (lastPage, allPages) => {
      return lastPage.length < PAGE_LIMIT ? undefined : allPages.length + 1;
    },
    enabled: filter !== undefined,
    throwOnError: true,
  });

  const history = data?.pages.flat() ?? [];

  return {
    history,
    isLoading,
    fetchNextPage,
    isFetchingNextPage,
    hasNextPage,
    error,
    refetch,
    isFetching,
  };
}

export function useTaskStart() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ uid }: Pick<Task, "uid">) => apiPost(`${ENDPOINT}/start/${uid}`),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["task"] });
      void queryClient.invalidateQueries({ queryKey: ["task_history"] });
    },
  });
}
