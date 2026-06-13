import { useAuth } from "@/lib/hooks/useAuth";
import { Button, Center, Paper, Stack, Text, Title } from "@mantine/core";
import { FaOpenid } from "react-icons/fa6"

export const Login = () => {
  const { login } = useAuth();

  return (
    <Center h="100vh">
      <Paper shadow="sm" p="xl" className="w-96">
        <Stack align="center">
          <Title>Login</Title>
          <Text c="gray">Beta-Log</Text>
          <Button onClick={login} size="lg" className="my-12">
            <FaOpenid size={"1.7rem"} className="mr-2" />
            OIDC
          </Button>
        </Stack>
      </Paper>
    </Center>
  );
};
