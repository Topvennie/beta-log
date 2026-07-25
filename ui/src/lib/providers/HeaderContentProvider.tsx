import type { ReactNode } from "react";
import type { HeaderContentAction, HeaderContentState } from "../contexts/headerContentContext";
import { useMemo, useReducer } from "react";
import { HeaderContentContext } from "../contexts/headerContentContext";

function headerContentReducer(state: HeaderContentState, action: HeaderContentAction): HeaderContentState {
  switch (action.type) {
    case "PUSH": {
      if (state.some(item => item.id === action.payload.id)) {
        return state.map(item =>
          item.id === action.payload.id ? action.payload : item
        );
      }
      return [...state, action.payload];
    }
    case "REMOVE":
      return state.filter(item => item.id !== action.payload.id);
    default:
      return state;
  }
}

export function HeaderContentProvider({ children }: { children: ReactNode }) {
  const [state, dispatch] = useReducer(headerContentReducer, []);

  // LIFO
  const content = state.length > 0 ? state[state.length - 1].content : null;

  const value = useMemo(() => ({ content, dispatch }), [content, dispatch]);

  return (
    <HeaderContentContext value={value}>
      {children}
    </HeaderContentContext>
  );
}
