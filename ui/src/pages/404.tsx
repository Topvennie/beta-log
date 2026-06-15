import { LinkButton } from "@/components/atoms/LinkButton"
import { Center, Stack, Text, Title } from "@mantine/core"

export const Error404 = () => {
  return (
    <div className="h-screen">
      <Center h="100%">
        <Stack align="center" gap={0}>
          <Text fw={600}>404</Text>
          <Title fw={600} className="mt-2">Page not found</Title>
          <LinkButton className="mt-6" to="/">
            To the startpage
          </LinkButton>
        </Stack>
      </Center>
    </div>
  )
}
