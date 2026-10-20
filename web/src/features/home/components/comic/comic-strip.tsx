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
import { Link } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'

import { STEPS } from './content'
import { SpeechBubble } from './decorations'
import { SectionHeading } from './primitives'

export function ComicStrip() {
  const { t } = useTranslation()
  return (
    <section className='hc-section'>
      <div className='hc-container'>
        <SectionHeading
          kicker={t('Quick Start')}
          title={t('Make your first call in three steps')}
        />
        <ol className='hc-strip'>
          {STEPS.map((s, i) => {
            const Icon = s.icon
            return (
              <li key={s.title} className='hc-strip-item'>
                <div className='hc-strip-reveal'>
                  <div className='hc-strip-panel'>
                    <span className='hc-strip-no'>{i + 1}</span>
                    <span className='hc-strip-sfx' aria-hidden='true'>
                      {s.sfx}
                    </span>
                    <div className='hc-strip-art' aria-hidden='true'>
                      <Icon strokeWidth={2.4} />
                    </div>
                    <SpeechBubble className='hc-strip-bubble'>
                      {s.to ? (
                        <Link to={s.to} className='hc-strip-link'>
                          {t(s.title)}
                        </Link>
                      ) : (
                        t(s.title)
                      )}
                    </SpeechBubble>
                    <p className='hc-strip-desc'>{t(s.description)}</p>
                  </div>
                </div>
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
