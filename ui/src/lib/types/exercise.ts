import { API } from "./api";
import { z } from "zod";
import { JSONBody } from "./general";
import { convertVariants, convertVariantUpdateSchema, Variant, variantCreateSchema, variantUpdateSchema } from "./variant";

export interface Exercise {
  id: number;
  name: string;
  variants: Variant[];
};

// Converts

export const convertExercise = (e: API.Exercise): Exercise => ({
  id: e.id,
  name: e.name,
  variants: convertVariants(e.variants ?? []),
});

export const convertExercises = (e: API.Exercise[]): Exercise[] => e.map(convertExercise);

export const convertExerciseUpdateSchema = (e: Exercise): ExerciseUpdate => ({
  id: e.id,
  name: e.name,
  variants: e.variants.map(convertVariantUpdateSchema),
})

// Schemas

export const exerciseCreateSchema = z.object({
  name: z.string().min(1),
  variants: z.array(variantCreateSchema),
});
export type ExerciseCreate = z.infer<typeof exerciseCreateSchema> & JSONBody;

export const exerciseUpdateSchema = z.object({
  id: z.number().positive(),
  name: z.string().min(1),
  variants: z.array(variantUpdateSchema),
});
export type ExerciseUpdate = z.infer<typeof exerciseUpdateSchema> & JSONBody;

