import { LoadingLayout } from "@/layout/LoadingLayout"
import useInfiniteScroll from "react-infinite-scroll-hook"
import { useClimbDayGetFiltered, useClimbGymGetAll } from "@/lib/api/climb"
import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb"
import { ClimbDay, ClimbFinish, ClimbGym } from "@/lib/types/climb"
import { ActionIcon, Avatar, Badge, BadgeProps, Button, Card, Divider, Group, Stack } from "@mantine/core"
import { format } from "date-fns"
import { Fragment } from "react"
import { FaGear, FaPencil, FaPlus, FaTrashCan } from "react-icons/fa6"
import { BottomOfPage } from "@/components/atoms/BottomOfPage"

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
  return (
    <Group>
      <p className="text-neutral-400">Gyms:</p>
      <p>{gyms.map(g => g.name).join(", ")}</p>
      <Button variant="subtle" leftSection={<FaGear />}>
        Manage
      </Button>
    </Group>
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
            <Avatar src={day.gym.iconPath} name={day.gym.name} />
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
                  <div className="w-4 h-4 rounded-full border border-neutral-200" style={{ background: c.holdColor }} />
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
