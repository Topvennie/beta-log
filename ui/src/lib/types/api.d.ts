export namespace API {
  interface Base extends JSON {
    id: number;
  }

  export interface User extends Base {
    uid: string;
    name: string;
  }

  export interface Exercise extends Base {
    name: string;
    variants?: Variant[];
  }

  export interface Variant extends Base {
    variant: string;
  }

  export interface Session extends Base {
    name: string;
    exercises: SessionExercise[];
  }

  export interface SessionExercise extends Base {
    exercise: Omit<Exercise, "variants">;
    variant?: Variant;
    position: number;
    sets: number;
    reps?: number;
    weight?: number;
    duration_s?: number;
  }
}
