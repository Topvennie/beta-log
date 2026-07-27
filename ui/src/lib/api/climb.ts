import { useInfiniteQuery, useQuery } from "@tanstack/react-query";
import { format } from "date-fns";
import { convertClimbDays, convertClimbGyms, convertClimbStats } from "../types/climb";
import { apiGet } from "./query";

const ENDPOINT = "climb";
const PAGE_LIMIT = 10;
const DATE_FORMAT = "dd-MM-yyyy";

export const useClimbDayGetFiltered = () => {
  const { data, isLoading, fetchNextPage, isFetchingNextPage, hasNextPage, error, refetch, isFetching } = useInfiniteQuery({
    queryKey: ["climb", "day"],
    queryFn: async ({ pageParam = 1 }) => {
      const queryParams = new URLSearchParams({
        page: pageParam.toString(),
        limit: PAGE_LIMIT.toString(),
      });

      const url = `${ENDPOINT}/day?${queryParams.toString()}`;
      return (await apiGet(url, convertClimbDays)).data;
    },
    initialPageParam: 1,
    getNextPageParam: (lastPage, allPages) => {
      return lastPage.length < PAGE_LIMIT ? undefined : allPages.length + 1;
    },
    throwOnError: true,
  });

  const days = data?.pages.flat() ?? [];

  return {
    days,
    isLoading,
    fetchNextPage,
    isFetchingNextPage,
    hasNextPage,
    error,
    refetch,
    isFetching,
  };
}

export const useClimbStatGetFiltered = (start?: Date, end?: Date) => {
  return useQuery({
    queryKey: ["climb", "stat", start, end],
    queryFn: async () => {
      const queryParams = new URLSearchParams();

      if (start) {
        queryParams.append("start", format(start, DATE_FORMAT));
      }
      if (end) {
        queryParams.append("end", format(end, DATE_FORMAT));
      }

      const url = `${ENDPOINT}/stat${queryParams.toString() ? `?${queryParams.toString()}` : ""}`;
      return (await apiGet(url, convertClimbStats)).data;
    },
    throwOnError: true,
  });
}

export const useClimbGymGetAll = () => {
  return useQuery({
    queryKey: ["climb", "gym"],
    queryFn: async () => (await apiGet(`${ENDPOINT}/gym`, convertClimbGyms)).data,
  })
}
