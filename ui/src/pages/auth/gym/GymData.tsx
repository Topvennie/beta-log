import { GymForm } from "@/components/gym/GymForm"
import { LoadingLayout } from "@/layout/LoadingLayout"
import { useGymCreate, useGymDelete, useGymGetAll, useGymUpdate } from "@/lib/api/gym"
import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb"
import { Gym, GymCreate, GymUpdate } from "@/lib/types/gym"
import { cn } from "@/lib/utils"
import { Avatar, Button, Card, CardProps, Group, Stack, Tooltip } from "@mantine/core"
import { notifications } from "@mantine/notifications"
import { useState } from "react"
import { FaPlus } from "react-icons/fa6"

export const GymData = () => {
  useBreadcrumb({ title: "Manage Data", weight: 20, link: { to: "/gym/data" } })

  const { data: gyms, isLoading } = useGymGetAll()

  const [selected, setSelected] = useState<Gym | null>()
  const [resetKey, setResetKey] = useState(0)

  const gymCreate = useGymCreate()
  const gymUpdate = useGymUpdate()
  const gymDelete = useGymDelete()

  const handleCreate = (gym: GymCreate) => {
    return gymCreate.mutateAsync(gym, {
      onSuccess: () => {
        notifications.show({ color: "green", title: "Gym", message: `Created ${gym.name}` })
        setSelected(null)
        setResetKey(k => k + 1)
      }
    })
  }

  const handleUpdate = (gym: GymUpdate) => {
    return gymUpdate.mutateAsync(gym, {
      onSuccess: () => {
        notifications.show({ color: "green", title: "Gym", message: `Updated ${gym.name}` })
      }
    })
  }

  const handleDelete = ({ id }: Pick<Gym, "id">) => {
    return gymDelete.mutateAsync({ id }, {
      onSuccess: () => {
        notifications.show({ color: "green", title: "Gym", message: `Deleted ${selected?.name}` })
        setSelected(null)
      },
    })
  }

  return (
    <LoadingLayout isLoading={isLoading}>
      <div className="flex gap-8">
        <Stack className="flex-1">
          <Group justify="space-between">
            <p className="text-xl font-bold">Gyms</p>
            <Button onClick={() => { setSelected(null); setResetKey(k => k + 1) }} variant="outline" leftSection={<FaPlus />}>Add Gym</Button>
          </Group>
          {gyms?.map(g => <GymCard key={g.id} gym={g} selected={g.id === selected?.id} onClick={setSelected} />)}
        </Stack>

        <Stack className="flex-1">
          <p className="text-xl font-bold">Gym Details</p>
          {selected
            ? <GymForm key={selected.id} gym={selected} onSubmit={handleUpdate} onCancel={() => setSelected(null)} onDelete={handleDelete} />
            : <GymForm key={`create-${resetKey}`} onSubmit={handleCreate} onCancel={() => setSelected(null)} />
          }
        </Stack>
      </div>
    </LoadingLayout>
  )
}

type GymCardProps = {
  gym: Gym;
  selected: boolean;
  onClick: (gym: Gym) => void;
} & CardProps

const GymCard = ({ gym, selected, onClick, className, ...props }: GymCardProps) => {
  const external = gym.source !== "manual"

  const handleClick = () => {
    if (external) return
    onClick(gym)
  }

  return (
    <Tooltip.Floating label="Gym is managed externally" disabled={!external}>
      <Card
        onClick={handleClick}
        className={cn(
          "border",
          selected ? "border-blue-500 bg-blue-500/20" : "border-neutral-200",
          external ? "opacity-40 cursor-not-allowed" : "cursor-pointer",
          className,
        )}
        {...props}
      >
        <Group>
          <Avatar src={gym.iconPath} name={gym.name} className="border" />
          <Stack gap={2}>
            <p>{gym.name}</p>
          </Stack>
        </Group>
      </Card>
    </Tooltip.Floating>
  )
}
