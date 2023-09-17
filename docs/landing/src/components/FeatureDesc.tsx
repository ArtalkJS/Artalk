import './FeatureDesc.scss'

export const FeatureDesc: React.FC<{
  children: React.ReactNode
}> = ({ children }) => {
  return (
    <div className="feature-desc">
      {children}
    </div>
  )
}
