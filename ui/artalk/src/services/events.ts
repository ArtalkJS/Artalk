import type { ArtalkPlugin, EventPayloadMap } from '@/types'
import { EventManager } from '@/lib/event-manager'

export const EventsService: ArtalkPlugin = (ctx) => {
  ctx.provide('events', () => new EventManager<EventPayloadMap>())
}
