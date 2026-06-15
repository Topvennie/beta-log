import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb"

export const Exercises = () => {
  useBreadcrumb({ title: "Exercises", weight: 10, link: { to: "/exercises" } })

  return (
    <p>Exercises</p>
  )
}
