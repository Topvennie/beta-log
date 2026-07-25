import { createContext, ReactNode } from "react";

export interface HeaderContentItem {
  id: string;
  content: ReactNode;
}

export type HeaderContentAction =
  | { type: "PUSH"; payload: HeaderContentItem }
  | { type: "REMOVE"; payload: { id: string } };

export type HeaderContentState = HeaderContentItem[];

export const HeaderContentContext = createContext<{
  content: ReactNode;
  dispatch: React.Dispatch<HeaderContentAction>;
} | undefined>(undefined);
