import { API } from "./api";
import { z } from "zod";
import { JSONBody } from "./general";

export interface Exercise {
  id: number;
  name: string;
  variant: string;
};

// Converts

export const convertExercise = (e: API.Exercise): Exercise => ({
  id: e.id,
  name: e.name,
  variant: e.variant,
});

export const convertExercises = (e: API.Exercise[]): Exercise[] => e.map(convertExercise);

// Schemas

export const ExerciseCreateSchema = z.object({
  name: z.string(),
  variant: z.string().optional(),
});
export type ExerciseCreate = z.infer<typeof ExerciseCreateSchema> & JSONBody;

export const ExerciseUpdateSchema = z.object({
  id: z.number().positive(),
  name: z.string(),
  variant: z.string().optional(),
});
export type ExerciseUpdate = z.infer<typeof ExerciseUpdateSchema> & JSONBody;
