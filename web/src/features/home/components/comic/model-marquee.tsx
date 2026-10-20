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

function Reel({ copy }: { copy: 'a' | 'b' }) {
  return (
    <ul className='hc-reel' aria-hidden={copy === 'b'}>
      {ICONS.map((name) => (
        <li key={`${copy}-${name}`} className='hc-frame'>
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
