import { isResponseNot200Error } from "@/lib/api/query";
import { Button, Text, Center, Stack, Title } from "@mantine/core";
import { ErrorComponentProps, useNavigate } from "@tanstack/react-router";
import { Error404 } from "./404";
import { Forbidden } from "./Forbidden";

export const Error = ({ error, reset }: ErrorComponentProps) => {
  const navigate = useNavigate()

  if (isResponseNot200Error(error)) {
    switch (error.response.status) {
      case 404:
        return <Error404 />
      case 403:
        return <Forbidden />
    }
  }

  const handleReturn = () => {
    reset()
    navigate({ to: "/" })
  }

  return (
    <Center className="h-screen">
      <Stack align="center" gap={0}>
        <Text fw={600}>500</Text>
        <Title fw={600} className="mt-2">Server Error</Title>
        <Button onClick={handleReturn} className="mt-6">
          Go back to the startpage
        </Button>
      </Stack>
    </Center>
  )
}
