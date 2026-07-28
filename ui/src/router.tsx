import { QueryClient } from "@tanstack/react-query";
import { createRootRouteWithContext, createRoute, createRouter } from "@tanstack/react-router";
import { App } from "./App";
import { queryClient } from "./lib/api/query";
import { Error404 } from "./pages/404";
import { Error } from "./pages/Error";
import { Index } from "./pages/auth/Index";
import { Exercises } from "./pages/auth/Exercises";
import { Dashboard } from "./pages/auth/Dashboard";
import { Sessions } from "./pages/auth/Sessions";
import { Settings } from "./pages/auth/Settings";
import { Tasks } from "./pages/auth/Tasks";
import { Climbing } from "./pages/auth/climbing/Climbing";
import { ClimbingData } from "./pages/auth/climbing/ClimbingData";
import { ClimbingDashboard } from "./pages/auth/climbing/ClimbingDashboard";
import { Gym } from "./pages/auth/gym/Gym";
import { GymDashboard } from "./pages/auth/gym/GymDashboard";
import { GymData } from "./pages/auth/gym/GymData";

type Context = {
  queryClient: QueryClient,
}

const root = createRootRouteWithContext<Context>()({
  component: App,
})

const index = createRoute({
  getParentRoute: () => root,
  id: "public-layout",
  component: Index,
})

const dashboard = createRoute({
  getParentRoute: () => index,
  path: "/",
  component: Dashboard,
})

const exercises = createRoute({
  getParentRoute: () => index,
  path: "/exercises",
  component: Exercises,
})

const sessions = createRoute({
  getParentRoute: () => index,
  path: "/sessions",
  component: Sessions,
})

const tasks = createRoute({
  getParentRoute: () => index,
  path: "/tasks",
  component: Tasks,
})

// Climbing

const climbing = createRoute({
  getParentRoute: () => index,
  path: "/climbing",
  component: Climbing,
})

const climbingDashboard = createRoute({
  getParentRoute: () => climbing,
  path: "/",
  component: ClimbingDashboard,
})

const climbingData = createRoute({
  getParentRoute: () => climbing,
  path: "/data",
  component: ClimbingData,
})

// Gym

const gym = createRoute({
  getParentRoute: () => index,
  path: "/gym",
  component: Gym,
})

const gymDashboard = createRoute({
  getParentRoute: () => gym,
  path: "/",
  component: GymDashboard,
})

const gymData = createRoute({
  getParentRoute: () => gym,
  path: "/data",
  component: GymData,
})

const settings = createRoute({
  getParentRoute: () => index,
  path: "/settings",
  component: Settings,
})

const routeTree = root.addChildren([
  index.addChildren([
    dashboard,
    exercises,
    sessions,
    tasks,
    climbing.addChildren([climbingDashboard, climbingData]),
    gym.addChildren([gymDashboard, gymData]),
    settings,
  ]),
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
