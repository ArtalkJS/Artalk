import React from 'react'
import { useTranslation } from 'react-i18next'
import { getFeatures, FeatureItem } from '../../features.ts'

export const FullList: React.FC = () => {
  const { t } = useTranslation()

  // let func list item group by every two items
  const FuncListGrouped = getFeatures(t).reduce<FeatureItem[][]>((result, current, index) => {
    if (index % 3 === 0) result.push([current])
    else result[result.length - 1].push(current)
    return result
  }, [])

  return (
    <div className="func-list">
      {FuncListGrouped.map((row, i) => (
        <div key={i} className="row">
          {row.map((item, j) => (
            <a key={j} className="item" href={item.link} target="_blank" rel="noreferrer">
              <div className="header">
                <div className="icon-wrap">
                  <item.icon />
                </div>
                <span className="text">{item.name}</span>
              </div>
              <div className="body">
                <div className="desc">{item.desc}</div>
              </div>
            </a>
          ))}
        </div>
      ))}
    </div>
  )
}
