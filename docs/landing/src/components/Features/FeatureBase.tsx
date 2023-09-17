import React, { useEffect, useRef } from 'react'
import './FeatureBase.scss'
import { useIsElementVisible } from '../../hooks/visible'

interface FeatureBaseProps extends React.ComponentProps<'div'> {
  onVisibleChange?: (visible: boolean) => void
}

export const FeatureBase: React.FC<FeatureBaseProps> = (props) => {
  const { className, onVisibleChange, ...otherProps } = props


  const refElem = useRef<HTMLDivElement>(null)
  const isSlightFeatureVisible = useIsElementVisible(refElem)

  const classNames = ['feature', 'container', isSlightFeatureVisible ? 'visible' : '', ...(className?.split(' ') || [])]

  useEffect(() => {
    onVisibleChange && onVisibleChange(isSlightFeatureVisible)
  }, [onVisibleChange, isSlightFeatureVisible])

  return (
    <div ref={refElem} className={classNames.join(' ')} {...otherProps}>
      {props.children}
    </div>
  )
}
