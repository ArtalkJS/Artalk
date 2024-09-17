import React, { useEffect, useRef, useState } from 'react'
import './Reveal.scss'

interface RevelProps {
  children: React.ReactNode
  duration?: number
  delay?: number
  threshold?: number
}

export const Reveal: React.FC<RevelProps> = (props) => {
  props = {
    ...{
      // default config
      threshold: 0.5,
      duration: 1000,
      delay: 0,
    },
    ...props,
  }

  const [isVisible, setIsVisible] = useState(false)
  const elementRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const observerRefValue = elementRef.current

    if (!('IntersectionObserver' in window)) {
      observerRefValue?.classList.add('show')
      return
    }

    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting) {
          setIsVisible(true)
          observerRefValue && observer.unobserve(observerRefValue)
        }
      },
      {
        threshold: props.threshold,
      },
    )

    observerRefValue && observer.observe(observerRefValue)

    return () => {
      observerRefValue && observer.unobserve(observerRefValue)
    }
  }, [props.threshold])

  useEffect(() => {
    if (!isVisible) return

    const observerRefValue = elementRef.current

    observerRefValue?.classList.add('animate')

    const animationEndHandler = () => {
      observerRefValue?.classList.remove('animate')
      observerRefValue?.classList.add('show')
    }

    observerRefValue?.addEventListener('animationend', animationEndHandler)

    return () => {
      observerRefValue?.removeEventListener('animationend', animationEndHandler)
    }
  }, [isVisible])

  return (
    <div
      ref={elementRef}
      className={['reveal'].join(' ')}
      style={{
        animationDuration: `${props.duration}ms`,
        animationDelay: `${props.delay}ms`,
      }}
    >
      {props.children}
    </div>
  )
}
