import { Group, Stack } from "@mantine/core"

type Props = {
  items: Item[];
  total: number;
  withDescription?: boolean
}

type Item = {
  label?: string;
  color: string;
  amount: number;
}

const mantineColor = (color: string) => {
  const [name, shade] = color.split("-")
  return `var(--mantine-color-${name}-${Number(shade) / 100})`
}

export const ProportionBar = ({ items, total, withDescription = false }: Props) => {
  return (
    <Stack>
      <div className="flex h-4 overflow-hidden rounded bg-neutral-200">
        {items.map(i => (
          <div
            key={i.color}
            className="h-4"
            style={{
              width: `${Math.round(i.amount / total * 100)}%`,
              backgroundColor: mantineColor(i.color),
            }}
          />
        ))}
      </div>
      {withDescription && (
        <Stack gap={2} className="text-neutral-400 text-md font-normal">
          {items.map(i => (
            <Group key={i.color}>
              <div
                className="w-2 h-2 rounded-full"
                style={{ backgroundColor: mantineColor(i.color) }}
              />
              <p className="text-black mr-auto">{i.label}</p>
              <p>{i.amount}</p>
            </Group>
          ))}
        </Stack>
      )}
    </Stack>
  )
}
