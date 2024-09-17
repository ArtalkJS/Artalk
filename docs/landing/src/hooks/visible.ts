import { useEffect, useState } from 'react'

export const useIsElementVisible = (target: React.RefObject<HTMLDivElement>) => {
  const [isVisible, setIsVisible] = useState(false)

  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        setIsVisible(entries[0].isIntersecting)
      },
      {
        threshold: 0.4,
      },
    )

    target.current && observer.observe(target.current)

    return () => {
      target.current && observer.unobserve(target.current)
    }
  }, [target])

  return isVisible
}
