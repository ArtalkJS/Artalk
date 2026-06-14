import { useEffect, useState } from 'react'

export const useIsElementVisible = (target: React.RefObject<HTMLDivElement | null>) => {
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

    const targetEl = target.current
    targetEl && observer.observe(targetEl)

    return () => {
      targetEl && observer.unobserve(targetEl)
    }
  }, [target])

  return isVisible
}
