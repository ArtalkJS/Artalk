import RenderCtx from '../render-ctx'
import Avatar from './avatar'
import Header from './header'
import Content from './content'
import ReplyAt from './reply-at'
import ReplyTo from './reply-to'
import Pending from './pending'
import Actions from './actions'

const Renders = {
  Avatar, Header, Content, ReplyAt,
  ReplyTo, Pending, Actions
}

export default function loadRenders(ctx: RenderCtx) {
  Object.entries(Renders).forEach(([ name, render ]) => {
    render(ctx)
  })
}
