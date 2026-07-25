import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb"
import { Outlet } from "@tanstack/react-router"

export const Climbing = () => {
  useBreadcrumb({ title: "Climbing", weight: 10, link: { to: "/climbing" } })

  return <Outlet />
}
