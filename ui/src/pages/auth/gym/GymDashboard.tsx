import { LinkButton } from "@/components/atoms/LinkButton"
import { Stat } from "@/components/atoms/Stat"
import { LoadingLayout } from "@/layout/LoadingLayout"
import { useGymGetStats } from "@/lib/api/gym"
import { useHeaderContent } from "@/lib/hooks/useHeaderContent"
import { GymStats } from "@/lib/types/gym"
import { BarChart, BarChartSeries } from "@mantine/charts"
import { LuDatabase } from "react-icons/lu"

export const GymDashboard = () => {
  useHeaderContent(<HeaderContent />)

  const { data: stats, isLoading } = useGymGetStats()

  return (
    <LoadingLayout isLoading={isLoading}>
      <div className="grid grid-cols-6 gap-4">
        <Stat
          title="Most Visited"
          stat={stats?.mostVisited}
          description={`${stats?.mostVisitedAmount} visits`}
          className="col-span-2"
        />
        <Stat
          title="Total Gyms"
          stat={stats?.total}
          className="col-span-2"
        />
        <Stat
          title="Total Sessions"
          stat={stats?.sessions}
          className="col-span-2"
        />
        <Stat
          title="Visits per Gym"
          stat={<GraphVisits graphVisits={stats?.graphVisits ?? []} />}
          className="col-span-3"
        />
        <Stat
          title="Best Top & Flash per Gym"
          stat={<GraphTop graphTop={stats?.graphTop ?? []} />}
          className="col-span-3"
        />
        <Stat
          title="Grade Distribution per Gym"
          stat={<GraphDistribution graphDistribution={stats?.graphDistribution ?? []} />}
          className="col-span-6"
        />
      </div>
    </LoadingLayout>
  )
}

const GraphVisits = ({ graphVisits }: Pick<GymStats, "graphVisits">) => {
  return (
    <BarChart
      data={graphVisits}
      dataKey="gym"
      series={[
        { name: "amount", label: "Visits", color: "blue" },
      ]}
      gridAxis="none"
      withYAxis={false}
      withBarValueLabel
      barProps={{ radius: 8 }}
      valueFormatter={(value) => value.toString()}
    />
  )
}

const GraphTop = ({ graphTop }: Pick<GymStats, "graphTop">) => {
  return (
    <BarChart
      data={graphTop}
      dataKey="gym"
      series={[
        { name: "top", label: "Top", color: "blue" },
        { name: "flash", label: "Flash", color: "blue.3" },
      ]}
      gridAxis="none"
      withYAxis={false}
      withBarValueLabel
      barProps={{ radius: 8 }}
      valueFormatter={(value) => value.toString()}
    />
  )
}

const GraphDistribution = ({ graphDistribution }: Pick<GymStats, "graphDistribution">) => {
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

  const grades = Array.from(
    new Set(graphDistribution.flatMap(g => Object.keys(g.distribution).map(Number))),
  ).sort((a, b) => a - b)

  const series: BarChartSeries[] = grades.map(grade => ({
    name: String(grade),
    color: `${gradeHue[grade] ?? "neutral"}.6`,
  }))

  const data = graphDistribution.map(g => {
    const entry: Record<string, number | string | null> = { gym: g.gym }
    grades.forEach(grade => {
      entry[String(grade)] = g.distribution[grade] ?? null
    })
    return entry
  })

  return (
    <BarChart
      data={data}
      dataKey="gym"
      type="stacked"
      orientation="vertical"
      series={series}
      gridAxis="none"
      withYAxis={true}
      withXAxis={false}
      withBarValueLabel
      barProps={{}}
      valueFormatter={(v) => (v ? `${v}%` : "")}
    />
  )
}

const HeaderContent = () => {
  return (
    <LinkButton variant="outline" to="/gym/data" leftSection={<LuDatabase />} size="xs">
      Manage Data
    </LinkButton>
  )
}
