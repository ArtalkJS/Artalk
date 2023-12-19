import type Render from '../render'
import Avatar from './avatar'
import Header from './header'
import Content from './content'
import ReplyAt from './reply-at'
import ReplyTo from './reply-to'
import Pending from './pending'
import Actions from './actions'

const Renders = {
  Avatar,
  Header,
  Content,
  ReplyAt,
  ReplyTo,
  Pending,
  Actions,
}

export default function loadRenders(r: Render) {
  Object.entries(Renders).forEach(([name, render]) => {
    render(r)
  })
}
