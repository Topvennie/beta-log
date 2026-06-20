import { API } from "./api";
import { z } from "zod";
import { JSONBody } from "./general";

export interface Exercise {
  id: number;
  name: string;
  variants: string[];
};

// Converts

export const convertExercise = (e: API.Exercise): Exercise => ({
  id: e.id,
  name: e.name,
  variants: e.variants,
});

export const convertExercises = (e: API.Exercise[]): Exercise[] => e.map(convertExercise);

export const convertExerciseUpdateSchema = (e: Exercise): ExerciseUpdate => ({
  id: e.id,
  name: e.name,
  variants: e.variants,
})

// Schemas

export const exerciseCreateSchema = z.object({
  name: z.string(),
  variants: z.array(z.string()),
});
export type ExerciseCreate = z.infer<typeof exerciseCreateSchema> & JSONBody;

export const exerciseUpdateSchema = z.object({
  id: z.number().positive(),
  name: z.string(),
  variants: z.array(z.string()),
});
export type ExerciseUpdate = z.infer<typeof exerciseUpdateSchema> & JSONBody;

