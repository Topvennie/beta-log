import { Button, Group, Modal, ModalBaseProps, Stack } from "@mantine/core";

interface Props extends Omit<ModalBaseProps, 'title' | 'content'> {
  onConfirm: () => void;
}

export const Confirm = ({ onConfirm, ...props }: Props) => {
  return (
    <Modal title="Confirm Deletion" {...props}>
      <Stack>
        <p>Are you sure you want to delete</p>
        <Group justify="end">
          <Button onClick={props.onClose} variant="outline">
            Cancel
          </Button>
          <Button onClick={onConfirm} color="red">
            Confirm
          </Button>
        </Group>
      </Stack>
    </Modal>
  )
}
