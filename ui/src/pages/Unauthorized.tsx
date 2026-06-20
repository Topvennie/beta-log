import { useAuth } from "@/lib/hooks/useAuth";
import { Text, Center, Stack, Title, Button } from "@mantine/core";
import { useNavigate } from "@tanstack/react-router";

export const Unauthorized = () => {
  const { logout } = useAuth()
  const navigate = useNavigate()

  const handleReturn = () => {
    logout()
    navigate({ to: "/" })
  }

  return (
    <Center className="h-screen">
      <Stack align="center" gap={0}>
        <Text fw={600}>401</Text>
        <Title fw={600} className="mt-2">Unauthorized</Title>
        <Button onClick={handleReturn} className="mt-6">
          Go back to the startpage and log out
        </Button>
      </Stack>
    </Center>
  );
}
