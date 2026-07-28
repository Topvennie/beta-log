import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb"
import { Outlet } from "@tanstack/react-router"

export const Gym = () => {
  useBreadcrumb({ title: "Gym", weight: 10, link: { to: "/gym" } })

  return <Outlet />
}
