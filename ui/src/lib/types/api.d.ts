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
    variant: string;
  }

  export interface Session extends Base {
    name: string;
    active: boolean;
    position?: number;
    exercises: SessionExercise[];
  }

  export interface SessionExercise extends Base {
    name: string;
    variant: string;
    position: number;
    sets: number;
    reps?: number;
    weight?: number;
    duration_s?: number;
  }
}
