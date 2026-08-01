import { SettingGradeSystem } from "@/components/setting/SettingGradeSystem"
import { SettingToplogger } from "@/components/setting/SettingToplogger"
import { LoadingLayout } from "@/layout/LoadingLayout"
import { useSettingGet } from "@/lib/api/setting"
import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb"
import { Stack } from "@mantine/core"

export const Settings = () => {
  useBreadcrumb({ title: "Settings", weight: 10, link: { to: "/settings" } })

  const { data: setting, isLoading } = useSettingGet()

  return (
    <LoadingLayout isLoading={isLoading || !setting}>
      <Stack gap="xl">
        <SettingGradeSystem setting={setting!} />
        <SettingToplogger setting={setting!} />
      </Stack>
    </LoadingLayout>
  )
}
