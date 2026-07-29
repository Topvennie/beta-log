import z from "zod";
import type { API } from "./api";
import { JSONBody } from "./general";

export enum Source {
  Toplogger = "toplogger",
  Manual = "manual",
}


export interface Gym {
  id: number;
  name: string;
  iconPath: string;
  source: Source;
}

// Converts

export const convertGym = (g: API.Gym): Gym => ({
  id: g.id,
  name: g.name,
  iconPath: g.icon_path,
  source: g.source as Source,
});

export const convertGyms = (g: API.Gym[]): Gym[] => g.map(convertGym);

export const convertGymUpdateSchema = (g: Gym): GymUpdate => {
  return {
    id: g.id,
    name: g.name,
    iconPath: g.iconPath,
  }
}

export interface GymStatsVisits {
  gym: string;
  amount: number;
}

export interface GymStatsTop {
  gym: string;
  top: number;
  flash: number;
}

export interface GymStatsDistribution {
  gym: string;
  distribution: Record<number, number>;
}

export interface GymStats {
  total: number;
  mostVisited: string;
  mostVisitedAmount: number;
  sessions: number;
  graphVisits: GymStatsVisits[];
  graphTop: GymStatsTop[];
  graphDistribution: GymStatsDistribution[];
}

export const convertGymStats = (s: API.GymStats): GymStats => ({
  total: s.total,
  mostVisited: s.most_visited,
  mostVisitedAmount: s.most_visited_amount,
  sessions: s.sessions,
  graphVisits: s.graph_visits,
  graphTop: s.graph_top,
  graphDistribution: s.graph_distribution,
});

// Schemas

export const gymCreateSchema = z.object({
  name: z.string().min(1),
  iconPath: z.string().optional(),
})
export type GymCreate = z.infer<typeof gymCreateSchema> & JSONBody

export const gymUpdateSchema = z.object({
  id: z.number().positive(),
  name: z.string().min(1),
  iconPath: z.string().optional(),
})
export type GymUpdate = z.infer<typeof gymUpdateSchema> & JSONBody
