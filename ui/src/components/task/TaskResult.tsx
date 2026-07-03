import { TaskResult as TaskResultType } from "@/lib/types/task"
import { Badge } from "@mantine/core"

type Props = {
  result: TaskResultType
}

export const TaskResult = ({ result }: Props) => {
  return (
    <Badge color={result === TaskResultType.Success ? "green" : "red"}>{result}</Badge>
  )
}
