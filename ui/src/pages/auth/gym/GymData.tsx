import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb"

export const GymData = () => {
  useBreadcrumb({ title: "Manage Data", weight: 20, link: { to: "/gym/data" } })

  return (
    <div>
      Gym
    </div>
  )
}
