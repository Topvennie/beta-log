export namespace API {
  interface Base extends JSON {
    id: number;
  }

  export interface User extends Base {
    uid: number;
    name: string;
  }
}
