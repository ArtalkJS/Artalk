import type { ArtalkPlugin } from '@/types'
import * as Utils from '@/lib/utils'
import $t from '@/i18n'

export const Dropdown: ArtalkPlugin = (ctx) => {
  const reloadUseParamsEditor = (func: (p: any) => void) => {
    ctx.conf.listFetchParamsModifier = func
    ctx.reload()
  }

  const initDropdown = ($dropdownOn: HTMLElement) => {
    renderDropdown({
      $dropdownWrap: $dropdownOn,
      dropdownList: [
        [
          $t('sortLatest'),
          () => {
            reloadUseParamsEditor((p) => {
              p.sort_by = 'date_desc'
            })
          },
        ],
        [
          $t('sortBest'),
          () => {
            reloadUseParamsEditor((p) => {
              p.sort_by = 'vote'
            })
          },
        ],
        [
          $t('sortOldest'),
          () => {
            reloadUseParamsEditor((p) => {
              p.sort_by = 'date_asc'
            })
          },
        ],
        [
          $t('sortAuthor'),
          () => {
            reloadUseParamsEditor((p) => {
              p.view_only_admin = true
            })
          },
        ],
      ],
    })
  }

  ctx.watchConf(['listSort', 'locale'], (conf) => {
    const list = ctx.get('list')

    const $count = list.$el.querySelector<HTMLElement>('.atk-comment-count')
    if (!$count) return

    // 评论列表排序 Dropdown 下拉选择层
    if (conf.listSort) {
      initDropdown($count)
    } else {
      removeDropdown({
        $dropdownWrap: $count,
      })
    }
  })
}

/** 评论排序方式选择下拉菜单 */
function renderDropdown(conf: {
  $dropdownWrap: HTMLElement
  dropdownList: [string, () => void][]
}) {
  const { $dropdownWrap, dropdownList } = conf
  if ($dropdownWrap.querySelector('.atk-dropdown')) return

  // 修改 class
  $dropdownWrap.classList.add('atk-dropdown-wrap')

  // 插入图标
  $dropdownWrap.append(Utils.createElement(`<span class="atk-arrow-down-icon"></span>`))

  // 列表项点击事件
  let curtActive = 0 // 当前选中
  const onItemClick = (i: number, $item: HTMLElement, name: string, action: Function) => {
    action()

    // set active
    curtActive = i
    $dropdown.querySelectorAll('.active').forEach((e) => {
      e.classList.remove('active')
    })
    $item.classList.add('active')

    // 关闭层 (临时消失，取消 :hover)
    $dropdown.style.display = 'none'
    setTimeout(() => {
      $dropdown.style.display = ''
    }, 80)
  }

  // 生成列表元素
  const $dropdown = Utils.createElement(`<ul class="atk-dropdown atk-fade-in"></ul>`)
  dropdownList.forEach((item, i) => {
    const name = item[0] as string
    const action = item[1] as Function

    const $item = Utils.createElement(`<li class="atk-dropdown-item"><span></span></li>`)
    const $link = $item.querySelector<HTMLElement>('span')!
    $link.innerText = name
    $link.onclick = () => {
      onItemClick(i, $item, name, action)
    }
    $dropdown.append($item)

    if (i === curtActive) $item.classList.add('active') // 默认选中项
  })

  $dropdownWrap.append($dropdown)
}

/** 删除评论排序方式选择下拉菜单 */
function removeDropdown(conf: { $dropdownWrap: HTMLElement }) {
  const { $dropdownWrap } = conf
  $dropdownWrap.classList.remove('atk-dropdown-wrap')
  $dropdownWrap.querySelector('.atk-arrow-down-icon')?.remove()
  $dropdownWrap.querySelector('.atk-dropdown')?.remove()
}
