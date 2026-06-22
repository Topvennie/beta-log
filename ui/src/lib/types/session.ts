import { z } from "zod";
import { API } from "./api";
import { convertExercise, Exercise } from "./exercise";
import { JSONBody } from "./general";

export interface Session {
  id: number;
  name: string;
  active: boolean;
  position?: number;
  exercises: SessionExercise[];
};

export interface SessionExercise {
  id: number;
  exercise: Exercise;
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
  active: s.active,
  position: s.position,
  exercises: convertSessionExercises(s.exercises),
});
export const convertSessions = (s: API.Session[]): Session[] => s.map(convertSession);

export const convertSessionExercise = (se: API.SessionExercise): SessionExercise => ({
  id: se.id,
  exercise: convertExercise(se.exercise),
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
  active: s.active,
  position: s.position,
  exercises: s.exercises.map((se) => ({
    ...convertSessionExerciseUpdateSchema(se),
    clientId: crypto.randomUUID(),
  })),
})

export const convertSessionExerciseUpdateSchema = (se: SessionExercise) => ({
  exerciseId: se.exercise.id,
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
  variant: z.string().optional(),
  position: z.number().positive(),
  sets: z.number().positive(),
  reps: z.number().optional(),
  weight: z.number().optional(),
  durationS: z.number().optional(),
});
export type SessionExerciseCreate = z.infer<typeof sessionCreateSchema> & JSONBody;

export const sessionCreateSchema = z.object({
  name: z.string().min(1),
  active: z.boolean(),
  position: z.number().positive().optional(),
  exercises: z.array(sessionExerciseCreateSchema).min(1),
})
export type SessionCreate = z.infer<typeof sessionCreateSchema> & JSONBody;

export const sessionExerciseUpdateSchema = z.object({
  clientId: z.string(),
  exerciseId: z.number().positive(),
  variant: z.string().optional(),
  position: z.number().positive(),
  sets: z.number().positive(),
  reps: z.number().optional(),
  weight: z.number().optional(),
  durationS: z.number().optional(),
});
export type SessionExerciseUpdateSchema = z.infer<typeof sessionExerciseUpdateSchema> & JSONBody;

export const sessionUpdateSchema = z.object({
  id: z.number().positive(),
  name: z.string().min(1),
  active: z.boolean(),
  position: z.number().positive().optional(),
  exercises: z.array(sessionExerciseUpdateSchema).min(1),
})
export type SessionUpdate = z.infer<typeof sessionUpdateSchema> & JSONBody;
