import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { convertSetting, SettingToploggerUpdate } from "../types/setting";
import { apiGet, apiPut } from "./query";

const ENDPOINT = "setting";

export const useSettingGet = () => {
  return useQuery({
    queryKey: ["setting"],
    queryFn: async () => (await apiGet(ENDPOINT, convertSetting)).data,
  });
}

export const useSettingToploggerUpdate = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (setting: SettingToploggerUpdate) => apiPut(`${ENDPOINT}/toplogger`, setting, convertSetting),
    onSuccess: () => {
      queryClient.invalidateQueries()

    },
  })
}
