import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import LazyImage from '../LazyImage'
import './StartCards.scss'

interface ICard {
  username: string
  avatar: string
  text: string
}

export const StartCards: React.FC = () => {
  const { t } = useTranslation()
  const [cards, setCards] = useState<ICard[]>([])

  useEffect(() => {
    setCards([
      {
        username: t('intro_stack_card_1_name'),
        avatar: '/images/avatar-1.webp',
        text: t('intro_stack_card_1_text'),
      },
      {
        username: t('intro_stack_card_2_name'),
        avatar: '/images/avatar-2.webp',
        text: t('intro_stack_card_2_text'),
      },
      {
        username: t('intro_stack_card_3_name'),
        avatar: '/images/avatar-3.webp',
        text: t('intro_stack_card_3_text'),
      },
      {
        username: '',
        avatar: '',
        text: '',
      },
    ])
  }, [t, setCards])

  return (
    <div className="start-cards">
      {cards.map((card, index) => (
        <div key={index} className="card">
          {card.username && (
            <div className="header">
              <div className="avatar">
                <LazyImage src={card.avatar} alt="avatar" referrerPolicy="no-referrer" />
              </div>
              <div className="username">{card.username}</div>
            </div>
          )}
          {card.text && <div className="content">{card.text}</div>}
          {!card.username && (
            <div className="placeholder">
              <div className="text-row"></div>
              <div className="text-row"></div>
              <div className="submit-button">
                <div className="text"></div>
              </div>
            </div>
          )}
        </div>
      ))}
    </div>
  )
}
