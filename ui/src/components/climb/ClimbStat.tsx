import { Card, Stack } from "@mantine/core";
import { ReactNode } from "react";

type Props = {
  title: ReactNode | string;
  stat: ReactNode | string | number;
  description?: ReactNode | string;
}

export const ClimbStat = ({ title, stat, description }: Props) => {
  return (
    <Card className="border border-gray-200">
      <Stack>
        <div className="text-neutral-400">{title}</div>
        <div className="font-bold text-2xl">{stat}</div>
        <div className="text-neutral-400">{description} </div>
      </Stack>
    </Card>
  )
}
