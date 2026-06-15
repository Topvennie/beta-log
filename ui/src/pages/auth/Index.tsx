import { AuthLayout } from "@/layout/AuthLayout"
import { NavLayout } from "@/layout/NavLayout"
import { AuthProvider } from "@/lib/providers/AuthProvider"
import { BreadcrumbProvider } from "@/lib/providers/BreadcrumbProvider"
import { Outlet } from "@tanstack/react-router"

export const Index = () => {
  return (
    <AuthProvider>
      <BreadcrumbProvider>
        <AuthLayout>
          <NavLayout>
            <Outlet />
          </NavLayout>
        </AuthLayout>
      </BreadcrumbProvider>
    </AuthProvider>
  )
}

