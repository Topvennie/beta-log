import { LinkButton } from "@/components/atoms/LinkButton";
import { Breadcrumb as BreadcrumbType } from "@/lib/contexts/breadcrumbContext";
import { useAuth } from "@/lib/hooks/useAuth";
import { useBreadcrumbs } from "@/lib/hooks/useBreadcrumb";
import { AppShell, Avatar, Burger, Group, Menu, ScrollArea, Stack } from "@mantine/core";
import { useDisclosure, useHover } from '@mantine/hooks';
import { LinkProps, useNavigate } from "@tanstack/react-router";
import { Fragment, PropsWithChildren, ReactNode } from "react";
import { FaArrowRightFromBracket, FaChevronRight, FaDumbbell, FaGear } from "react-icons/fa6";
import { LuLayoutDashboard, LuLayoutList } from "react-icons/lu";

type Props = PropsWithChildren

const breakpoint = "md"

type Route = {
  title: string;
  icon: ReactNode;
  link: LinkProps;
}

const routes: Route[] = [
  {
    title: "Dashboard",
    icon: <LuLayoutDashboard />,
    link: { to: "/" },
  },
  {
    title: "Sessions",
    icon: <LuLayoutList />,
    link: { to: "/sessions" },
  },
  {
    title: "Exercises",
    icon: <FaDumbbell />,
    link: { to: "/exercises" },
  },
]

const Breadcrumb = ({ breadcrumb: { title, link } }: { breadcrumb: BreadcrumbType }) => {
  return (
    <LinkButton
      variant="transparent"
      color="black"
      p={0}
      className="hover:underline underline-offset-2"
      {...link}
    >
      <p className="font-semibold">{title}</p>
    </LinkButton>
  )
}

const NavLink = ({ route: { title, icon, link } }: { route: Route }) => {
  return (
    <LinkButton
      fullWidth
      leftSection={icon}
      variant="subtle"
      justify="start"
      color="black"
      size="md"
      className="border-l-4"
      {...link}
      activeProps={{ bg: "blue.1", className: "border-l-blue-500" }}
    >
      <p>{title}</p>
    </LinkButton>
  )
}

export const NavLayout = ({ children }: Props) => {
  const { user, logout } = useAuth()
  const { state: breadcrumbs } = useBreadcrumbs();

  const [opened, { toggle }] = useDisclosure();
  const { hovered, ref } = useHover()

  const navigate = useNavigate()

  return (
    <AppShell
      layout="alt"
      header={{ height: 60 }}
      navbar={{ width: 200, breakpoint, collapsed: { mobile: !opened } }}
    >
      <AppShell.Header>
        <Group h="100%" px="xl">
          <Burger opened={opened} onClick={toggle} hiddenFrom={breakpoint} size="sm" />
          <Group>
            {breadcrumbs.map((breadcrumb, idx) => (
              <Fragment key={breadcrumb.title}>
                <Breadcrumb breadcrumb={breadcrumb} />
                {idx < breadcrumbs.length - 1 && <p>/</p>}
              </Fragment>
            ))}
          </Group>
        </Group>
      </AppShell.Header>

      <AppShell.Navbar p={0}>
        <AppShell.Section p="md">
          <p className="uppercase font-bold text-2xl">Beta Log</p>
        </AppShell.Section>
        <AppShell.Section grow my="sm" component={ScrollArea} px="md">
          <Stack gap={4}>
            {routes.map(route => <NavLink key={route.title} route={route} />)}
          </Stack>
        </AppShell.Section>
        <AppShell.Section p="md" className="border-t border-gray-200">
          <Menu position="right-end" trigger="click-hover">
            <Menu.Target>
              <Group ref={ref} className="cursor-pointer">
                <Avatar radius="sm">{user.name?.[0] ?? ""}</Avatar>
                {user.name}
                <FaChevronRight className={`ml-auto text-neutral-400 duration-300 ${hovered ? "translate-x-1" : ""}`} />
              </Group>
            </Menu.Target>
            <Menu.Dropdown>
              <Menu.Item onClick={() => navigate({ to: "/settings" })} leftSection={<FaGear />}>
                Settings
              </Menu.Item>
              <Menu.Divider />
              <Menu.Item color="red" onClick={logout} leftSection={<FaArrowRightFromBracket />}>
                Logout
              </Menu.Item>
            </Menu.Dropdown>
          </Menu>
        </AppShell.Section>
      </AppShell.Navbar>

      <AppShell.Main>
        <ScrollArea p="xl" h="calc(100vh - var(--app-shell-header-height, 0px) - var(--app-shell-footer-height, 0px))">
          {children}
        </ScrollArea>
      </AppShell.Main>
    </AppShell>
  )
}
