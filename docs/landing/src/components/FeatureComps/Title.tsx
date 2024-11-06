import './Title.scss'
import { Reveal } from '../Reveal'

interface TitleProps {
  text: string
  color: string
}

export const Title: React.FC<TitleProps> = (props) => {
  return (
    <Reveal>
      <div className="feature-title" style={{ color: props.color }}>
        {props.text.slice(0, 1)}
        {props.text.slice(1, props.text.length)}
      </div>
    </Reveal>
  )
}
