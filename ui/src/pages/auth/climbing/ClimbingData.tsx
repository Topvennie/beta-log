import { BottomOfPage } from "@/components/atoms/BottomOfPage"
import { LinkButton } from "@/components/atoms/LinkButton"
import { ClimbDayForm } from "@/components/climb/ClimbDayForm"
import { LoadingLayout } from "@/layout/LoadingLayout"
import { useClimbDayCreate, useClimbDayDelete, useClimbDayGetFiltered, useClimbDayUpdate } from "@/lib/api/climb"
import { useGymGetAll } from "@/lib/api/gym"
import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb"
import { ClimbDay, ClimbDayCreate, ClimbDayUpdate, ClimbFinish } from "@/lib/types/climb"
import { Gym, Source } from "@/lib/types/gym"
import { ActionIcon, Avatar, Badge, BadgeProps, Button, Card, ColorSwatch, Divider, Group, Modal, Stack, Tooltip } from "@mantine/core"
import { useDisclosure } from "@mantine/hooks"
import { notifications } from "@mantine/notifications"
import { format } from "date-fns"
import { Fragment, useState } from "react"
import { FaGear, FaPencil, FaPlus } from "react-icons/fa6"
import useInfiniteScroll from "react-infinite-scroll-hook"

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

  const [opened, { open, close }] = useDisclosure()
  const [selected, setSelected] = useState<ClimbDay | null>(null)

  const dayCreate = useClimbDayCreate()
  const dayUpdate = useClimbDayUpdate()
  const dayDelete = useClimbDayDelete()

  const handleCreate = (day: ClimbDayCreate) => {
    return dayCreate.mutateAsync(day, {
      onSuccess: () => {
        notifications.show({ color: "green", title: "Climb Day", message: `Created` })
        handleClose()
      }
    })
  }

  const handleUpdate = (day: ClimbDayUpdate) => {
    return dayUpdate.mutateAsync(day, {
      onSuccess: () => {
        notifications.show({ color: "green", title: "Climb Day", message: `Updated` })
        handleClose()
      }
    })
  }

  const handleDelete = ({ id }: Pick<ClimbDay, "id">) => {
    return dayDelete.mutateAsync({ id }, {
      onSuccess: () => {
        notifications.show({ color: "green", title: "ClimbDay", message: `Deleted` })
        handleClose()
      },
    })
  }

  const handleSelect = (s: ClimbDay | null = null) => {
    setSelected(s)
    open()
  }

  const handleClose = () => {
    setSelected(null)
    close()
  }

  return (
    <LoadingLayout isLoading={isLoadingGyms || isLoadingDays}>
      <Stack>
        <Gyms gyms={gyms ?? []} />

        <Group justify="space-between">
          <p className="font-bold">Climbing Days</p>
          <Button onClick={() => handleSelect()} variant="outline" leftSection={<FaPlus />}>
            Add Day
          </Button>
        </Group>

        {days.map(d => <Day key={d.id} onClick={handleSelect} day={d} />)}

        <BottomOfPage ref={sentryRef} showLoading={isFetchingNextPage} hasNextPage={hasNextPage} />
      </Stack>

      <Modal title="Climbing Day" opened={opened} onClose={handleClose}>
        {selected
          ? <ClimbDayForm climbDay={selected} onSubmit={handleUpdate} onCancel={handleClose} onDelete={handleDelete} />
          : <ClimbDayForm onSubmit={handleCreate} onCancel={handleClose} />
        }
      </Modal>
    </LoadingLayout>
  )
}

const Gyms = ({ gyms }: { gyms: Gym[] }) => {
  return (
    <Group>
      <p className="text-neutral-400">Gyms:</p>
      <p>{gyms.map(g => g.name).join(", ")}</p>
      <LinkButton to="/gym/data" variant="subtle" leftSection={<FaGear />}>
        Manage
      </LinkButton>
    </Group>
  )
}

const Day = ({ day, onClick }: { day: ClimbDay, onClick: (day: ClimbDay) => void }) => {
  const external = day.source !== Source.Manual

  const finishProps: Record<ClimbFinish, Partial<BadgeProps>> = {
    [ClimbFinish.Flash]: {},
    [ClimbFinish.Top]: { variant: "light" },
    [ClimbFinish.Repeat]: { variant: "light", color: "black" },
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

          <Group>
            <p className="text-neutral-400">{`${day.climbs.length} climbs`}</p>
            <Tooltip label="Day is managed externally" disabled={!external}>
              <ActionIcon onClick={() => onClick(day)} color="black" variant="subtle" disabled={external}><FaPencil /></ActionIcon>
            </Tooltip>
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
