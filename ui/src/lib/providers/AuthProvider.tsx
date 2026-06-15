import { notifications } from "@mantine/notifications"
import { PropsWithChildren, useCallback, useMemo } from "react"
import { isResponseNot200Error } from "../api/query"
import { useUser, useUserLogin, useUserLogout } from "../api/user"
import { AuthContext } from "../contexts/authContext"

export const AuthProvider = ({ children }: PropsWithChildren) => {
  const { data, isLoading, error } = useUser()
  const { mutate: logoutMutation } = useUserLogout()

  const isUnauthorized = !!error && isResponseNot200Error(error) && error.response.status === 401
  const forbidden = !!error && isResponseNot200Error(error) && error.response.status === 403

  const user = isUnauthorized ? undefined : (data ?? undefined)

  const logout = useCallback(() => {
    logoutMutation(undefined, {
      onSuccess: () =>
        notifications.show({
          color: "green",
          message: "Logged out",
        }),
      onError: (err) => {
        throw new Error(`Logout failed ${err}`)
      },
    })
  }, [logoutMutation])

  const value = useMemo(() => ({
    user: user ?? { id: 0, uid: "", name: "" },
    isLoading,
    forbidden,
    login: useUserLogin,
    logout,
  }), [user, isLoading, forbidden, logout])

  return <AuthContext value={value}>{children}</AuthContext>
}
