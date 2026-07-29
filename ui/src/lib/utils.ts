import { type ClassValue, clsx } from "clsx";
import { format } from "date-fns";
import { twMerge } from "tailwind-merge";

export function camelToSnake(obj: unknown): unknown {
  if (obj === null || obj === undefined) {
    return obj;
  }

  if (Array.isArray(obj)) {
    return obj.map(camelToSnake);
  }

  if (obj instanceof Date) {
    return obj.toISOString()
  }

  if (typeof obj === "object") {
    return Object.fromEntries(
      Object.entries(obj).filter(([key]) => key?.[0] !== "_").map(([key, value]) => [
        stringCamelToSnake(key),
        camelToSnake(value),
      ]),
    );
  }

  return obj;
}

function stringCamelToSnake(str: string) {
  return str.replace(/[A-Z]+/g, l => `_${l.toLowerCase()}`);
}

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function formatDate(date: Date) {
  return format(date, "eee dd MMMM, HH:mm");
}

export function capitalize(text: string) {
  if (text.length <= 1) return text.toUpperCase()

  return text[0].toUpperCase() + text.slice(1)
}
