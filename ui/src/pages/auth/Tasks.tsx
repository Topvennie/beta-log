import { TaskHistory } from "@/components/task/TaskHistoryTable"
import { TaskTable } from "@/components/task/TaskTable"
import { LoadingLayout } from "@/layout/LoadingLayout"
import { useTaskGetAll } from "@/lib/api/task"
import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb"
import { TaskHistoryFilter, TaskResult } from "@/lib/types/task"
import { capitalize } from "@/lib/utils"
import { Group, SegmentedControl, Select, Stack, Title } from "@mantine/core"
import { useState } from "react"

export const Tasks = () => {
  useBreadcrumb({ title: "Tasks", weight: 10, link: { to: "/tasks" } })

  const { data: tasks, isLoading: isLoadingTasks } = useTaskGetAll()

  const [filter, setFilter] = useState<TaskHistoryFilter>({});

  const handleTaskChange = (value: string | null) => {
    setFilter({ ...filter, uid: value ? value : undefined })
  }

  const handleResultChange = (value: string) => {
    let result: TaskResult | undefined = undefined

    if (value !== "all") result = value as TaskResult

    setFilter({ ...filter, result: result })
  }

  const handleRecurringChange = (value: string) => {
    let recurring: boolean | undefined = undefined

    if (value !== "all") recurring = value === "true"

    setFilter({ ...filter, recurring })
  }

  return (
    <LoadingLayout isLoading={isLoadingTasks}>
      <Stack gap="xl">
        <Group justify="space-between">
          <Title order={1}>Tasks</Title>
          <p className="text-neutral-400">{tasks?.length}</p>
        </Group>

        <TaskTable tasks={tasks ?? []} />

        <Group justify="space-between">
          <Title order={2}>History</Title>
          <Group>
            <Select
              data={tasks?.map(t => ({ value: t.uid, label: t.name }))}
              value={filter.uid}
              onChange={handleTaskChange}
              placeholder="Filter by task name..."
              disabled={isLoadingTasks}
            />
            <SegmentedControl
              data={[
                { value: "all", label: "All" },
                ...Object.values(TaskResult).map(r => ({ value: r, label: capitalize(r) }))
              ]}
              value={filter.result ? filter.result : "all"}
              onChange={handleResultChange}
            />
            <SegmentedControl
              data={[
                { value: "all", label: "All" },
                { value: "true", label: "Recurring" },
                { value: "false", label: "Non-recurring" },
              ]}
              value={filter.recurring !== undefined ? String(filter.recurring) : "all"}
              onChange={handleRecurringChange}
            />
          </Group>
        </Group>

        <TaskHistory filter={filter} />
      </Stack>
    </LoadingLayout>
  )
}
