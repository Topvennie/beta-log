import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { convertUserToModel } from "../types/user";
import { apiGet, apiPost, NO_CONVERTER, NO_DATA, NO_FILES } from "./query";

const ENDPOINT_AUTH = "auth"
const ENDPOINT_USER = "user"

export const useUser = () => {
  return useQuery({
    queryKey: ["user"],
    queryFn: async () => (await apiGet(`${ENDPOINT_USER}/me`, convertUserToModel)).data,
    staleTime: Infinity,
    throwOnError: false,
  })
}

export const useUserLogin = () => {
  window.location.href = `/api/${ENDPOINT_AUTH}/login/openid-connect`
}

export const useUserLogout = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => (await apiPost(`${ENDPOINT_AUTH}/logout`, NO_DATA, NO_CONVERTER, NO_FILES, false)).data,
    onSuccess: async () => {
      await queryClient.cancelQueries({ queryKey: ["user"], exact: true })
      queryClient.setQueryData(["user"], null)
    },
  })
}

