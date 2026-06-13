import { AuthLayout } from "@/layout/AuthLayout"
import { useAuth } from "@/lib/hooks/useAuth"
import { AuthProvider } from "@/lib/providers/AuthProvider"
import { Button } from "@mantine/core"

export const Index = () => {
  return (
    <AuthProvider>
      <Inner />
    </AuthProvider>
  )
}

const Inner = () => {
  const { user, logout } = useAuth()

  return (
    <AuthLayout>
      <div>
        <h1>Admin</h1>
        <p>{`Hi ${user?.name}`}</p>
        <Button onClick={logout}>Logout</Button>
      </div>
    </AuthLayout>
  )
}
