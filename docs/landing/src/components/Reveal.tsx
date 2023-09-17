import React, { useEffect, useRef, useState } from 'react'
import './Reveal.scss'

interface RevelProps {
  children: React.ReactNode
  duration?: number
  delay?: number
  threshold?: number
}

export const Reveal: React.FC<RevelProps> = (props) => {
  props = {...{
    // default config
    threshold: 0.5,
    duration: 1000,
    delay: 0
  }, ...props}

  const [isVisible, setIsVisible] = useState(false)
  const elementRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const observer = new IntersectionObserver((entries) => {
      if (entries[0].isIntersecting) {
        setIsVisible(true)
        elementRef.current && observer.unobserve(elementRef.current)
      }
    }, {
      threshold: props.threshold
    })

    elementRef.current && observer.observe(elementRef.current)

    return () => {
      elementRef.current && observer.unobserve(elementRef.current)
    }
  }, [props.threshold])

  useEffect(() => {
    if (isVisible) {
      elementRef.current?.classList.add('animate')

      const animationEndHandler = () => {
        elementRef.current?.classList.remove('animate')
        elementRef.current?.classList.add('show')
      }

      elementRef.current?.addEventListener('animationend', animationEndHandler)

      return () => {
        elementRef.current?.removeEventListener('animationend', animationEndHandler)
      }
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
