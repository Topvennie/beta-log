import type { ReactNode } from "react";
import { use, useId, useLayoutEffect } from "react";
import { HeaderContentContext } from "../contexts/headerContentContext";

export function useHeaderContents() {
  const context = use(HeaderContentContext);
  if (!context) {
    throw new Error("useHeaderContents must be used within a HeaderContentProvider");
  }

  return context;
}

export function useHeaderContent(content: ReactNode) {
  const { dispatch } = useHeaderContents();
  const id = useId();

  // useLayoutEffect so dispatches flush before the browser paints —
  // keeps header content updates in the same frame as the route change.
  useLayoutEffect(() => {
    dispatch({ type: "PUSH", payload: { id, content } });

    return () => {
      dispatch({ type: "REMOVE", payload: { id } });
    };

    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [id]);
}
