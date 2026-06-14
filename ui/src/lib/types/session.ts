import { API } from "./api";
import { z } from "zod";
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
  name: string;
  variant: string;
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
  exercises: s.exercises.map(convertSessionExercise),
});

export const convertSessions = (s: API.Session[]): Session[] => s.map(convertSession);

const convertSessionExercise = (se: API.SessionExercise): SessionExercise => ({
  id: se.id,
  name: se.name,
  variant: se.variant,
  position: se.position,
  sets: se.sets,
  reps: se.reps,
  weight: se.weight,
  durationS: se.duration_s,
});

// Schemas

export const SessionExerciseCreateSchema = z.object({
  exercise_id: z.number().positive(),
  position: z.number().positive(),
  sets: z.number().positive(),
  reps: z.number().optional(),
  weight: z.number().optional(),
  duration_s: z.number().optional(),
});

export const SessionCreateSchema = z.object({
  name: z.string(),
  active: z.boolean(),
  position: z.number().positive().optional(),
  exercises: z.array(SessionExerciseCreateSchema).min(1),
}).refine(data => data.active || (data.position ?? 0) > 0, {
  message: "A position is required if session is active",
  path: ["position"],
});
export type SessionCreate = z.infer<typeof SessionCreateSchema> & JSONBody;

export const SessionExerciseUpdateSchema = z.object({
  id: z.number().positive(),
  position: z.number().positive(),
  sets: z.number().positive(),
  reps: z.number().optional(),
  weight: z.number().optional(),
  duration_s: z.number().optional(),
});

export const SessionUpdateSchema = z.object({
  id: z.number().positive(),
  name: z.string(),
  active: z.boolean(),
  position: z.number().positive().optional(),
  exercises: z.array(SessionExerciseUpdateSchema).min(1),
}).refine(data => data.active || (data.position ?? 0) > 0, {
  message: "A position is required if session is active",
  path: ["position"],
});
export type SessionUpdate = z.infer<typeof SessionUpdateSchema> & JSONBody;
