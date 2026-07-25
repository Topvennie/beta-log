import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb"

export const ClimbingData = () => {
  useBreadcrumb({ title: "Manage Data", weight: 20, link: { to: "/climbing/data" } })

  return (
    <div>Data</div>
  )
}
