import './Start.scss'
import { Base } from '../FeatureComps/Base'
import { StartCarousel } from './StartCarousel'
import { StartCards } from './StartCards'

export const StartFeature: React.FC = () => {
  return (
    <Base className="start-feature">
      <aside>
        <StartCards />
      </aside>

      <main>
        <StartCarousel />
      </main>
    </Base>
  )
}
