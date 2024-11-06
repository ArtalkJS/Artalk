import './Desc.scss'

export const Desc: React.FC<{
  children: React.ReactNode
}> = ({ children }) => {
  return <div className="feature-desc">{children}</div>
}
