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

  export interface Climb {
    id: number;
    grade: number;
    hold_color: string;
    climb_type: string;
    finish_type: string;
  }

  export interface ClimbGym {
    id: number;
    name: string;
    icon_path: string;
  }

  export interface ClimbDay {
    id: number;
    date: string;
    gym: ClimbGym;
    climbs: Climb[];
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
    total_unique: number;
    flash: number;
    top: number;
    repeat: number;
    best: number;
    best_amount: number;
    best_flash: number;
    best_flash_amount: number;
    sessions: number;
    climbs_per_session: number;
    graph_progress: ClimbStatsGraphProgress[];
    graph_per_grade: ClimbStatsGraphGrade[];
  }
}
