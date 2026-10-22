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
import { useTranslation } from 'react-i18next'

import { Section, SectionHeading } from '@/features/marketing/primitives'

import { CAPABILITIES } from './content'

export function CapabilityWall() {
  const { t } = useTranslation()
  return (
    <Section>
      <SectionHeading
        kicker={t('Core Capabilities')}
        title={t('Everything you need to run AI in production')}
      />
      <div className='hc-wall'>
        {CAPABILITIES.map((c, i) => {
          const Icon = c.icon
          return (
            <div key={c.no} className={`hc-wall-cell hc-wall-cell-${i + 1}`}>
              <article className={`hc-card hc-tone-${c.tone}`}>
                <span className='hc-card-no'>{c.no}</span>
                <div className='hc-card-icon' aria-hidden='true'>
                  <Icon strokeWidth={2.5} />
                </div>
                <h3 className='hc-card-title'>{t(c.title)}</h3>
                <p className='hc-card-desc'>{t(c.description)}</p>
              </article>
            </div>
          )
        })}
      </div>
    </Section>
  )
}
