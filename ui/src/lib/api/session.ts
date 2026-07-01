import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import type { Session, SessionCreate, SessionUpdate } from "../types/session";
import { convertSession, convertSessions } from "../types/session";
import { apiDelete, apiGet, apiPost, apiPut, NO_CONVERTER } from "./query";

const ENDPOINT = "session";

export const useSessionGetAll = () => {
  return useQuery({
    queryKey: ["session"],
    queryFn: async () => (await apiGet(ENDPOINT, convertSessions)).data,
  });
};

export const useSessionCreate = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (session: SessionCreate) => await apiPost(ENDPOINT, session, convertSession),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["session"] })
    },
  });
};

export const useSessionUpdate = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (session: SessionUpdate) => apiPut(`${ENDPOINT}/${session.id}`, session, convertSession),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["session"] })
    },
  });
};

export const useSessionDelete = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id }: Pick<Session, "id">) => apiDelete(`${ENDPOINT}/${id}`, NO_CONVERTER),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["session"] })
    },
  });
};
