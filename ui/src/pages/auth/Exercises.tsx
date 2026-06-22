import { ExerciseForm } from "@/components/exercise/ExerciseForm"
import { useExerciseCreate, useExerciseDelete, useExerciseGetAll, useExerciseUpdate } from "@/lib/api/exercise"
import { useBreadcrumb } from "@/lib/hooks/useBreadcrumb"
import { ExerciseCreate, Exercise as ExerciseType, ExerciseUpdate } from "@/lib/types/exercise"
import { Badge, Group, Stack, Title } from "@mantine/core"
import { useHover } from "@mantine/hooks"
import { notifications } from "@mantine/notifications"
import { useState } from "react"
import { FaChevronRight } from "react-icons/fa6"

export const Exercises = () => {
  useBreadcrumb({ title: "Exercises", weight: 10, link: { to: "/exercises" } })

  const { data: exercises, isLoading } = useExerciseGetAll()
  const exerciseCreate = useExerciseCreate()
  const exerciseUpdate = useExerciseUpdate()
  const exerciseDelete = useExerciseDelete()

  const [selected, setSelected] = useState<ExerciseType | undefined>(undefined)

  if (isLoading || exercises === undefined) return null

  const handleCreate = (exercise: ExerciseCreate) => {
    return exerciseCreate.mutateAsync(exercise, {
      onSuccess: () => {
        notifications.show({ color: "green", title: "Exercise", message: `Created ${exercise.name}` })
        setSelected(undefined)
      },
    })
  }

  const handleUpdate = (exercise: ExerciseUpdate) => {
    return exerciseUpdate.mutateAsync(exercise, {
      onSuccess: () => {
        notifications.show({ color: "green", title: "Exercise", message: `Updated ${exercise.name}` })
        setSelected(undefined)
      },
    })
  }

  const handleDelete = ({ id }: Pick<ExerciseType, "id">) => {
    return exerciseDelete.mutateAsync({ id }, {
      onSuccess: () => {
        notifications.show({ color: "green", title: "Exercise", message: `Deleted` })
        setSelected(undefined)
      },
    })
  }

  return (
    <div className="flex md:flex-row gap-8 h-full">
      <Stack gap="xl" className="flex-3 h-full">
        <Group justify="space-between">
          <Title order={1}>Exercises</Title>
          <p className="text-neutral-400">{exercises.length}</p>
        </Group>

        <Stack gap="md">
          {exercises.map(e => <Exercise key={e.id} exercise={e} onClick={setSelected} />)}
        </Stack>
      </Stack>

      <div className="border-l border-gray-200" />

      <div className="flex-2">
        {selected
          ? <ExerciseForm key={selected.id} exercise={selected} onSubmit={handleUpdate} onCancel={() => setSelected(undefined)} onDelete={handleDelete} />
          : <ExerciseForm onSubmit={handleCreate} onCancel={() => null} />
        }
      </div>
    </div>
  )
}

type ExerciseProps = {
  exercise: ExerciseType;
  onClick: (exercise: ExerciseType) => void;
}

const Exercise = ({ exercise, onClick }: ExerciseProps) => {
  const { hovered, ref } = useHover()

  return (
    <Group ref={ref} onClick={() => onClick(exercise)} p="xs" justify="space-between" className="rounded-md cursor-pointer hover:bg-blue-50">
      <Stack>
        <p>{exercise.name}</p>
        <Group>
          {exercise.variants.map(v => <Badge key={v} variant="light" size="sm">{v}</Badge>)}
        </Group>
      </Stack>
      <FaChevronRight className={`text-blue-500 duration-300 ${hovered ? "translate-x-1" : ""}`} />
    </Group>
  )
}
