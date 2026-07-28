import { QueryReponse } from "@/lib/api/query";
import { Gym, GymCreate, gymCreateSchema, GymUpdate, gymUpdateSchema, convertGymUpdateSchema } from "@/lib/types/gym";
import { ActionIcon, Button, Group, Stack, TextInput } from "@mantine/core"
import { useForm, UseFormReturnType } from "@mantine/form";
import { zod4Resolver } from "mantine-form-zod-resolver";
import { useState } from "react";
import { FaTrashCan } from "react-icons/fa6";
import { Confirm } from "../atoms/Confirm";
import { useDisclosure } from "@mantine/hooks";

type CreateProps = {
  gym?: undefined;
  onSubmit: (gym: GymCreate) => Promise<QueryReponse<Gym>>;
  onCancel: () => void;
  onDelete?: undefined;
}

type UpdateProps = {
  gym: Gym;
  onSubmit: (gym: GymUpdate) => Promise<QueryReponse<Gym>>;
  onCancel: () => void;
  onDelete: ({ id }: Pick<Gym, "id">) => Promise<QueryReponse<null>>;
}

type Props = CreateProps | UpdateProps;

export const GymForm = ({ gym, onSubmit, onCancel, onDelete }: Props) => {
  if (gym) return <Update gym={gym} onSubmit={onSubmit} onCancel={onCancel} onDelete={onDelete} />
  return <Create onSubmit={onSubmit} onCancel={onCancel} />
}

const Create = (props: CreateProps) => {
  const form = useForm<GymCreate>({
    initialValues: {
      name: "",
      iconPath: "",
    },
    validate: zod4Resolver(gymCreateSchema)
  })

  return <Form form={form} {...props} />
}

const Update = ({ gym, ...props }: UpdateProps) => {
  const form = useForm<GymUpdate>({
    initialValues: convertGymUpdateSchema(gym),
    validate: zod4Resolver(gymUpdateSchema),
  })

  return <Form form={form} {...props} />
}

type FormProps<T extends GymCreate | GymUpdate> = {
  form: UseFormReturnType<T>;
  onSubmit: (t: T) => Promise<QueryReponse<Gym>>;
  onCancel: () => void
  onDelete?: ({ id }: Pick<Gym, "id">) => Promise<QueryReponse<null>>
}

const Form = <T extends GymCreate | GymUpdate>({ form, onSubmit, onCancel, onDelete }: FormProps<T>) => {
  const [opened, { open, close }] = useDisclosure()
  const [submitting, setSubmitting] = useState(false)

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
    onDelete?.(values as GymUpdate)
      .then(() => form.reset())
      .finally(() => setSubmitting(false))
  }

  return (
    <Stack>
      <TextInput
        label="Name"
        placeholder="B-PUMP"
        required
        {...form.getInputProps("name")}
      />
      <TextInput
        label="Icon URL"
        placeholder="https://b-pump_logo.png"
        {...form.getInputProps("iconPath")}
      />
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
