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

import { getLobeIcon } from '@/lib/lobe-icon'

import { AI_APPLICATIONS, AI_MODELS } from '../../constants'
import { SectionHeading } from './primitives'

const ICONS = [...AI_MODELS, ...AI_APPLICATIONS]

/**
 * The track holds two identical reels and slides by exactly one reel width
 * (-50%), so the loop is seamless only while one reel is wider than the film
 * window. Ten logos are ~1500px, which left a blank tail on wide screens;
 * repeating them inside each reel keeps a reel ~4500px wide (4K-safe).
 */
const REEL_ROUNDS = ['r1', 'r2', 'r3'] as const
const REEL_ITEMS = REEL_ROUNDS.flatMap((round) =>
  ICONS.map((name) => ({ id: `${round}-${name}`, name }))
)

function Reel({ copy }: { copy: 'a' | 'b' }) {
  return (
    <ul className='hc-reel' aria-hidden={copy === 'b'}>
      {REEL_ITEMS.map(({ id, name }) => (
        <li key={`${copy}-${id}`} className='hc-frame'>
          {getLobeIcon(name, 40)}
          <span className='hc-frame-label'>{name.split('.')[0]}</span>
        </li>
      ))}
    </ul>
  )
}

/** Model logos running through a film strip. */
export function ModelMarquee() {
  const { t } = useTranslation()
  return (
    <section className='hc-section'>
      <div className='hc-container'>
        <SectionHeading
          kicker={t('Model Lineup')}
          title={t('Plugged into the models you love')}
        />
      </div>
      <div className='hc-film'>
        <div className='hc-film-track'>
          <Reel copy='a' />
          <Reel copy='b' />
        </div>
      </div>
    </section>
  )
}
