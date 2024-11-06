import React, { useEffect, useRef } from 'react'
import LazyLoad, { type ILazyLoadInstance } from 'vanilla-lazyload'
import './LazyImage.scss'

let lazyLoadInstance: ILazyLoadInstance

const LazyImage: React.FC<React.ImgHTMLAttributes<HTMLImageElement>> = ({
  src,
  alt,
  className = '',
  ...props
}) => {
  const imageRef = useRef(null)

  useEffect(() => {
    if (!lazyLoadInstance) {
      lazyLoadInstance = new LazyLoad({
        elements_selector: '.lazyload',
        threshold: 0,
      })
    }
    lazyLoadInstance.update()
  }, [src])

  return (
    <img ref={imageRef} className={`lazyload ${className}`} data-src={src} alt={alt} {...props} />
  )
}

export default LazyImage
