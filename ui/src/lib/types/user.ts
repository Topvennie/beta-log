import { API } from "./api";

export interface User {
  id: number;
  uid: number;
  name: string;
}

// Converters

export const convertUserToModel = (u: API.User): User => ({
  id: u.id,
  uid: u.uid,
  name: u.name,
})

