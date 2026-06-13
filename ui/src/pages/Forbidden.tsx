import { useAuth } from "@/lib/hooks/useAuth";
import { Text, Center, Stack, Title, Button } from "@mantine/core";
import { useNavigate } from "@tanstack/react-router";

export const Forbidden = () => {
  const { logout } = useAuth()
  const navigate = useNavigate()

  const handleReturn = () => {
    logout()
    navigate({ to: "/" })
  }

  return (
    <Center className="h-screen">
      <Stack align="center" gap={0}>
        <Text fw={600}>403</Text>
        <Title fw={600} className="mt-2">Forbidden</Title>
        <Button onClick={handleReturn} className="mt-6">
          Naar de startpagina
        </Button>
      </Stack>
    </Center>
  );
}
