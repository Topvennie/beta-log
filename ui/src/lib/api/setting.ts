import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { convertSetting, SettingUpdateGradeSystem, SettingUpdateToplogger } from "../types/setting";
import { apiGet, apiPut } from "./query";

const ENDPOINT = "setting";

export const useSettingGet = () => {
  return useQuery({
    queryKey: ["setting"],
    queryFn: async () => (await apiGet(ENDPOINT, convertSetting)).data,
  });
}

export const useSettingUpdateGradeSystem = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (setting: SettingUpdateGradeSystem) => apiPut(`${ENDPOINT}/grade_system`, setting, convertSetting),
    onSuccess: () => {
      queryClient.invalidateQueries()

    },
  })
}

export const useSettingUpdateToplogger = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (setting: SettingUpdateToplogger) => apiPut(`${ENDPOINT}/toplogger`, setting, convertSetting),
    onSuccess: () => {
      queryClient.invalidateQueries()

    },
  })
}
