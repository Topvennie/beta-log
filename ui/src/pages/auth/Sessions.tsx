import { useSessionCreate, useSessionDelete, useSessionGetAll, useSessionUpdate } from "@/lib/api/session"
import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb"
import { Group, Stack, Title } from "@mantine/core"
import { SessionCreate, SessionExercise, Session as SessionType, SessionUpdate } from "@/lib/types/session"
import { SessionForm } from "@/components/session/SessionForm"
import { notifications } from "@mantine/notifications"
import { useState } from "react"
import { useHover } from "@mantine/hooks"
import { FaChevronRight } from "react-icons/fa6"
import { useExerciseGetAll } from "@/lib/api/exercise"
import { LoadingLayout } from "@/layout/LoadingLayout"

export const Sessions = () => {
  useBreadcrumb({ title: "Sessions", weight: 10, link: { to: "/sessions" } })

  const { data: sessions, isLoading: isLoadingSessions } = useSessionGetAll()
  const { data: exercises, isLoading: isLoadingExercises } = useExerciseGetAll()

  const sessionCreate = useSessionCreate()
  const sessionUpdate = useSessionUpdate()
  const sessionDelete = useSessionDelete()

  const [selected, setSelected] = useState<SessionType | undefined>(undefined)

  if (!isLoadingExercises && exercises?.length === 0) return <NoExercises />

  const handleCreate = (session: SessionCreate) => {
    return sessionCreate.mutateAsync(session, {
      onSuccess: () => {
        notifications.show({ color: "green", title: "Session", message: `Created ${session.name}` })
        setSelected(undefined)
      },
    })
  }

  const handleUpdate = (session: SessionUpdate) => {
    return sessionUpdate.mutateAsync(session, {
      onSuccess: () => {
        notifications.show({ color: "green", title: "Session", message: `Updated ${session.name}` })
        setSelected(undefined)
      },
    })
  }

  const handleDelete = ({ id }: Pick<SessionType, "id">) => {
    return sessionDelete.mutateAsync({ id }, {
      onSuccess: () => {
        notifications.show({ color: "green", title: "Session", message: `Deleted` })
        setSelected(undefined)
      },
    })
  }

  return (
    <LoadingLayout isLoading={isLoadingSessions || isLoadingExercises}>
      <div className="flex md:flex-row gap-8 h-full">
        <Stack gap="xl" className="flex-3 h-full">
          <Group justify="space-between">
            <Title order={1}>Sessions</Title>
            <p className="text-neutral-400">{sessions?.length}</p>
          </Group>

          <Stack gap="md">
            {sessions?.map(s => <Session key={s.id} session={s} onClick={setSelected} />)}
          </Stack>
        </Stack>

        <div className="border-l border-gray-200" />

        <div className="flex-2">
          {selected
            ? <SessionForm key={selected.id} session={selected} onSubmit={handleUpdate} onCancel={() => setSelected(undefined)} onDelete={handleDelete} />
            : <SessionForm onSubmit={handleCreate} onCancel={() => null} />
          }
        </div>

      </div>
    </LoadingLayout>
  )
}

const exerciseDescription = (exercise: SessionExercise): string => {
  let base = exercise.variant ? `${exercise.exercise.name} - ${exercise.variant.variant}` : exercise.exercise.name

  base += ` | ${exercise.sets} set${exercise.sets !== 1 ? "s" : ""}`
  if (exercise.reps) base += ` of ${exercise.reps} rep${exercise.reps !== 1 ? "s" : ""}`
  if (exercise.durationS) base += ` for ${exercise.durationS} second${exercise.durationS !== 1 ? "s" : ""}`
  if (exercise.weight) base += ` with +${exercise.weight} kg`

  return base
}

type SessionProps = {
  session: SessionType;
  onClick: (session: SessionType) => void;
}

const Session = ({ session, onClick }: SessionProps) => {
  const { hovered, ref } = useHover()

  return (
    <Group ref={ref} onClick={() => onClick(session)} p="xs" justify="space-between" className="rounded-sm cursor-pointer hover:bg-blue-50">
      <Stack>
        <p>{session.name}</p>
        <ol>
          {session.exercises.map((e, idx) => <li key={e.id}>{`${idx + 1}. ${exerciseDescription(e)}`}</li>)}
        </ol>
      </Stack>
      <FaChevronRight className={`text-blue-500 duration-300 ${hovered ? "translate-x-1" : ""}`} />
    </Group>
  )
}

const NoExercises = () => {
  return (
    <div>
      <p>There are no exercises yet.</p>
      <p>Make some first and then return to this page.</p>
    </div>
  )
}
