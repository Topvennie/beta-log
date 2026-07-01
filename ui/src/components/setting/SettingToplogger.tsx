import { useSettingToploggerUpdate } from "@/lib/api/setting";
import { convertSettingToploggerUpdateSchema, Setting, settingToploggerUpdateSchema } from "@/lib/types/setting";
import { Alert, Button, PasswordInput, Stack, TextInput, Title } from "@mantine/core";
import { useForm } from "@mantine/form";
import { notifications } from "@mantine/notifications";
import { zod4Resolver } from "mantine-form-zod-resolver";
import { useState } from "react";
import { FaCircleInfo, FaTriangleExclamation } from "react-icons/fa6";

type Props = {
  setting: Setting;
}

export const SettingToplogger = ({ setting }: Props) => {
  const settingUpdate = useSettingToploggerUpdate()

  const [disabled, setDisabled] = useState(false)
  const [submitting, setSubmitting] = useState(false)

  const form = useForm({
    initialValues: convertSettingToploggerUpdateSchema(setting),
    validate: zod4Resolver(settingToploggerUpdateSchema),
  })

  const handleSubmit = () => {
    if (form.validate().hasErrors) return

    setSubmitting(true)

    settingUpdate.mutate(form.getValues(), {
      onSuccess: () => {
        notifications.show({ color: "green", title: "Setting", message: "Toplogger settings saved" })
        setDisabled(true)
      },
      onSettled: () => setSubmitting(false)
    })
  }

  return (
    <div className="flex md:flex-row gap-8 h-full">
      <Stack className="flex-1">
        <Title order={2}>Toplogger Integration</Title>
        <p className="text-pretty whitespace-pre-wrap">{`Automatically import your climbs. This method will log you out of the web app (NOT the mobile app). If you ever log back in on the webapp then you will have to enter this information again.`}</p>
      </Stack>

      <div className="border-l border-gray-200" />

      <Stack className="flex-2">
        <TextInput
          label="User id"
          {...form.getInputProps("climbToploggerUserId")}
        />
        <PasswordInput
          label="Auth token"
          {...form.getInputProps("climbToploggerAuthToken")}
        />
        <PasswordInput
          label="Refresh token"
          {...form.getInputProps("climbToploggerRefreshToken")}
        />
        <Button onClick={handleSubmit} loading={submitting} disabled={disabled}>Save</Button>
        {disabled && (
          <Alert variant="light" color="red" icon={<FaTriangleExclamation />} title="Disabled">
            <p>To prevent unsynced updates, this button is disabled until a refresh</p>
            <p>Everything is working as intended, no need to worry :)</p>
          </Alert>
        )}
        <Alert variant="light" icon={<FaCircleInfo />} title="Get your data" className="text-pretty">
          <ol className="list-decimal pl-5">
            <li>Go to <a href="https://app.toplogger.nu" target="_blank" rel="noreferrer noopener" className="underline">app.toplogger.nu</a> and log in</li>
            <li>Open the dev tools (right click on the screen -&gt; inspect)</li>
            <li>Go to the storage tab</li>
            <li>Go to local storage</li>
            <li>Open the key <code className="px-1 bg-gray-100 rounded">tl-auth</code>, you can find your access and refresh token in the object</li>
            <li>Open the key <code className="px-1 bg-gray-100 rounded">tl-user-states</code>, you can find your user_id underneath the version</li>
          </ol>
        </Alert>
      </Stack>
    </div>
  )
}
