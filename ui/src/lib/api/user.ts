import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { convertUserToModel } from "../types/user";
import { apiGet, apiPost } from "./query";

const ENDPOINT_AUTH = "auth"
const ENDPOINT_USER = "user"

export const useUser = () => {
  return useQuery({
    queryKey: ["user"],
    queryFn: async () => (await apiGet(`${ENDPOINT_USER}/me`, convertUserToModel, true)).data,
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
    mutationFn: async () => (await apiPost(`${ENDPOINT_AUTH}/logout`)).data,
    onSuccess: async () => {
      await queryClient.cancelQueries({ queryKey: ["user"], exact: true })
      queryClient.setQueryData(["user"], null)
    },
  })
}

