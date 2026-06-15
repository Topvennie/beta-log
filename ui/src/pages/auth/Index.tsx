import { AuthLayout } from "@/layout/AuthLayout"
import { NavLayout } from "@/layout/NavLayout"
import { useAuth } from "@/lib/hooks/useAuth"
import { AuthProvider } from "@/lib/providers/AuthProvider"
import { Button } from "@mantine/core"

export const Index = () => {
  return (
    <AuthProvider>
      <AuthLayout>
        <NavLayout>
          <Inner />
        </NavLayout>
      </AuthLayout>
    </AuthProvider>
  )
}

const Inner = () => {
  const { user, logout } = useAuth()

  return (
    <div>
      <h1>Admin</h1>
      <p>{`Hi ${user?.name}`}</p>
      <Button onClick={logout}>Logout</Button>
    </div>
  )
}
