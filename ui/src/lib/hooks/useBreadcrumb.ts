import type { Breadcrumb } from "../contexts/breadcrumbContext";
import { use, useLayoutEffect } from "react";
import { BreadcrumbContext } from "../contexts/breadcrumbContext";

export function useBreadcrumbs() {
  const context = use(BreadcrumbContext);
  if (!context) {
    throw new Error("useBreadcrumbs must be used within a BreadcrumbProvider");
  }

  return context;
}

export function useBreadcrumb(crumb: Breadcrumb) {
  const { dispatch } = useBreadcrumbs();

  useLayoutEffect(() => {
    dispatch({ type: "ADD", payload: crumb });

    return () => {
      dispatch({ type: "REMOVE", payload: crumb });
    };

    // Use stringify to avoid infinite rerenders
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [JSON.stringify(crumb), dispatch]);
}
