import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Gym, GymCreate, GymUpdate, convertGym, convertGymStats, convertGyms } from "../types/gym";
import { apiDelete, apiGet, apiPost, apiPut, NO_CONVERTER } from "./query";

const ENDPOINT = "gym";

export const useGymGetAll = () => {
  return useQuery({
    queryKey: ["gym"],
    queryFn: async () => (await apiGet(`${ENDPOINT}`, convertGyms)).data,
  })
}

export const useGymGetStats = () => {
  return useQuery({
    queryKey: ["gym", "stat"],
    queryFn: async () => (await apiGet(`${ENDPOINT}/stat`, convertGymStats)).data,
    throwOnError: true,
  })
}

export const useGymCreate = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (gym: GymCreate) => apiPost(`${ENDPOINT}`, gym, convertGym),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["gym"] })
  })
}

export const useGymUpdate = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (gym: GymUpdate) => apiPut(`${ENDPOINT}/${gym.id}`, gym, convertGym),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["gym"] })
  })
}

export const useGymDelete = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id }: Pick<Gym, "id">) => apiDelete(`${ENDPOINT}/${id}`, NO_CONVERTER),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["gym"] })
  })
}