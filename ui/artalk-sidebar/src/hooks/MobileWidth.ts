export const useMobileWidth = () => {
  const isMobileWidth = ref(false)
  const maxMobileWidth = 1023

  const debounce = (fn: () => void, delay = 100) => {
    let timer: any
    return () => {
      clearTimeout(timer)
      timer = setTimeout(fn, delay)
    }
  }

  const checkMobileWidth = debounce(() => {
    isMobileWidth.value = window.innerWidth <= maxMobileWidth
  })

  onMounted(() => {
    checkMobileWidth()
    window.addEventListener('resize', checkMobileWidth)
  })

  onBeforeUnmount(() => {
    window.removeEventListener('resize', checkMobileWidth)
  })

  return isMobileWidth
}
