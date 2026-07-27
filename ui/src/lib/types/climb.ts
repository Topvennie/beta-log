import type { API } from "./api";
import { convertGym } from "./gym";
import type { Gym } from "./gym";

export type ClimbType = "boulder" | "lead";
export type ClimbFinish = "flash" | "top" | "repeat";
export type ClimbSource = "toplogger" | "manual";

export interface Climb {
  id: number;
  grade: number;
  holdColor: string;
  climbType: ClimbType;
  finishType: ClimbFinish;
  source: ClimbSource;
}

export interface ClimbDay {
  id: number;
  date: Date;
  gym: Gym;
  climbs: Climb[];
  source: ClimbSource;
}

export interface ClimbStatsGraphProgress {
  date: string;
  grade: number;
  volume: number;
}

export interface ClimbStatsGraphGrade {
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
  medianClimbsPerSession: number;
  graphProgress: ClimbStatsGraphProgress[];
  graphPerGrade: ClimbStatsGraphGrade[];
}

// Converts

export const convertClimb = (c: API.Climb): Climb => ({
  id: c.id,
  grade: c.grade,
  holdColor: c.hold_color,
  climbType: c.climb_type as ClimbType,
  finishType: c.finish_type as ClimbFinish,
  source: c.source as ClimbSource,
});

export const convertClimbs = (c: API.Climb[]): Climb[] => c.map(convertClimb);

export const convertClimbDay = (d: API.ClimbDay): ClimbDay => ({
  id: d.id,
  date: new Date(d.date),
  gym: convertGym(d.gym),
  climbs: convertClimbs(d.climbs),
  source: d.source as ClimbSource,
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
  medianClimbsPerSession: s.median_climbs_per_session,
  graphProgress: s.graph_progress,
  graphPerGrade: s.graph_per_grade,
});
