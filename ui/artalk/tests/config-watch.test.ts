import { expect, it, vi } from 'vitest'
import { Renderer } from 'marked'
import Context from '@/context'
import { Defaults } from '@/defaults'
import { EventManager } from '@/lib/event-manager'
import { watchConf } from '@/lib/watch-conf'
import { ConfigService } from '@/services/config'
import { EventsService } from '@/services/events'
import type { Config, EventPayloadMap } from '@/types'

function observe<K extends keyof Config>(key: K, initial: Config[K]) {
  const events = new EventManager<EventPayloadMap>()
  let conf: Config = { ...Defaults, [key]: initial }
  const changed = vi.fn()
  watchConf({ keys: [key], effect: changed, getConf: () => conf, getEvents: () => events })
  events.trigger('mounted')
  changed.mockClear()
  return {
    changed,
    update(value: Config[K]) {
      conf = { ...conf, [key]: value }
      events.trigger('updated', conf)
    },
  }
}

it('observes function configuration by reference without resetting unchanged functions', () => {
  const ctx = new Context(document.createElement('div'))
  EventsService(ctx)
  ConfigService(ctx)
  ctx.updateConf({ fetchCommentsOnInit: false })
  const changed = vi.fn()
  ctx.watchConf(['imgUploader'], changed)
  ctx.trigger('mounted')
  const first = async () => 'first.png'
  const second = async () => 'second.png'
  ctx.updateConf({ imgUploader: first })
  ctx.updateConf({ imgUploader: first })
  ctx.updateConf({ imgUploader: second })
  ctx.updateConf({ imgUploader: undefined })
  expect(changed.mock.calls.map(([conf]) => conf.imgUploader)).toEqual([
    undefined,
    first,
    second,
    undefined,
  ])
  ctx.destroy()
})

it('compares nested configuration by value regardless of property order', () => {
  const { changed, update } = observe('emoticons', {
    name: 'Faces',
    items: [{ key: ':smile:', value: '/smile.png' }],
  })
  update({ items: [{ value: '/smile.png', key: ':smile:' }], name: 'Faces' })
  expect(changed).not.toHaveBeenCalled()
  update({ name: 'Faces', items: [{ key: ':smile:', value: '/new-smile.png' }] })
  expect(changed).toHaveBeenCalledTimes(1)
})

it('observes nested callback changes and removal without calling the callbacks', () => {
  const first = vi.fn()
  const second = vi.fn()
  const { changed, update } = observe('markedOptions', { walkTokens: first, breaks: true })
  update({ breaks: true, walkTokens: first })
  expect(changed).not.toHaveBeenCalled()
  update({ breaks: true, walkTokens: second })
  expect(changed).toHaveBeenCalledTimes(1)
  update({ breaks: true })
  expect(changed).toHaveBeenCalledTimes(2)
  expect(first).not.toHaveBeenCalled()
  expect(second).not.toHaveBeenCalled()
})

it('compares replacer arrays by callback identity and order', () => {
  const replace = (value: string) => (text: string) => text.replace('sample', value)
  const first = replace('first')
  const second = replace('second')
  const { changed, update } = observe('markedReplacers', [first, second])
  update([first, second])
  expect(changed).not.toHaveBeenCalled()
  update([second, first])
  expect(changed).toHaveBeenCalledTimes(1)
  update([second, replace('first')])
  expect(changed).toHaveBeenCalledTimes(2)
})

it.each([
  { name: 'empty and sparse arrays', before: [], after: Array(1) },
  { name: 'trailing array holes', before: ['smile'], after: ['smile'].concat(Array(1)) },
  {
    name: 'different occupied indexes',
    before: ['smile'].concat(Array(1)),
    after: Array(1).concat('smile'),
  },
  { name: 'arrays and numeric-key objects', before: ['smile'], after: { 0: 'smile' } },
])('detects changes between $name', ({ before, after }) => {
  const { changed, update } = observe('emoticons', before)
  update(after)
  expect(changed).toHaveBeenCalledTimes(1)
})

it('treats a constructor field as data in a plain configuration object', () => {
  const { changed, update } = observe('emoticons', { constructor: 'custom', name: 'Faces' })
  update({ name: 'Faces', constructor: 'custom' })
  expect(changed).not.toHaveBeenCalled()
  update({ name: 'Faces', constructor: 'other' })
  expect(changed).toHaveBeenCalledTimes(1)
})

it('compares Markdown renderer instances by reference', () => {
  const renderer = new Renderer()
  const { changed, update } = observe('markedOptions', { renderer })
  update({ renderer })
  expect(changed).not.toHaveBeenCalled()
  update({ renderer: new Renderer() })
  expect(changed).toHaveBeenCalledTimes(1)
})

it('compares DOM configuration values by reference', () => {
  const element = document.createElement('div')
  const { changed, update } = observe('el', element)
  update(element)
  expect(changed).not.toHaveBeenCalled()
  update(document.createElement('div'))
  expect(changed).toHaveBeenCalledTimes(1)
})
