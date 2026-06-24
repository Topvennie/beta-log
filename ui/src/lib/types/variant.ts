import z from "zod";
import { API } from "./api";
import { JSONBody } from "./general";

export interface Variant {
  id: number;
  variant: string;
}

// Converts

export const convertVariant = (v: API.Variant): Variant => ({
  id: v.id,
  variant: v.variant,
})

export const convertVariants = (v: API.Variant[]): Variant[] => v.map(convertVariant)

export const convertVariantUpdateSchema = (v: Variant): VariantUpdateSchema => ({
  id: v.id,
  variant: v.variant,
})

// Schemas

export const variantCreateSchema = z.object({
  variant: z.string().min(1),
})
export type VariantCreateSchema = z.infer<typeof variantCreateSchema> & JSONBody

export const variantUpdateSchema = z.object({
  id: z.number().positive().optional(),
  variant: z.string().min(1),
})
export type VariantUpdateSchema = z.infer<typeof variantUpdateSchema> & JSONBody
