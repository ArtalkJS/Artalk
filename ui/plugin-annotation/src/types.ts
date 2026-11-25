import type { CommentData } from 'artalk'

export interface Anchor {
  quote: {
    type: string
    exact: string
    prefix?: string
    suffix?: string
  }
  position?: { start: number; end: number }
  section?: string
}

export interface AnnotationConfig {
  enabled: boolean
  root: string | HTMLElement
  hideInList: boolean
  allowGuest: boolean
}

export interface AnnotationPluginOptions extends Partial<AnnotationConfig> {
  excludeSelector?: string
  panelContainer?: HTMLElement | string
}

export type AnnotationCommentData = CommentData & {
  anchor?: Anchor
  is_annotation?: boolean
  anchor_resolved?: boolean
}

export interface AnnotationThread {
  root: AnnotationCommentData
  replies: AnnotationCommentData[]
  orphan: boolean
}

export type Filter = 'all' | 'mine' | 'resolved' | 'pending'
export type SortMode = 'position' | 'newest'
