import { GradeSystem } from "./types/setting"

export function gradeValue(grade: string, system: GradeSystem): number {
  if (grade === "" || grade === "?" || grade === "n/a") {
    return 0
  }

  switch (system) {
    case GradeSystem.Font:
      return fontValue(grade)
    case GradeSystem.V:
      return vValue(grade)
  }
}

function fontValue(grade: string): number {
  const m = grade.match(/^(\d)([abc])(\+?)$/)
  if (!m) return 0
  const base = parseInt(m[1], 10)
  const letter = m[2].charCodeAt(0) - 97 // a=0, b=1, c=2
  const plus = m[3] === "+" ? 1 : 0
  return base * 6 + letter * 2 + plus
}

function vValue(grade: string): number {
  const m = grade.match(/^V(\d+)$/)
  if (!m) return 0
  return parseInt(m[1], 10)
}

export function formatGradeValue(value: number, system: GradeSystem): string {
  switch (system) {
    case GradeSystem.Font:
      return formatFontValue(value)
    case GradeSystem.V:
      return formatVValue(value)
  }
}

function formatFontValue(value: number): string {
  if (value < 6) return ""
  const base = Math.floor(value / 6)
  const rem = value % 6
  const letterIdx = Math.floor(rem / 2)
  const plus = rem % 2 === 1
  if (base < 1 || base > 9 || letterIdx < 0 || letterIdx > 2) return ""
  return `${base}${"abc"[letterIdx]}${plus ? "+" : ""}`
}

function formatVValue(value: number): string {
  if (value < 0 || value > 17) return ""
  return `V${value}`
}

export function gradeHue(grade: string, system: GradeSystem) {
  if (grade == "" || grade == "?" || grade == "n/a") {
    return "neutral"
  }

  switch (system) {
    case GradeSystem.Font:
      return gradeHueFont(grade)
    case GradeSystem.V:
      return gradeHueV(grade)
  }
}

function gradeHueFont(grade: string) {
  switch (grade) {
    case "2a": case "2a+": case "2b": case "2b+": case "2c": case "2c+":
    case "3a": case "3a+": case "3b": case "3b+": case "3c": case "3c+":
      return "green"
    case "4a": case "4a+": case "4b": case "4b+": case "4c": case "4c+":
    case "5a": case "5a+":
      return "yellow"
    case "5b": case "5b+": case "5c": case "5c+":
    case "6a": case "6a+":
      return "orange"
    case "6b": case "6b+":
      return "blue"
    case "6c": case "6c+":
      return "red"
    case "7a":
      return "black"
    case "7a+":
      return "white"
    case "7b": case "7b+": case "7c": case "7c+":
    case "8a": case "8a+": case "8b": case "8b+": case "8c": case "8c+":
    case "9a": case "9a+":
      return "purple"
    default:
      return "neutral"
  }
}

function gradeHueV(grade: string) {
  switch (grade) {
    case "V0":
      return "green"
    case "V1":
      return "yellow"
    case "V2": case "V3":
      return "orange"
    case "V4":
      return "blue"
    case "V5":
      return "red"
    case "V6":
      return "black"
    case "V7":
      return "white"
    case "V8": case "V9": case "V10": case "V11": case "V12": case "V13": case "V14": case "V15": case "V16": case "V17":
      return "purple"
    default:
      return "neutral"
  }
}
