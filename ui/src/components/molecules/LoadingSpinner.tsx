import { Center, Group, Text } from "@mantine/core"
import { CgSpinner } from "react-icons/cg";

type Props = {
  title?: string
}

export const LoadingSpinner = ({ title = "" }: Props) => {
  return (
    <Center>
      <Group gap="sm">
        <CgSpinner className="h-8 w-8 animate-spin" />
        <Text>{title}</Text>
      </Group>
    </Center>
  )
}
