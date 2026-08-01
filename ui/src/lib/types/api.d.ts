export namespace API {
  export interface User {
    id: number;
    uid: string;
    name: string;
  }

  export interface Exercise {
    id: number;
    name: string;
    variants?: Variant[];
  }

  export interface Variant {
    id: number;
    variant: string;
  }

  export interface Session {
    id: number;
    name: string;
    exercises: SessionExercise[];
  }

  export interface SessionExercise {
    id: number;
    exercise: Omit<Exercise, "variants">;
    variant?: Variant;
    position: number;
    sets: number;
    reps?: number;
    weight?: number;
    duration_s?: number;
  }

  export interface Setting {
    grade_system: string;
    climb_toplogger_user_id?: string;
    climb_toplogger_auth_token?: string;
    climb_toplogger_refresh_token?: string;
  }

  export interface Task {
    uid: string;
    name: string;
    status: string;
    next_run?: string;
    last_status?: string;
    last_run?: string;
    last_message?: string;
    last_error?: string;
    interval?: number;
    recurring: boolean;
  }

  export interface TaskHistory {
    id: number;
    name: string;
    result: string;
    run_at: string;
    message?: string;
    error?: string;
    duration: number;
  }

  export interface Gym {
    id: number;
    name: string;
    icon_path: string;
    source: string;
  }

  export interface GymStats {
    total: number;
    most_visited: string;
    most_visited_amount: number;
    sessions: number;
    graph_visits: GymStatsVisits[];
    graph_top: GymStatsTop[];
    graph_distribution: GymStatsDistribution[];
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

  export interface Climb {
    id: number;
    grade: string;
    hold_color: string;
    climb_type: string;
    finish_type: string;
    source: string;
  }

  export interface ClimbDay {
    id: number;
    date: string;
    gym: Gym;
    climbs: Climb[];
    source: string;
  }

  export interface ClimbStatsProgress {
    date: string;
    grade: string;
    volume: number;
  }

  export interface ClimbStatsGrade {
    grade: string;
    flash: number;
    top: number;
    repeat: number;
  }

  export interface ClimbStats {
    total: number;
    total_unique: number;
    flash: number;
    top: number;
    repeat: number;
    best: string;
    best_amount: number;
    best_flash: string;
    best_flash_amount: number;
    sessions: number;
    boulder: number;
    lead: number;
    median_climbs_per_session: number;
    graph_progress: ClimbStatsProgress[];
    graph_per_grade: ClimbStatsGrade[];
  }
}
