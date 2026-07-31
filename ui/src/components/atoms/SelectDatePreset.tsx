import { Select } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { FaCalendarDay } from "react-icons/fa6";

type Props = {
  value: DatePreset;
  setValue: (value: DatePreset) => void;
}

const datePresets = ["1 Month", "3 Months", "1 Year", "All Time"] as const
export type DatePreset = typeof datePresets[number];

export const SelectDatePreset = ({ value, setValue }: Props) => {
  const [dropdownOpened, { open, close }] = useDisclosure()

  const onChange = (v: DatePreset | null) => {
    if (!v) return

    setValue(v)
    close()
  }

  return (
    <Select
      data={datePresets.map(p => ({ value: p, label: p }))}
      value={value}
      onChange={onChange}
      dropdownOpened={dropdownOpened}
      onDropdownOpen={open}
      onDropdownClose={close}
      clearable={false}
      allowDeselect={false}
      leftSection={<FaCalendarDay />}
      size="xs"
    />
  )
}