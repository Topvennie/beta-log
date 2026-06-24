import { PropsWithChildren, useEffect, useRef, useState } from "react"

type Props = {
  isLoading: boolean
} & PropsWithChildren

const SHOW_DELAY_MS = 400
const MIN_VISIBLE_MS = 900
const FADE_DURATION_MS = 250

type Phase = "blank" | "loading" | "fading" | "content"

// LoadingLayout provides a nice looking loading screen
// The logic is as following:
// 1. If isLoading === true show a blank screen
// 2. If isLoading is still true after SHOW_DELAY_MS then show the loading screen (prevents unnecessary flickering for fast loading data)
// 3. When isLoading === false and after MIN_VISIBLE_MS start fading out the loading screen for FADE_DURATION_MS (again to prevent flickering)
// 4. Show the child
export const LoadingLayout = ({ isLoading, children }: Props) => {
  const [phase, setPhase] = useState<Phase>(isLoading ? "blank" : "content")

  const shownAtRef = useRef<number | null>(null)
  const timerRef = useRef<number | null>(null)

  useEffect(() => {
    const clearTimer = () => {
      if (timerRef.current === null) return
      clearTimeout(timerRef.current)
      timerRef.current = null
    }

    clearTimer()

    // Child is still loading data
    if (isLoading) {
      // Child was already shown
      // But it is loading some extra data
      // So hide it and restart the entire cycle
      // You probably don't want this to happen though. You can avoid it by not passing
      // the isLoading props from the extra data
      if (phase === "content" || phase === "fading") {
        shownAtRef.current = null
        setTimeout(() => {
          setPhase("blank")
        }, 0)
        return clearTimer
      }

      // Show the loading screen if the data isn't loaded in by SHOW_DELAY_MS
      if (phase === "blank") {
        timerRef.current = setTimeout(() => {
          shownAtRef.current = Date.now()
          setTimeout(() => {
            setPhase("loading")
          }, 0)
        }, SHOW_DELAY_MS)
        return clearTimer
      }

      return clearTimer
    }

    // The data is loaded in
    // The phase determines what we do next

    // We're still before SHOW_DELAY_MS
    // So show the child immediately
    if (phase === "blank") {
      setTimeout(() => {
        setPhase("content")
      }, 0)
      return clearTimer
    }

    // We're already showing the loading screen
    // Start the fading if MIN_VISIBLE_MS is passed
    if (phase === "loading") {
      const shownAt = shownAtRef.current ?? Date.now()
      const elapsed = Date.now() - shownAt
      const remaining = Math.max(0, MIN_VISIBLE_MS - elapsed)

      timerRef.current = setTimeout(() => {
        setTimeout(() => {
          setPhase("fading")
        }, 0)
      }, remaining)

      return clearTimer
    }

    // We're fading the loading screen
    // Show the child after FADE_DURATION_MS
    if (phase === "fading") {
      timerRef.current = setTimeout(() => {
        setTimeout(() => {
          setPhase("content")
        }, 0)
        shownAtRef.current = null
      }, FADE_DURATION_MS)

      return clearTimer
    }

    return clearTimer
  }, [isLoading, phase])

  if (phase === "content") return <>{children}</> // Keep the return type normal (thx typescript)
  if (phase === "blank") return null

  return (
    <div
      className="transition-opacity"
      style={{
        opacity: phase === "fading" ? 0 : 1,
        transitionDuration: `${FADE_DURATION_MS}ms`,
      }}
    >
      <p className="text-sm md:text-base lg:text-lg text-center">
        Loading...
      </p>
    </div>
  )
}
