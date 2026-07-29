import z from "zod";
import type { API } from "./api";
import { JSONBody } from "./general";
import type { Gym, Source } from "./gym";
import { convertGym } from "./gym";

export enum ClimbType {
  Boulder = "boulder",
  Lead = "lead",
}

export enum ClimbFinish {
  Flash = "flash",
  Top = "top",
  Repeat = "repeat",
}

export interface Climb {
  id: number;
  grade: number;
  holdColor: string;
  climbType: ClimbType;
  finishType: ClimbFinish;
  source: Source;
}

export interface ClimbDay {
  id: number;
  date: Date;
  gym: Gym;
  climbs: Climb[];
  source: Source;
}

export interface ClimbStatsProgress {
  date: string;
  grade: number;
  volume: number;
}

export interface ClimbStatsGrade {
  grade: number;
  flash: number;
  top: number;
  repeat: number;
}

export interface ClimbStats {
  total: number;
  totalUnique: number;
  flash: number;
  top: number;
  repeat: number;
  best: number;
  bestAmount: number;
  bestFlash: number;
  bestFlashAmount: number;
  sessions: number;
  boulder: number;
  lead: number;
  medianClimbsPerSession: number;
  graphProgress: ClimbStatsProgress[];
  graphPerGrade: ClimbStatsGrade[];
}

// Converts

export const convertClimb = (c: API.Climb): Climb => ({
  id: c.id,
  grade: c.grade,
  holdColor: c.hold_color,
  climbType: c.climb_type as ClimbType,
  finishType: c.finish_type as ClimbFinish,
  source: c.source as Source,
});

export const convertClimbs = (c: API.Climb[]): Climb[] => c.map(convertClimb);

export const convertClimbDay = (d: API.ClimbDay): ClimbDay => ({
  id: d.id,
  date: new Date(d.date),
  gym: convertGym(d.gym),
  climbs: convertClimbs(d.climbs),
  source: d.source as Source,
});

export const convertClimbDays = (d: API.ClimbDay[]): ClimbDay[] => d.map(convertClimbDay);

export const convertClimbStats = (s: API.ClimbStats): ClimbStats => ({
  total: s.total,
  totalUnique: s.total_unique,
  flash: s.flash,
  top: s.top,
  repeat: s.repeat,
  best: s.best,
  bestAmount: s.best_amount,
  bestFlash: s.best_flash,
  bestFlashAmount: s.best_flash_amount,
  sessions: s.sessions,
  boulder: s.boulder,
  lead: s.lead,
  medianClimbsPerSession: s.median_climbs_per_session,
  graphProgress: s.graph_progress,
  graphPerGrade: s.graph_per_grade,
});

export const convertClimbCreateSchema = (c: Climb): ClimbCreate => ({
  _clientId: crypto.randomUUID(),
  grade: c.grade,
  holdColor: c.holdColor,
  climbType: c.climbType,
  finishType: c.finishType,
})

export const convertClimbDayUpdateSchema = (c: ClimbDay): ClimbDayUpdate => ({
  id: c.id,
  date: c.date,
  gymId: c.gym.id,
  climbs: c.climbs.map(convertClimbCreateSchema),
})

// Schemas

export const climbCreateSchema = z.object({
  _clientId: z.string(),
  grade: z.number().positive(),
  holdColor: z.string(),
  climbType: z.enum(ClimbType),
  finishType: z.enum(ClimbFinish),
})
export type ClimbCreate = z.infer<typeof climbCreateSchema> & JSONBody

export const climbDayCreateSchema = z.object({
  date: z.date(),
  gymId: z.number().positive(),
  climbs: z.array(climbCreateSchema).min(1),
})
export type ClimbDayCreate = z.infer<typeof climbDayCreateSchema> & JSONBody

export const climbDayUpdateSchema = z.object({
  id: z.number().positive(),
  date: z.date(),
  gymId: z.number().positive(),
  climbs: z.array(climbCreateSchema).min(1),
})
export type ClimbDayUpdate = z.infer<typeof climbDayUpdateSchema> & JSONBody
