/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/
import { KeyRound, Send, UserPlus, type LucideIcon } from 'lucide-react'
import { useTranslation } from 'react-i18next'

import { SpeechBubble } from './decorations'
import { Reveal, SectionHeading } from './primitives'

type Step = {
  title: string
  description: string
  icon: LucideIcon
  sfx: string
}

const STEPS: Step[] = [
  {
    title: 'Sign Up',
    description: 'Create your account in seconds.',
    icon: UserPlus,
    sfx: 'POP!',
  },
  {
    title: 'Get a Key',
    description: 'Generate an API key in the console.',
    icon: KeyRound,
    sfx: 'CLICK!',
  },
  {
    title: 'Send a Request',
    description: 'Point your SDK at the gateway and go.',
    icon: Send,
    sfx: 'ZOOM!',
  },
]

export function ComicStrip() {
  const { t } = useTranslation()
  return (
    <section className='hc-section'>
      <div className='hc-container'>
        <SectionHeading
          kicker={t('Three Panels')}
          title={t('Your first call in three panels')}
        />
        <ol className='hc-strip'>
          {STEPS.map((s, i) => {
            const Icon = s.icon
            return (
              <li key={s.title} className='hc-strip-item'>
                <Reveal className='hc-strip-reveal' delay={i * 0.12} tilt={0}>
                  <div className='hc-strip-panel'>
                    <span className='hc-strip-no'>{i + 1}</span>
                    <span className='hc-strip-sfx' aria-hidden='true'>
                      {s.sfx}
                    </span>
                    <div className='hc-strip-art' aria-hidden='true'>
                      <Icon strokeWidth={2.4} />
                    </div>
                    <SpeechBubble className='hc-strip-bubble'>
                      {t(s.title)}
                    </SpeechBubble>
                    <p className='hc-strip-desc'>{t(s.description)}</p>
                  </div>
                </Reveal>
                {i < STEPS.length - 1 && (
                  <svg
                    className='hc-strip-arrow'
                    viewBox='0 0 48 32'
                    aria-hidden='true'
                  >
                    <path d='M2 10h28V2l16 14-16 14v-8H2z' />
                  </svg>
                )}
              </li>
            )
          })}
        </ol>
      </div>
    </section>
  )
}
