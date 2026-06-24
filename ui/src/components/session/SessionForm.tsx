/* eslint-disable @typescript-eslint/no-explicit-any */

import { useExerciseGetAll } from "@/lib/api/exercise";
import { QueryReponse } from "@/lib/api/query";
import { convertSessionUpdateSchema, Session, SessionCreate, sessionCreateSchema, SessionUpdate, sessionUpdateSchema } from "@/lib/types/session";
import { closestCenter, DndContext, PointerSensor, useSensor, useSensors } from "@dnd-kit/core";
import { SortableContext, useSortable, verticalListSortingStrategy } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { ActionIcon, Button, Group, NumberInput, Select, SimpleGrid, Stack, TextInput } from "@mantine/core";
import { useForm, UseFormReturnType } from "@mantine/form";
import { useDisclosure } from "@mantine/hooks";
import { zod4Resolver } from "mantine-form-zod-resolver";
import { useState } from "react";
import { FaGripVertical, FaTrashCan } from "react-icons/fa6";
import { Confirm } from "../atoms/Confirm";

type CreateProps = {
  session?: undefined;
  onSubmit: (session: SessionCreate) => Promise<QueryReponse<Session>>;
  onCancel: () => void;
  onDelete?: undefined;
}

type UpdateProps = {
  session: Session;
  onSubmit: (session: SessionUpdate) => Promise<QueryReponse<Session>>;
  onCancel: () => void;
  onDelete: ({ id }: Pick<Session, "id">) => Promise<QueryReponse<null>>;
}

type Props = CreateProps | UpdateProps

export const SessionForm = ({ session, onSubmit, onCancel, onDelete }: Props) => {
  if (session) return <SessionFormUpdate session={session} onSubmit={onSubmit} onCancel={onCancel} onDelete={onDelete} />
  return <SessionFormCreate onSubmit={onSubmit} onCancel={onCancel} />
}

const SessionFormCreate = (props: CreateProps) => {
  const form = useForm<SessionCreate>({
    initialValues: {
      name: "",
      exercises: [],
    },
    validate: zod4Resolver(sessionCreateSchema)
  })

  return <SessionFormInner form={form} {...props} />
}

const SessionFormUpdate = ({ session, ...props }: UpdateProps) => {
  const form = useForm<SessionUpdate>({
    initialValues: convertSessionUpdateSchema(session),
    validate: zod4Resolver(sessionUpdateSchema),
  })

  return <SessionFormInner form={form} {...props} />
}

type SessionFormInnerProps<T extends SessionCreate | SessionUpdate> = {
  form: UseFormReturnType<T>
  onSubmit: (session: T) => Promise<QueryReponse<Session>>
  onCancel: () => void
  onDelete?: ({ id }: Pick<Session, "id">) => Promise<QueryReponse<null>>
}

const SessionFormInner = <T extends SessionCreate | SessionUpdate>({ form, onSubmit, onCancel, onDelete }: SessionFormInnerProps<T>) => {
  const [opened, { open, close }] = useDisclosure()

  const [submitting, setSubmitting] = useState(false)

  const sensors = useSensors(useSensor(PointerSensor))

  const handleDragEnd = (event: any) => {
    const { active, over } = event
    if (active.id === over.id) return

    const exercises = form.getValues().exercises
    const oldIndex = exercises.findIndex((e) => e.clientId === active.id)
    const newIndex = exercises.findIndex((e) => e.clientId === over.id)

    const newExercises = [...exercises]
    const [removed] = newExercises.splice(oldIndex, 1)
    newExercises.splice(newIndex, 0, removed)

    form.setFieldValue("exercises", newExercises.map((e, i) => ({ ...e, position: i + 1 })) as any)
  }

  const handleExerciseAdd = () => {
    const exercises = form.getValues().exercises
    form.setFieldValue("exercises", [...exercises, {
      clientId: crypto.randomUUID(),
      exerciseId: 0,
      variantId: 0,
      position: exercises.length + 1,
      sets: 0,
    }] as any)
  }

  const handleExerciseDelete = (clientId: string) => {
    const exercises = form.getValues().exercises.filter((e) => e.clientId !== clientId)
    form.setFieldValue("exercises", exercises.map((e, i) => ({ ...e, position: i + 1 })) as any)
  }

  const handleSubmit = () => {
    if (form.validate().hasErrors) return

    const values = form.getValues()
    const { exercises, ...rest } = values
    const cleanValues = {
      ...rest,
      exercises: exercises.map((e) => {
        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        const { clientId, ...rest } = e
        return rest
      }),
    } as T

    setSubmitting(true)
    onSubmit(cleanValues)
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
    onDelete?.(values as SessionUpdate)
      .then(() => form.reset())
      .finally(() => setSubmitting(false))
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
        <Stack gap={0}>
          <p>Exercises</p>
          <p className="text-red-500 text-xs">{form.errors.exercises}</p>
        </Stack>
        <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
          <SortableContext items={form.getValues().exercises.map((e) => e.clientId)} strategy={verticalListSortingStrategy}>
            {form.getValues().exercises.map((exercise) => <Exercise key={exercise.clientId} clientId={exercise.clientId} form={form} onDelete={handleExerciseDelete} />)}
          </SortableContext>
        </DndContext>
        <div onClick={handleExerciseAdd} className="cursor-pointer rounded-sm border-4 border-dotted border-gray-200 py-2 px-4">
          <p className="text-center text-neutral-400">Add exercise</p>
        </div>
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

type ExerciseProps<T extends SessionCreate | SessionUpdate> = {
  form: UseFormReturnType<T>;
  clientId: string;
  onDelete: (clientId: string) => void;
}

const Exercise = <T extends SessionCreate | SessionUpdate>({ form, clientId, onDelete }: ExerciseProps<T>) => {
  const { data: exercises } = useExerciseGetAll()

  const { attributes, listeners, setNodeRef, transform, transition } = useSortable({ id: clientId })

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  }

  const idx = form.getValues().exercises.findIndex(e => e.clientId === clientId)
  const sessionExercise = form.getValues().exercises[idx]

  const options = exercises?.flatMap(e => {
    return [{ label: e.name, value: String(e.id) }, ...e.variants.map(v => ({ label: `${e.name} - ${v.variant}`, value: `${e.id}:${v.id}` }))]
  }) ?? []

  const selectedValue = sessionExercise.variantId ? `${sessionExercise.exerciseId}:${sessionExercise.variantId}` : String(sessionExercise.exerciseId)

  return (
    <Stack ref={setNodeRef} style={style} p="xs" className="border border-gray-200 rounded-sm">
      <Group justify="space-between">
        <ActionIcon variant="subtle" color="gray" {...attributes} {...listeners}>
          <FaGripVertical />
        </ActionIcon>
        <Select
          value={selectedValue}
          data={options}
          onChange={(value) => {
            if (!value) return
            const [id, variant] = value.split(":")
            form.setFieldValue(`exercises.${idx}.exerciseId`, Number(id) as any)
            form.setFieldValue(`exercises.${idx}.variantId`, variant ? Number(variant) : undefined as any)
          }}
          error={form.errors[`exercises.${idx}.exerciseId`]}
        />
        <ActionIcon onClick={() => onDelete(clientId)} variant="subtle" color="red">
          <FaTrashCan />
        </ActionIcon>
      </Group>
      <SimpleGrid cols={{ base: 1, md: 2 }}>
        <NumberInput
          label="Sets"
          description="Amount of sets"
          placeholder="1"
          min={1}
          required
          {...form.getInputProps(`exercises.${idx}.sets`)}
        />
        <NumberInput
          label="Reps"
          description="Amount of reps"
          {...form.getInputProps(`exercises.${idx}.reps`)}
        />
        <NumberInput
          label="Weight"
          description="Added weight"
          suffix=" kg"
          {...form.getInputProps(`exercises.${idx}.weight`)}
        />
        <NumberInput
          label="Duration"
          description="Amount of seconds to hold"
          suffix=" sec"
          {...form.getInputProps(`exercises.${idx}.durationS`)}
        />
      </SimpleGrid>
    </Stack>
  )
}
