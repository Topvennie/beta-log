import { cn } from "@/lib/utils";
import { Card, CardProps, Stack } from "@mantine/core";
import { ReactNode } from "react";

type Props = {
  title: ReactNode | string;
  stat: ReactNode | string | number;
  description?: ReactNode | string;
} & CardProps

export const Stat = ({ title, stat, description, className, ...props }: Props) => {
  return (
    <Card className={cn("border border-neutral-200", className)} {...props}>
      <Stack>
        <div className="text-neutral-400">{title}</div>
        <div className="font-bold text-2xl">{stat}</div>
        <div className="text-neutral-400">{description} </div>
      </Stack>
    </Card>
  )
}
