import { AuthLayout } from "@/layout/AuthLayout"
import { NavLayout } from "@/layout/NavLayout"
import { AuthProvider } from "@/lib/providers/AuthProvider"
import { BreadcrumbProvider } from "@/lib/providers/BreadcrumbProvider"
import { HeaderContentProvider } from "@/lib/providers/HeaderContentProvider"
import { Outlet } from "@tanstack/react-router"

export const Index = () => {
  return (
    <AuthProvider>
      <BreadcrumbProvider>
        <HeaderContentProvider>
          <AuthLayout>
            <NavLayout>
              <Outlet />
            </NavLayout>
          </AuthLayout>
        </HeaderContentProvider>
      </BreadcrumbProvider>
    </AuthProvider>
  )
}

