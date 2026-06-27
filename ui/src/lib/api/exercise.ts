import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import type { Exercise, ExerciseCreate, ExerciseUpdate } from "../types/exercise";
import { convertExercise, convertExercises } from "../types/exercise";
import { apiDelete, apiGet, apiPost, apiPut, NO_CONVERTER, NO_FILES } from "./query";

const ENDPOINT = "exercise";

export const useExerciseGetAll = () => {
  return useQuery({
    queryKey: ["exercise"],
    queryFn: async () => (await apiGet(ENDPOINT, convertExercises, true)).data,
  });
};

export const useExerciseCreate = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (exercise: ExerciseCreate) => apiPost(ENDPOINT, exercise, convertExercise, NO_FILES, true),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["exercise"] })
      queryClient.invalidateQueries({ queryKey: ["session"] })
    },
  });
};

export const useExerciseUpdate = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (exercise: ExerciseUpdate) => apiPut(`${ENDPOINT}/${exercise.id}`, exercise, convertExercise, NO_FILES, true),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["exercise"] })
      queryClient.invalidateQueries({ queryKey: ["session"] })
    },
  });
};

export const useExerciseDelete = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id }: Pick<Exercise, "id">) => apiDelete(`${ENDPOINT}/${id}`, NO_CONVERTER, true),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["exercise"] })
      queryClient.invalidateQueries({ queryKey: ["session"] })
    },
  });
};
