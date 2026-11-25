declare module 'dom-anchor-text-quote' {
  export function fromRange(
    root: Node,
    range: Range,
  ): { exact: string; prefix?: string; suffix?: string }

  export function toRange(
    root: Node,
    selector: { exact: string; prefix?: string; suffix?: string },
  ): Range | null
}
