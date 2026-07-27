import { LoadingLayout } from "@/layout/LoadingLayout"
import useInfiniteScroll from "react-infinite-scroll-hook"
import { useClimbDayGetFiltered, useClimbGymCreate, useClimbGymGetAll, useClimbGymUpdate } from "@/lib/api/climb"
import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb"
import { ClimbDay, ClimbFinish, ClimbGym, ClimbGymCreate, ClimbGymUpdate } from "@/lib/types/climb"
import { ActionIcon, Avatar, Badge, BadgeProps, Button, Card, ColorSwatch, Divider, Group, Modal, Scroller, Stack, Tabs } from "@mantine/core"
import { format } from "date-fns"
import { Fragment, useState } from "react"
import { FaGear, FaPencil, FaPlus, FaTrashCan } from "react-icons/fa6"
import { BottomOfPage } from "@/components/atoms/BottomOfPage"
import { useDisclosure } from "@mantine/hooks"
import { ClimbGymForm } from "@/components/climb/ClimbGymForm"
import { notifications } from "@mantine/notifications"

export const ClimbingData = () => {
  useBreadcrumb({ title: "Manage Data", weight: 20, link: { to: "/climbing/data" } })

  const { days, isLoading: isLoadingDays, isFetchingNextPage, hasNextPage, fetchNextPage } = useClimbDayGetFiltered()
  const { data: gyms, isLoading: isLoadingGyms } = useClimbGymGetAll()

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
          <Button leftSection={<FaPlus />}>
            Add Day
          </Button>
        </Group>
        {days.map(d => <Day key={d.id} day={d} />)}
        <BottomOfPage ref={sentryRef} showLoading={isFetchingNextPage} hasNextPage={hasNextPage} />
      </Stack>
    </LoadingLayout>
  )
}

const Gyms = ({ gyms }: { gyms: ClimbGym[] }) => {
  const [opened, { open, close }] = useDisclosure()

  const [activeTab, setActiveTab] = useState<string | null>(null)
  const [selected, setSelected] = useState<ClimbGym | undefined>()

  const handleActiveTab = (value: string | null) => {
    setActiveTab(value)
    setSelected(value ? gyms.find(g => g.id === Number(value)) : undefined)
  }

  const gymCreate = useClimbGymCreate()
  const gymUpdate = useClimbGymUpdate()

  const handleCreate = (gym: ClimbGymCreate) => {
    return gymCreate.mutateAsync(gym, {
      onSuccess: (resp) => {
        notifications.show({ color: "green", title: "Gym", message: `Created ${gym.name}` })
        handleActiveTab(resp.data.id.toString())
      }
    })
  }

  const handleUpdate = (gym: ClimbGymUpdate) => {
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
          <Group justify="space-between">
            <Tabs value={activeTab} onChange={handleActiveTab}>
              <Scroller>
                {gyms.map(g => (
                  <Tabs.Tab key={g.id} value={g.id.toString()} disabled={g.source !== "manual"}>{g.name}</Tabs.Tab>
                ))}
              </Scroller>
            </Tabs>
            <ActionIcon onClick={() => handleActiveTab(null)}><FaPlus /></ActionIcon>
          </Group>
          {selected
            ? <ClimbGymForm key={selected.id} gym={selected} onSubmit={handleUpdate} />
            : <ClimbGymForm gym={undefined} onSubmit={handleCreate} />
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
    <Card className="border border-gray-200">
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
        <Card.Section className="bg-neutral-200/40 border-t border-b border-gray-200">
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
