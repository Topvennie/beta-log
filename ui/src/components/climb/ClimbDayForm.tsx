/* eslint-disable @typescript-eslint/no-explicit-any */

import { useGymGetAll } from "@/lib/api/gym";
import { QueryReponse } from "@/lib/api/query";
import { ClimbDay, ClimbDayCreate, climbDayCreateSchema, ClimbDayUpdate, climbDayUpdateSchema, ClimbFinish, ClimbType, convertClimbDayUpdateSchema } from "@/lib/types/climb";
import { ActionIcon, Button, Card, ColorInput, Group, NumberInput, ScrollArea, Select, Stack } from "@mantine/core";
import { DatePickerInput } from '@mantine/dates';
import { useForm, UseFormReturnType } from "@mantine/form";
import { useDisclosure } from "@mantine/hooks";
import { zod4Resolver } from "mantine-form-zod-resolver";
import { useState } from "react";
import { FaPlus, FaTrashCan } from "react-icons/fa6";
import { Confirm } from "../atoms/Confirm";

type CreateProps = {
  climbDay?: undefined;
  onSubmit: (climbDay: ClimbDayCreate) => Promise<QueryReponse<ClimbDay>>;
  onCancel: () => void;
  onDelete?: undefined;
}

type UpdateProps = {
  climbDay: ClimbDay;
  onSubmit: (climbDay: ClimbDayUpdate) => Promise<QueryReponse<ClimbDay>>;
  onCancel: () => void;
  onDelete: ({ id }: Pick<ClimbDay, "id">) => Promise<QueryReponse<null>>;
}

type Props = CreateProps | UpdateProps;

export const ClimbDayForm = ({ climbDay, onSubmit, onCancel, onDelete }: Props) => {
  if (climbDay) return <Update climbDay={climbDay} onSubmit={onSubmit} onCancel={onCancel} onDelete={onDelete} />
  return <Create onSubmit={onSubmit} onCancel={onCancel} />
}

const Create = (props: CreateProps) => {
  const form = useForm<ClimbDayCreate>({
    initialValues: {
      date: new Date(),
      gymId: 0,
      climbs: [],
    },
    validate: zod4Resolver(climbDayCreateSchema)
  })

  return <Form form={form} {...props} />
}

const Update = ({ climbDay, ...props }: UpdateProps) => {
  const form = useForm<ClimbDayUpdate>({
    initialValues: convertClimbDayUpdateSchema(climbDay),
    validate: zod4Resolver(climbDayUpdateSchema),
  })

  return <Form form={form} {...props} />
}

type FormProps<T extends ClimbDayCreate | ClimbDayUpdate> = {
  form: UseFormReturnType<T>;
  onSubmit: (t: T) => Promise<QueryReponse<ClimbDay>>;
  onCancel: () => void
  onDelete?: ({ id }: Pick<ClimbDay, "id">) => Promise<QueryReponse<null>>
}

const Form = <T extends ClimbDayCreate | ClimbDayUpdate>({ form, onSubmit, onCancel, onDelete }: FormProps<T>) => {
  const { data: gyms, isLoading } = useGymGetAll()

  const [opened, { open, close }] = useDisclosure()
  const [submitting, setSubmitting] = useState(false)

  const handleClimbAdd = () => {
    const climbs = form.getValues().climbs
    form.setFieldValue("climbs", [...climbs, {
      _clientId: crypto.randomUUID(),
      grade: 0,
      holdColor: "",
      climbType: "",
      finishType: "",
    }] as any)
  }

  const handleClimbDelete = (clientId: string) => {
    const climbs = form.getValues().climbs.filter(c => c._clientId !== clientId)
    form.setFieldValue("climbs", climbs as any)
  }

  const handleSubmit = () => {
    const hasErrors = form.validate().hasErrors

    const values = form.getValues()

    if (hasErrors) return

    setSubmitting(true)
    onSubmit(values).finally(() => setSubmitting(false))
  }

  const handleCancel = () => {
    form.reset()
    onCancel()
  }

  const handleDelete = () => {
    const values = form.getValues()
    if (!("id" in values)) return

    setSubmitting(true)
    onDelete?.(values as ClimbDayUpdate)
      .then(() => form.reset())
      .finally(() => setSubmitting(false))
  }

  return (
    <Stack>
      <DatePickerInput
        label="Date"
        required
        value={form.values.date}
        onChange={value => form.setFieldValue("date", new Date(value ?? "") as any)}
        error={form.errors.date}
      />
      <Select
        label="Gym"
        required
        data={gyms?.map(g => ({ value: g.id.toString(), label: g.name }))}
        value={form.values.gymId.toString()}
        onChange={value => form.setFieldValue("gymId", Number(value) as any)}
        error={form.errors.gymId}
        loading={isLoading}
      />
      <Group justify="space-between">
        <Stack>
          <p>{`Climbs (${form.values.climbs.length})`}</p>
          {form.errors.climbs && <p className="text-sm text-red-500">{form.errors.climbs}</p>}
        </Stack>
        <Button onClick={handleClimbAdd} variant="subtle" leftSection={<FaPlus />}>
          Add Climb
        </Button>
      </Group>
      <Card className="border border-neutral-200">
        <ScrollArea h={250} offsetScrollbars>
          <Stack>
            {form.values.climbs.map(c => <ClimbForm key={c._clientId} form={form} clientId={c._clientId} onDelete={handleClimbDelete} />)}
          </Stack>
        </ScrollArea>
      </Card>
      <Button onClick={handleSubmit} loading={submitting}>
        Submit
      </Button>
      <Group>
        <Button onClick={handleCancel} variant="outline" className="flex-1" loading={submitting}>Cancel</Button>
        <ActionIcon onClick={open} variant="subtle" size="lg" color="red" loading={submitting} disabled={!onDelete}>
          <FaTrashCan />
        </ActionIcon>
      </Group>
      <Confirm
        opened={opened}
        onClose={close}
        onConfirm={handleDelete}
        onAbort={close}
      />
    </Stack>
  )
}

type ClimbFormProps<T extends ClimbDayCreate | ClimbDayUpdate> = {
  form: UseFormReturnType<T>;
  clientId: string;
  onDelete: (clientId: string) => void;
}

const ClimbForm = <T extends ClimbDayCreate | ClimbDayUpdate>({ form, clientId, onDelete }: ClimbFormProps<T>) => {
  const idx = form.getValues().climbs.findIndex(c => c._clientId === clientId)
  const climb = form.getValues().climbs[idx]

  return (
    <Group wrap="nowrap">
      <NumberInput
        min={1}
        required
        {...form.getInputProps(`climbs.${idx}.grade`)}
      />
      <ColorInput
        required
        value={climb.holdColor}
        onChangeEnd={value => form.setFieldValue(`climbs.${idx}.holdColor`, value as any)}
        error={form.errors[`climbs.${idx}.holdColor`]}
      />
      <Select
        data={Object.values(ClimbType).map(c => ({ value: c, label: c }))}
        value={climb.climbType}
        onChange={value => form.setFieldValue(`climbs.${idx}.climbType`, value as any)}
        error={form.errors[`climbs.${idx}.climbType`]}
      />
      <Select
        data={Object.values(ClimbFinish).map(c => ({ value: c, label: c }))}
        value={climb.finishType}
        onChange={value => form.setFieldValue(`climbs.${idx}.finishType`, value as any)}
        error={form.errors[`climbs.${idx}.finishType`]}
      />
      <ActionIcon onClick={() => onDelete(clientId)}><FaTrashCan /></ActionIcon>
    </Group>
  )
}
