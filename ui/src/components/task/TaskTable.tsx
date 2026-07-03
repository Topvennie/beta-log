import { useTaskStart } from "@/lib/api/task";
import { Task, TaskStatus } from "@/lib/types/task";
import { Button } from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { formatDistanceToNow, formatDuration, intervalToDuration } from "date-fns";
import { FaPlay } from "react-icons/fa6";
import { Table } from "../molecules/Table";
import { TaskResult } from "./TaskResult";

type Props = {
  tasks: Task[];
}

const formatInterval = (ms: number) => {
  const duration = intervalToDuration({ start: 0, end: ms / 1000000 });
  return formatDuration(duration);
}

export const TaskTable = ({ tasks }: Props) => {
  const taskRun = useTaskStart()

  const handleClick = (task: Task) => {
    taskRun.mutate(task, {
      onSuccess: () => notifications.show({ title: task.name, message: "Started" })
    })
  }

  return (
    <Table
      idAccessor="uid"
      columns={[
        { accessor: "name", title: "Task" },
        { accessor: "lastRun", render: ({ lastRun }) => <p className="text-muted">{lastRun ? formatDistanceToNow(lastRun, { addSuffix: true }) : ""}</p> },
        { accessor: "nextRun", render: ({ nextRun }) => <p className="text-muted">{nextRun ? formatDistanceToNow(nextRun, { addSuffix: true }) : ""}</p> },
        { accessor: "interval", render: ({ interval }) => <p className="text-muted">{interval ? formatInterval(interval) : ""}</p> },
        {
          accessor: "lastStatus",
          title: "Last result",
          render: ({ lastStatus }) => lastStatus && <TaskResult result={lastStatus} />
        },
        {
          accessor: "uid",
          title: "",
          width: 100,
          render: task => <Button size="compact-sm" leftSection={<FaPlay />} onClick={() => handleClick(task)} loading={task.status === TaskStatus.Running} disabled={!task.recurring}>Run now</Button>,
        }
      ]}
      records={tasks}
    />
  )
}
