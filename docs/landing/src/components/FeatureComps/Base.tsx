import React, { useEffect, useRef } from 'react'
import './Base.scss'
import { useIsElementVisible } from '../../hooks/visible'

interface FeatureBaseProps extends React.ComponentProps<'div'> {
  onVisibleChange?: (visible: boolean) => void
}

export const Base: React.FC<FeatureBaseProps> = (props) => {
  const { className, onVisibleChange, ...otherProps } = props

  const refElem = useRef<HTMLDivElement>(null)
  const isVisible = useIsElementVisible(refElem)

  const classNames = [
    'feature',
    'container',
    isVisible ? 'visible' : '',
    ...(className?.split(' ') || []),
  ]

  useEffect(() => {
    onVisibleChange && onVisibleChange(isVisible)
  }, [onVisibleChange, isVisible])

  return (
    <div ref={refElem} className={classNames.join(' ')} {...otherProps}>
      {props.children}
    </div>
  )
}
