/* eslint-disable @typescript-eslint/no-explicit-any */

import { QueryReponse } from "@/lib/api/query";
import { convertExerciseUpdateSchema, Exercise, ExerciseCreate, exerciseCreateSchema, ExerciseUpdate, exerciseUpdateSchema } from "@/lib/types/exercise";
import { ActionIcon, Button, Group, Stack, TextInput } from "@mantine/core";
import { useForm, UseFormReturnType } from "@mantine/form";
import { useDisclosure } from "@mantine/hooks";
import { zod4Resolver } from "mantine-form-zod-resolver";
import { useState } from "react";
import { FaTrashCan } from "react-icons/fa6";
import { Confirm } from "../atoms/Confirm";

type CreateProps = {
  exercise?: undefined;
  onSubmit: (exercise: ExerciseCreate) => Promise<QueryReponse<Exercise>>;
  onCancel: () => void;
  onDelete?: undefined;
}

type UpdateProps = {
  exercise: Exercise;
  onSubmit: (exercise: ExerciseUpdate) => Promise<QueryReponse<Exercise>>;
  onCancel: () => void;
  onDelete: ({ id }: Pick<Exercise, "id">) => Promise<QueryReponse<null>>;
}

type Props = CreateProps | UpdateProps

export const ExerciseForm = ({ exercise, onSubmit, onCancel, onDelete }: Props) => {
  if (exercise) return <ExerciseFormUpdate exercise={exercise} onSubmit={onSubmit} onCancel={onCancel} onDelete={onDelete} />
  return <ExerciseFormCreate onSubmit={onSubmit} onCancel={onCancel} />
}

const ExerciseFormCreate = (props: CreateProps) => {
  const form = useForm<ExerciseCreate>({
    initialValues: {
      name: "",
      variants: [],
    },
    validate: zod4Resolver(exerciseCreateSchema)
  })

  return <ExerciseFormInner form={form} {...props} />
}

const ExerciseFormUpdate = ({ exercise, ...props }: UpdateProps) => {
  const form = useForm<ExerciseUpdate>({
    initialValues: convertExerciseUpdateSchema(exercise),
    validate: zod4Resolver(exerciseUpdateSchema),
  })

  return <ExerciseFormInner form={form} {...props} />
}

type ExerciseFormInnerProps<T extends ExerciseCreate | ExerciseUpdate> = {
  form: UseFormReturnType<T>
  onSubmit: (exercise: T) => Promise<QueryReponse<Exercise>>
  onCancel: () => void
  onDelete?: ({ id }: Pick<Exercise, "id">) => Promise<QueryReponse<null>>
}

const ExerciseFormInner = <T extends ExerciseCreate | ExerciseUpdate>({ form, onSubmit, onCancel, onDelete }: ExerciseFormInnerProps<T>) => {
  const [opened, { open, close }] = useDisclosure()

  const [submitting, setSubmitting] = useState(false)

  const handleSubmit = () => {
    if (form.validate().hasErrors) return
    const values = form.getValues()

    setSubmitting(true)
    onSubmit(values)
      .then(() => form.reset())
      .finally(() => setSubmitting(false))
  }

  const handleCancel = () => {
    form.reset()
    onCancel()
  }

  const handleDeleteInit = () => {
    open()
  }

  const handleDelete = () => {
    const values = form.getValues()
    if (!("id" in values)) return

    setSubmitting(true)
    onDelete?.(values as ExerciseUpdate)
      .then(() => form.reset())
      .finally(() => setSubmitting(false))
  }

  const handleAddVariant = () => {
    const current = form.getValues().variants
    form.setFieldValue("variants", [...current, { variant: "" }] as any)
  }

  const handleDeleteVariant = (idx: number) => {
    const current = [...form.getValues().variants] as Array<{ id?: number; variant: string }>
    current.splice(idx, 1)
    form.setFieldValue("variants", current as any)
  }

  return (
    <>
      <Stack>
        <TextInput
          label="Name"
          placeholder="No Hang"
          required
          {...form.getInputProps("name")}
        />
        <Stack gap="xs">
          <p>Variants</p>
          {form.getValues().variants.map((_, idx) => (
            <Group key={idx} gap="xs">
              <TextInput
                placeholder="Right hand"
                {...form.getInputProps(`variants.${idx}.variant`)}
                className="grow"
              />
              <ActionIcon variant="subtle" color="red" onClick={() => handleDeleteVariant(idx)}>
                <FaTrashCan />
              </ActionIcon>
            </Group>
          ))}
          <div onClick={handleAddVariant} className="cursor-pointer rounded-sm border-4 border-dotted border-gray-200 py-2 px-4">
            <p className="text-center text-neutral-400">Add Variant</p>
          </div>
        </Stack>
        <Button onClick={handleSubmit} loading={submitting}>Submit</Button>
        <Group>
          <Button onClick={handleCancel} variant="outline" className="flex-1" loading={submitting}>Cancel</Button>
          <ActionIcon onClick={handleDeleteInit} variant="subtle" size="lg" color="red" loading={submitting} disabled={!onDelete}>
            <FaTrashCan />
          </ActionIcon>
        </Group>
      </Stack>
      <Confirm
        opened={opened}
        onClose={close}
        onConfirm={handleDelete}
        onAbort={close}
      />
    </>
  )
}
