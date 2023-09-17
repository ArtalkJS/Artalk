import './FeatureTitle.scss'
import { Reveal } from './Reveal'

interface FeatureTitleProps {
  text: string
  color: string
}

export const FeatureTitle: React.FC<FeatureTitleProps> = (props) => {
  return (
    <Reveal>
      <div className='feature-title' style={{color: props.color}}>
        <span className='bold'>{props.text.slice(0, 1)}</span>{props.text.slice(1, props.text.length)}
      </div>
    </Reveal>
  )
}
