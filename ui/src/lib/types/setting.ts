import z from "zod";
import { API } from "./api";
import { JSONBody } from "./general";

export enum GradeSystem {
  Font = "font",
  V = "v",
}

export interface Setting {
  gradeSystem: GradeSystem;
  climbToploggerUserId?: string;
  climbToploggerAuthToken?: string;
  climbToploggerRefreshToken?: string;
}

// Converts

export const convertSetting = (s: API.Setting): Setting => ({
  gradeSystem: s.grade_system as GradeSystem,
  climbToploggerUserId: s.climb_toplogger_user_id,
  climbToploggerAuthToken: s.climb_toplogger_auth_token,
  climbToploggerRefreshToken: s.climb_toplogger_refresh_token,
})

export const convertSettingUpdateGradeSystem = (s: Setting): SettingUpdateGradeSystem => ({
  gradeSystem: s.gradeSystem,
})

export const convertSettingUpdateToploggerSchema = (s: Setting): SettingUpdateToplogger => ({
  climbToploggerUserId: s.climbToploggerUserId,
  climbToploggerAuthToken: s.climbToploggerAuthToken,
  climbToploggerRefreshToken: s.climbToploggerRefreshToken
})

// Schemas

export const settingUpdateGradeSystem = z.object({
  gradeSystem: z.enum(GradeSystem),
})
export type SettingUpdateGradeSystem = z.infer<typeof settingUpdateGradeSystem> & JSONBody

export const settingUpdateToploggerSchema = z.object({
  climbToploggerUserId: z.string().optional(),
  climbToploggerAuthToken: z.string().optional(),
  climbToploggerRefreshToken: z.string().optional(),
}).superRefine((args, ctx) => {
  if (!!args.climbToploggerUserId !== !!args.climbToploggerAuthToken || !!args.climbToploggerUserId !== !!args.climbToploggerRefreshToken) {
    ctx.addIssue({
      code: "custom",
      path: ["climbToploggerUserId"],
      message: "User id, auth token and refresh need all to be empty or filled in."
    })
    ctx.addIssue({
      code: "custom",
      path: ["climbToploggerAuthToken"],
      message: "User id, auth token and refresh need all to be empty or filled in."
    })
    ctx.addIssue({
      code: "custom",
      path: ["climbToploggerRefreshToken"],
      message: "User id, auth token and refresh need all to be empty or filled in."
    })
  }
})
export type SettingUpdateToplogger = z.infer<typeof settingUpdateToploggerSchema> & JSONBody;
