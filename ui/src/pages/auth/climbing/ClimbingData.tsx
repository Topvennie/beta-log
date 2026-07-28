import { LoadingLayout } from "@/layout/LoadingLayout"
import useInfiniteScroll from "react-infinite-scroll-hook"
import { useClimbDayGetFiltered } from "@/lib/api/climb"
import { useGymCreate, useGymGetAll, useGymUpdate } from "@/lib/api/gym"
import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb"
import { ClimbDay, ClimbFinish } from "@/lib/types/climb"
import { Gym, GymCreate, GymUpdate } from "@/lib/types/gym"
import { ActionIcon, Avatar, Badge, BadgeProps, Button, Card, ColorSwatch, Divider, Group, Modal, Scroller, Stack, Tabs } from "@mantine/core"
import { format } from "date-fns"
import { Fragment, useState } from "react"
import { FaGear, FaPencil, FaPlus, FaTrashCan } from "react-icons/fa6"
import { BottomOfPage } from "@/components/atoms/BottomOfPage"
import { useDisclosure } from "@mantine/hooks"
import { GymForm } from "@/components/gym/GymForm"
import { notifications } from "@mantine/notifications"

export const ClimbingData = () => {
  useBreadcrumb({ title: "Manage Data", weight: 20, link: { to: "/climbing/data" } })

  const { days, isLoading: isLoadingDays, isFetchingNextPage, hasNextPage, fetchNextPage } = useClimbDayGetFiltered()
  const { data: gyms, isLoading: isLoadingGyms } = useGymGetAll()

  const [sentryRef] = useInfiniteScroll({
    loading: isFetchingNextPage,
    hasNextPage: Boolean(hasNextPage),
    onLoadMore: fetchNextPage,
    rootMargin: "0px",
  });

  return (
    <LoadingLayout isLoading={isLoadingGyms || isLoadingDays}>
      <Stack>
        <Gyms gyms={gyms ?? []} />
        <Group justify="space-between">
          <p className="font-bold">Climbing Days</p>
          <Button variant="outline" leftSection={<FaPlus />}>
            Add Day
          </Button>
        </Group>
        {days.map(d => <Day key={d.id} day={d} />)}
        <BottomOfPage ref={sentryRef} showLoading={isFetchingNextPage} hasNextPage={hasNextPage} />
      </Stack>
    </LoadingLayout>
  )
}

const Gyms = ({ gyms }: { gyms: Gym[] }) => {
  const [opened, { open, close }] = useDisclosure()

  const [activeTab, setActiveTab] = useState<string | null>(null)
  const [selected, setSelected] = useState<Gym | undefined>()

  const handleActiveTab = (value: string | null) => {
    setActiveTab(value)
    setSelected(value ? gyms.find(g => g.id === Number(value)) : undefined)
  }

  const gymCreate = useGymCreate()
  const gymUpdate = useGymUpdate()

  const handleCreate = (gym: GymCreate) => {
    return gymCreate.mutateAsync(gym, {
      onSuccess: (resp) => {
        notifications.show({ color: "green", title: "Gym", message: `Created ${gym.name}` })
        setSelected(resp.data)
        setActiveTab(resp.data.id.toString())
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

  return (
    <>
      <Group>
        <p className="text-neutral-400">Gyms:</p>
        <p>{gyms.map(g => g.name).join(", ")}</p>
        <Button onClick={open} variant="subtle" leftSection={<FaGear />}>
          Manage
        </Button>
      </Group>

      <Modal opened={opened} onClose={close}>
        <Stack>
          <Group justify="space-between" wrap="nowrap" gap="xs">
            <Tabs value={activeTab} onChange={handleActiveTab} style={{ flex: 1, minWidth: 0, overflow: "hidden" }}>
              <Tabs.List style={{ minWidth: 0 }}>
                <Scroller draggable style={{ flex: 1, minWidth: 0 }}>
                  {gyms.map(g => (
                    <Tabs.Tab key={g.id} value={g.id.toString()} disabled={g.source !== "manual"}>{g.name}</Tabs.Tab>
                  ))}
                </Scroller>
              </Tabs.List>
            </Tabs>
            <ActionIcon onClick={() => handleActiveTab(null)}><FaPlus /></ActionIcon>
          </Group>
          {selected
            ? <GymForm key={selected.id} gym={selected} onSubmit={handleUpdate} />
            : <GymForm gym={undefined} onSubmit={handleCreate} />
          }

        </Stack>
      </Modal>
    </>
  )
}

const Day = ({ day }: { day: ClimbDay }) => {
  const finishProps: Record<ClimbFinish, Partial<BadgeProps>> = {
    "flash": {},
    "top": { variant: "light" },
    "repeat": { variant: "light", color: "black" },
  }

  return (
    <Card className="border border-neutral-200">
      <Stack>
        <Group justify="space-between">
          <Group>
            <Avatar src={day.gym.iconPath} name={day.gym.name} className="border" />
            <Stack gap={2}>
              <p>{format(day.date, "EEE dd MMM yyyy")}</p>
              <p className="text-neutral-400">{day.gym.name}</p>
            </Stack>
          </Group>
          <Group gap={2}>
            <p className="text-neutral-400">{`${day.climbs.length} climbs`}</p>
            <ActionIcon color="black" variant="subtle"><FaPencil /></ActionIcon>
            <ActionIcon color="black" variant="subtle"><FaTrashCan /></ActionIcon>
          </Group>
        </Group>
        <Card.Section className="bg-neutral-200/40 border-t border-b border-neutral-200">
          <div className="grid grid-cols-4 px-md py-xs text-neutral-400">
            <p>Grade</p>
            <p>Hold</p>
            <p>Type</p>
            <p>Finish</p>
          </div>
        </Card.Section>
        <div className="overflow-y-auto max-h-56 scroll-">
          <Stack gap={0}>
            {day.climbs.map((c, i) => (
              <Fragment key={c.id}>
                {i > 0 && <Divider />}
                <div className="grid grid-cols-4 py-xs">
                  <p>{c.grade}</p>
                  <ColorSwatch color={c.holdColor} size={18} />
                  <p className="text-neutral-400">{c.climbType}</p>
                  <Badge size="sm" {...finishProps[c.finishType]}>{c.finishType}</Badge>
                </div>
              </Fragment>
            ))}
          </Stack>
        </div>
      </Stack>
    </Card>
  )
}
