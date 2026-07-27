import { QueryReponse } from "@/lib/api/query";
import { Gym, GymCreate, gymCreateSchema, GymUpdate, gymUpdateSchema, convertGymUpdateSchema } from "@/lib/types/gym";
import { Button, Stack, TextInput } from "@mantine/core"
import { useForm, UseFormReturnType } from "@mantine/form";
import { zod4Resolver } from "mantine-form-zod-resolver";
import { useState } from "react";

type CreateProps = {
  gym?: undefined;
  onSubmit: (gym: GymCreate) => Promise<QueryReponse<Gym>>;
}

type UpdateProps = {
  gym: Gym;
  onSubmit: (gym: GymUpdate) => Promise<QueryReponse<Gym>>;
}

type Props = CreateProps | UpdateProps;

export const GymForm = ({ gym, onSubmit }: Props) => {
  if (gym) return <Update gym={gym} onSubmit={onSubmit} />
  return <Create onSubmit={onSubmit} />
}

const Create = (props: CreateProps) => {
  const form = useForm<GymCreate>({
    initialValues: {
      name: "",
      iconPath: undefined,
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
}

const Form = <T extends GymCreate | GymUpdate>({ form, onSubmit }: FormProps<T>) => {
  const [submitting, setSubmitting] = useState(false)

  const handleSubmit = () => {
    const hasErrors = form.validate().hasErrors

    const values = form.getValues()

    if (hasErrors) return

    setSubmitting(true)
    onSubmit(values).finally(() => setSubmitting(false))
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
      <Button onClick={handleSubmit} loading={submitting} className="ml-auto">
        Submit
      </Button>
    </Stack>
  )
}
