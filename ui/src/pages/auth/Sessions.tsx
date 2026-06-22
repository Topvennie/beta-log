import { useSessionCreate, useSessionDelete, useSessionGetAll, useSessionUpdate } from "@/lib/api/session"
import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb"
import { Group, Stack, Title } from "@mantine/core"
import { SessionCreate, Session as SessionType, SessionUpdate } from "@/lib/types/session"
import { SessionForm } from "@/components/session/SessionForm"
import { notifications } from "@mantine/notifications"
import { useState } from "react"

export const Sessions = () => {
  useBreadcrumb({ title: "Sessions", weight: 10, link: { to: "/sessions" } })

  const { data: sessions, isLoading } = useSessionGetAll()
  const sessionCreate = useSessionCreate()
  const sessionUpdate = useSessionUpdate()
  const sessionDelete = useSessionDelete()

  const [selected, setSelected] = useState<SessionType | undefined>(undefined)

  if (isLoading || sessions === undefined) return null

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
    <div className="flex md:flex-row gap-8 h-full">
      <Stack gap="xl" className="flex-3 h-full">
        <Group justify="space-between">
          <Title order={1}>Sessions</Title>
          <p className="text-neutral-400">{sessions.length}</p>
        </Group>
      </Stack>

      <div className="border-l border-gray-200" />

      <div className="flex-2">
        {selected
          ? <SessionForm key={selected.id} session={selected} onSubmit={handleUpdate} onCancel={() => setSelected(undefined)} onDelete={handleDelete} />
          : <SessionForm onSubmit={handleCreate} onCancel={() => null} />
        }
      </div>

    </div>
  )
}
