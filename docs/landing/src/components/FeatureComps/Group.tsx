import './Group.scss'

export const Group: React.FC<{
  children: React.ReactNode
}> = ({ children }) => {
  return <div className="feature-group container">{children}</div>
}
