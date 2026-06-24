import { z } from "zod";
import { API } from "./api";
import { convertExercise, Exercise } from "./exercise";
import { JSONBody } from "./general";
import { convertVariant, Variant } from "./variant";

export interface Session {
  id: number;
  name: string;
  exercises: SessionExercise[];
};

export interface SessionExercise {
  id: number;
  exercise: Omit<Exercise, "variants">;
  variant?: Variant;
  position: number;
  sets: number;
  reps?: number;
  weight?: number;
  durationS?: number;
};

// Converts

export const convertSession = (s: API.Session): Session => ({
  id: s.id,
  name: s.name,
  exercises: convertSessionExercises(s.exercises),
});
export const convertSessions = (s: API.Session[]): Session[] => s.map(convertSession);

export const convertSessionExercise = (se: API.SessionExercise): SessionExercise => ({
  id: se.id,
  exercise: convertExercise(se.exercise),
  variant: se.variant ? convertVariant(se.variant) : undefined,
  position: se.position,
  sets: se.sets,
  reps: se.reps,
  weight: se.weight,
  durationS: se.duration_s,
});
export const convertSessionExercises = (se: API.SessionExercise[]): SessionExercise[] => se.map(convertSessionExercise)

export const convertSessionUpdateSchema = (s: Session): SessionUpdate => ({
  id: s.id,
  name: s.name,
  exercises: s.exercises.map(convertSessionExerciseUpdateSchema),
})

export const convertSessionExerciseUpdateSchema = (se: SessionExercise): SessionExerciseUpdate => ({
  clientId: crypto.randomUUID(),
  exerciseId: se.exercise.id,
  variantId: se.variant ? se.variant.id : undefined,
  position: se.position,
  sets: se.sets,
  reps: se.reps,
  weight: se.weight,
  durationS: se.durationS,
})

// Schemas

export const sessionExerciseCreateSchema = z.object({
  clientId: z.string(),
  exerciseId: z.number().positive(),
  variantId: z.number().optional(),
  position: z.number().positive(),
  sets: z.number().positive(),
  reps: z.number().positive().optional(),
  weight: z.number().positive().optional(),
  durationS: z.number().positive().optional(),
});
export type SessionExerciseCreate = z.infer<typeof sessionCreateSchema> & JSONBody;

export const sessionCreateSchema = z.object({
  name: z.string().min(1),
  exercises: z.array(sessionExerciseCreateSchema).min(1),
})
export type SessionCreate = z.infer<typeof sessionCreateSchema> & JSONBody;

export const sessionExerciseUpdateSchema = z.object({
  clientId: z.string(),
  exerciseId: z.number().positive(),
  variantId: z.number().optional(),
  position: z.number().positive(),
  sets: z.number().positive(),
  reps: z.number().positive().optional(),
  weight: z.number().positive().optional(),
  durationS: z.number().positive().optional(),
});
export type SessionExerciseUpdate = z.infer<typeof sessionExerciseUpdateSchema> & JSONBody;

export const sessionUpdateSchema = z.object({
  id: z.number().positive(),
  name: z.string().min(1),
  exercises: z.array(sessionExerciseUpdateSchema).min(1),
})
export type SessionUpdate = z.infer<typeof sessionUpdateSchema> & JSONBody;
