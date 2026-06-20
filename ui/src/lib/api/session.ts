import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import type { Session, SessionCreate, SessionUpdate } from "../types/session";
import { convertSessions } from "../types/session";
import { apiDelete, apiGet, apiPost, apiPut, NO_CONVERTER, NO_FILES } from "./query";

const ENDPOINT = "session";

export const useSessionGetAll = () => {
  return useQuery({
    queryKey: ["session"],
    queryFn: async () => (await apiGet(ENDPOINT, convertSessions, true)).data,
  });
};

export const useSessionCreate = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (session: SessionCreate) => await apiPut(ENDPOINT, session, NO_CONVERTER, NO_FILES, true),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["session"] })
    },
  });
};

export const useSessionUpdate = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (session: SessionUpdate) => apiPost(`${ENDPOINT}/${session.id}`, session, NO_CONVERTER, NO_FILES, true),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["session"] })
    },
  });
};

export const useSessionDelete = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id }: Pick<Session, "id">) => apiDelete(`${ENDPOINT}/${id}`, NO_CONVERTER, true),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["session"] })
    },
  });
};
