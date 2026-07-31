import { LinkButton } from "@/components/atoms/LinkButton";
import { DatePreset, SelectDatePreset } from "@/components/atoms/SelectDatePreset";
import { Stat } from "@/components/atoms/Stat";
import { ProportionBar } from "@/components/molecules/ProportionBar";
import { LoadingLayout } from "@/layout/LoadingLayout";
import { useClimbStatGetFiltered } from "@/lib/api/climb";
import { useHeaderContent } from "@/lib/hooks/useHeaderContent";
import { ClimbStats } from "@/lib/types/climb";
import { BarChart, BarChartSeries, ChartTooltip, CompositeChart } from '@mantine/charts';
import { getThemeColor, Group, useMantineTheme } from "@mantine/core";
import { subMonths, subYears } from "date-fns";
import { useState } from "react";
import { LuDatabase } from "react-icons/lu";
import { BarShapeProps, Rectangle } from "recharts";

export const ClimbingDashboard = () => {
  const [dates, setDates] = useState<[Date | null, Date | null]>([null, null])
  useHeaderContent(<HeaderContent setValue={setDates} />)

  const { data: stats, isLoading } = useClimbStatGetFiltered(dates[0] ?? undefined, dates[1] ?? undefined)

  return (
    <LoadingLayout isLoading={isLoading}>
      <div className="grid grid-cols-4 gap-4">
        <Stat
          title="Total Climbs"
          stat={stats?.total}
          description={`${stats?.totalUnique} Unique`}
        />
        <Stat
          title="Top Grade"
          stat={stats?.best}
          description={`${stats?.bestAmount} Times`}
        />
        <Stat
          title="Top Flash"
          stat={stats?.bestFlash}
          description={`${stats?.bestFlashAmount} Times`}
        />
        <Stat
          title="Sessions"
          stat={stats?.sessions}
          description={`Med. ${stats?.medianClimbsPerSession} Climbs / Session`}
        />
        <div className="col-span-3 row-span-2">
          <Stat
            title="Grade Progress (Top per Session)"
            stat={<GraphProgress graphProgress={stats?.graphProgress ?? []} />}
          />
        </div>
        <Stat
          title="Finish Types"
          stat={<GraphFinishType total={stats?.total ?? 1} flash={stats?.flash ?? 0} top={stats?.top ?? 0} repeat={stats?.repeat ?? 0} />}
        />
        <Stat
          title="Climb Types"
          stat={<GraphClimbType boulder={stats?.boulder ?? 1} lead={stats?.lead ?? 1} />}
        />
        <div className="col-span-4 row-span-2">
          <Stat
            title="Climbs per Grade"
            stat={<GraphGrade graphPerGrade={stats?.graphPerGrade ?? []} />}
          />
        </div>
      </div>
    </LoadingLayout>
  )
}

const GraphFinishType = ({ total, flash, top, repeat }: Pick<ClimbStats, "total" | "flash" | "top" | "repeat">) => {
  return (
    <ProportionBar
      items={[
        { label: "Flash", color: "blue-700", amount: flash },
        { label: "Top", color: "blue-500", amount: top },
        { label: "Repeat", color: "blue-300", amount: repeat },
      ]}
      total={total}
      withDescription
    />
  )
}

const GraphClimbType = ({ boulder, lead }: Pick<ClimbStats, "boulder" | "lead">) => {
  return (
    <ProportionBar
      items={[
        { label: "Boulder", color: "blue-500", amount: boulder },
        { label: "Lead", color: "blue-300", amount: lead },
      ]}
      total={boulder + lead}
      withDescription
    />
  )
}

const GraphProgress = ({ graphProgress }: Pick<ClimbStats, "graphProgress">) => {
  return (
    <CompositeChart
      data={graphProgress}
      dataKey="date"
      maxBarWidth={30}
      series={[
        { name: "grade", label: "Grade", color: "blue.7", type: "line" },
        { name: "volume", label: "Volume", color: "rgba(18, 129, 255, 0.2)", type: "bar", yAxisId: "right" },
      ]}
      tickLine="none"
      withXAxis={false}
      withRightYAxis
      withLegend
      legendProps={{ verticalAlign: "bottom" }}
    />
  )
}

const GraphGrade = ({ graphPerGrade }: Pick<ClimbStats, "graphPerGrade">) => {
  const gradeHue: Record<number, string> = {
    0: "green",
    200: "green",
    250: "green",
    300: "green",
    333: "green",
    367: "green",
    400: "yellow",
    433: "yellow",
    467: "yellow",
    500: "yellow",
    517: "orange",
    533: "orange",
    550: "orange",
    567: "orange",
    583: "orange",
    600: "orange",
    617: "orange",
    633: "blue",
    650: "blue",
    667: "red",
  }

  const finishShade: Record<string, number> = {
    flash: 7,
    top: 5,
    repeat: 3,
  }

  const getColor = (grade: number, series: BarChartSeries) => {
    const hue = gradeHue[grade] ?? "neutral"
    const shade = finishShade[series.name] ?? 5
    return `${hue}.${shade}`
  }

  const theme = useMantineTheme()

  const series = [
    { name: "flash", label: "Flash", color: `green.${finishShade.flash}` },
    { name: "top", label: "Top", color: `green.${finishShade.top}` },
    { name: "repeat", label: "Repeat", color: `green.${finishShade.repeat}` },
  ] satisfies BarChartSeries[]

  const seriesByName = Object.fromEntries(series.map(s => [s.name, s]))

  return (
    <BarChart
      data={graphPerGrade}
      dataKey="grade"
      type="stacked"
      series={series}
      barProps={(series) => ({
        shape: (props: BarShapeProps) => {
          const grade = props.payload?.grade ?? 0
          return (
            <Rectangle
              {...props}
              fill={getThemeColor(getColor(grade, series), theme)}
            />
          )
        },
      })}
      tooltipProps={{
        content: ({ label, payload }) => (
          <ChartTooltip
            label={label}
            payload={payload?.map(item => ({
              ...item,
              color: getColor(item.payload?.grade ?? 0, seriesByName[item.name ?? ""] ?? { name: item.name ?? "" }),
            }))}
            series={series}
          />
        ),
      }}
      withLegend
      legendProps={{ verticalAlign: "bottom", itemSorter: null }}
    />
  )
}

type HeaderContentProps = {
  setValue: (value: [Date | null, Date | null]) => void
}

const HeaderContent = ({ setValue }: HeaderContentProps) => {
  const [preset, setPreset] = useState<DatePreset>("All Time")

  const onChange = (p: DatePreset) => {
    const now = new Date()

    let start: Date | null = null
    let end: Date | null = null

    switch (p) {
      case "1 Month":
        end = now
        start = subMonths(now, 1)
        break
      case "3 Months":
        end = now
        start = subMonths(now, 3)
        break
      case "1 Year":
        end = now
        start = subYears(now, 1)
        break
      case "All Time":
        break
    }

    setPreset(p)
    setValue([start, end])
  }

  return (
    <Group>
      <SelectDatePreset value={preset} setValue={onChange} />
      <LinkButton variant="outline" to="/climbing/data" leftSection={<LuDatabase />} size="xs">
        Manage Data
      </LinkButton>
    </Group>
  )
}
