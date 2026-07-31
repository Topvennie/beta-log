import { useSettingUpdateGradeSystem } from "@/lib/api/setting";
import { convertSettingUpdateGradeSystem, GradeSystem, Setting, settingUpdateGradeSystem } from "@/lib/types/setting";
import { capitalize } from "@/lib/utils";
import { Button, Select, Stack, Title } from "@mantine/core";
import { useForm } from "@mantine/form";
import { notifications } from "@mantine/notifications";
import { zod4Resolver } from "mantine-form-zod-resolver";
import { useState } from "react";

type Props = {
  setting: Setting;
}

export const SettingGradeSystem = ({ setting }: Props) => {
  const settingUpdate = useSettingUpdateGradeSystem()

  const [submitting, setSubmitting] = useState(false)

  const form = useForm({
    initialValues: convertSettingUpdateGradeSystem(setting),
    validate: zod4Resolver(settingUpdateGradeSystem),
  })

  const handleSubmit = () => {
    if (form.validate().hasErrors) return

    setSubmitting(true)

    settingUpdate.mutate(form.getValues(), {
      onSuccess: () => {
        notifications.show({ color: "green", title: "Setting", message: "Grade system saved" })
      },
      onSettled: () => setSubmitting(false)
    })
  }

  return (
    <div className="flex md:flex-row gap-8 h-full">
      <Stack className="flex-1">
        <Title order={2}>Grade System</Title>
        <p className="text-pretty whitespace-pre-wrap">{`Select which grade system to use.`}</p>
      </Stack>

      <div className="border-l border-neutral-200" />

      <Stack className="flex-2">
        <Select
          data={Object.values(GradeSystem).map(g => ({ value: g, label: capitalize(g) }))}
          value={form.values.gradeSystem}
          onChange={value => form.setFieldValue("gradeSystem", value ? value as GradeSystem : GradeSystem.Font)}
          error={form.errors.gradeSystem}
        />
        <Button onClick={handleSubmit} loading={submitting}>Save</Button>
      </Stack>
    </div>
  )
}
