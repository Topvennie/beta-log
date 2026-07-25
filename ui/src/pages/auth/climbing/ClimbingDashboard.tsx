import { LinkButton } from "@/components/atoms/LinkButton";
import { ClimbStat } from "@/components/climb/ClimbStat";
import { LoadingLayout } from "@/layout/LoadingLayout";
import { useClimbGetStats } from "@/lib/api/climb";
import { useHeaderContent } from "@/lib/hooks/useHeaderContent";
import { ClimbStats } from "@/lib/types/climb";
import { BarChart, BarChartSeries, ChartTooltip, CompositeChart } from '@mantine/charts';
import { getThemeColor, Group, Stack, useMantineTheme } from "@mantine/core";
import { LuDatabase } from "react-icons/lu";
import { BarShapeProps, Rectangle } from "recharts";

export const ClimbingDashboard = () => {
  useHeaderContent(<HeaderContent />)

  const { data: stats, isLoading } = useClimbGetStats()

  return (
    <LoadingLayout isLoading={isLoading}>
      <div className="grid grid-cols-4 gap-4">
        <ClimbStat
          title="Total Climbs"
          stat={stats?.total}
          description={`${stats?.totalUnique} Unique`}
        />
        <ClimbStat
          title="Top Grade"
          stat={stats?.best}
          description={`${stats?.bestAmount} Times`}
        />
        <ClimbStat
          title="Top Flash"
          stat={stats?.bestFlash}
          description={`${stats?.bestFlashAmount} Times`}
        />
        <ClimbStat
          title="Sessions"
          stat={stats?.sessions}
          description={`Med. ${stats?.medianClimbsPerSession} Climbs / Session`}
        />
        <div className="col-span-3 row-span-2">
          <ClimbStat
            title="Grade Progress (Top per Session)"
            stat={<ClimbGraph graphProgress={stats?.graphProgress ?? []} />}
          />
        </div>
        <ClimbStat
          title="Finish Types"
          stat={
            <div className="flex h-4 overflow-hidden rounded bg-gray-200">
              <div
                className="h-4 bg-blue-700"
                style={{ width: `${Math.round((stats?.flash ?? 0) / (stats?.total ?? 0) * 100)}%` }}
              />
              <div
                className="h-4 bg-blue-500"
                style={{ width: `${Math.round((stats?.top ?? 0) / (stats?.total ?? 0) * 100)}%` }}
              />
              <div
                className="h-4 bg-blue-300"
                style={{ width: `${Math.round((stats?.repeat ?? 0) / (stats?.total ?? 0) * 100)}%` }}
              />
            </div>
          }
          description={
            <Stack gap={2}>
              <Group>
                <div className="w-2 h-2 rounded-full bg-blue-700" />
                <p className="text-black mr-auto">Flash</p>
                <p>{stats?.flash}</p>
              </Group>
              <Group>
                <div className="w-2 h-2 rounded-full bg-blue-500" />
                <p className="text-black mr-auto">Top</p>
                <p>{stats?.top}</p>
              </Group>
              <Group>
                <div className="w-2 h-2 rounded-full bg-blue-300" />
                <p className="text-black mr-auto">Repeat</p>
                <p>{stats?.repeat}</p>
              </Group>
            </Stack>
          }
        />
        <div className="col-span-4 row-span-2">
          <ClimbStat
            title="Climbs per Grade"
            stat={<ClimbBar graphPerGrade={stats?.graphPerGrade ?? []} />}
          />
        </div>
      </div>
    </LoadingLayout>
  )
}

const ClimbGraph = ({ graphProgress }: Pick<ClimbStats, "graphProgress">) => {
  return (
    <CompositeChart
      h={250}
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
      textColor="gray.4"
      withLegend
      legendProps={{ verticalAlign: "bottom" }}
    />
  )
}

const ClimbBar = ({ graphPerGrade }: Pick<ClimbStats, "graphPerGrade">) => {
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
  }

  const finishShade: Record<string, number> = {
    flash: 7,
    top: 5,
    repeat: 3,
  }

  const getColor = (grade: number, series: BarChartSeries) => {
    const hue = gradeHue[grade] ?? "gray"
    const shade = finishShade[series.name] ?? 5
    return `${hue}.${shade}`
  }

  const theme = useMantineTheme()

  const series = [
    { name: "flash", label: "Flash", color: `blue.${finishShade.flash}` },
    { name: "top", label: "Top", color: `blue.${finishShade.top}` },
    { name: "repeat", label: "Repeat", color: `blue.${finishShade.repeat}` },
  ] as BarChartSeries[]

  const seriesByName = Object.fromEntries(series.map(s => [s.name, s]))

  return (
    <BarChart
      h={250}
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
      textColor="gray.4"
    />
  )
}

const HeaderContent = () => {
  return (
    <LinkButton variant="outline" to="/climbing/data" leftSection={<LuDatabase />} size="xs">
      Manage Data
    </LinkButton>
  )
}
