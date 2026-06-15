import { QueryClient } from "@tanstack/react-query";
import { createRootRouteWithContext, createRoute, createRouter } from "@tanstack/react-router";
import { App } from "./App";
import { queryClient } from "./lib/api/query";
import { Error404 } from "./pages/404";
import { Error } from "./pages/Error";
import { Index } from "./pages/auth/Index";

type Context = {
  queryClient: QueryClient,
}

const root = createRootRouteWithContext<Context>()({
  component: App,
})

const index = createRoute({
  getParentRoute: () => root,
  path: "/",
  component: Index,
})

const routeTree = root.addChildren([
  index,
])

export const router = createRouter({
  routeTree,
  context: {
    queryClient,
  },
  defaultPreload: "render",
  defaultPreloadStaleTime: 0, // Data is immediatly marked as stale and will refetch when the user navigates to the page
  scrollRestoration: true,
  defaultErrorComponent: Error,
  defaultNotFoundComponent: Error404,
})

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}
