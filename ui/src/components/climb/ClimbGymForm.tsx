import { QueryReponse } from "@/lib/api/query";
import { ClimbGym, ClimbGymCreate, climbGymCreateSchema, ClimbGymUpdate, convertClimbGymUpdateSchema } from "@/lib/types/climb";
import { Button, Stack, TextInput } from "@mantine/core"
import { useForm, UseFormReturnType } from "@mantine/form";
import { zod4Resolver } from "mantine-form-zod-resolver";
import { useState } from "react";

type CreateProps = {
  gym?: undefined;
  onSubmit: (gym: ClimbGymCreate) => Promise<QueryReponse<ClimbGym>>;
}

type UpdateProps = {
  gym: ClimbGym;
  onSubmit: (gym: ClimbGymUpdate) => Promise<QueryReponse<ClimbGym>>;
}

type Props = CreateProps | UpdateProps;

export const ClimbGymForm = ({ gym, onSubmit }: Props) => {
  if (gym) return <Update gym={gym} onSubmit={onSubmit} />
  return <Create onSubmit={onSubmit} />
}

const Create = (props: CreateProps) => {
  const form = useForm<ClimbGymCreate>({
    initialValues: {
      name: "",
      iconPath: undefined,
    },
    validate: zod4Resolver(climbGymCreateSchema)
  })

  return <Form form={form} {...props} />
}

const Update = ({ gym, ...props }: UpdateProps) => {
  const form = useForm<ClimbGymUpdate>({
    initialValues: convertClimbGymUpdateSchema(gym),
    validate: zod4Resolver(climbGymCreateSchema),
  })

  return <Form form={form} {...props} />
}

type FormProps<T extends ClimbGymCreate | ClimbGymUpdate> = {
  form: UseFormReturnType<T>;
  onSubmit: (t: T) => Promise<QueryReponse<ClimbGym>>;
}

const Form = <T extends ClimbGymCreate | ClimbGymUpdate>({ form, onSubmit }: FormProps<T>) => {
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
