import type { Context } from 'artalk'
import type { AnnotationConfig, AnnotationPluginOptions } from './types'

export const DEFAULT_ANNOTATION_CONFIG: AnnotationConfig = {
  enabled: false,
  root: 'article',
  hideInList: true,
  allowGuest: true,
}

export function resolveConfig(
  ctx: Context,
  options: AnnotationPluginOptions = {},
): AnnotationConfig {
  const conf = ctx.getConf()
  return {
    ...DEFAULT_ANNOTATION_CONFIG,
    ...conf.annotation,
    ...options,
  }
}

export function resolveElement(
  value: string | HTMLElement | undefined,
  fallback: HTMLElement | null = null,
): HTMLElement | null {
  if (!value) return fallback
  if (value instanceof HTMLElement) return value
  return document.querySelector<HTMLElement>(value) || fallback
}
