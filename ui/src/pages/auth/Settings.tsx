import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb"

export const Settings = () => {
  useBreadcrumb({ title: "Settings", weight: 10, link: { to: "/settings" } })

  return (
    <div>Settings</div>
  )
}
