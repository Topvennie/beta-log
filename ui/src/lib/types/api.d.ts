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
}
