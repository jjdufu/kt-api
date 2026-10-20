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
import {
  Boxes,
  Coins,
  KeyRound,
  Plug,
  Route,
  Timer,
  type LucideIcon,
} from 'lucide-react'
import { useTranslation } from 'react-i18next'

import { SectionHeading } from './primitives'

type Capability = {
  no: string
  title: string
  description: string
  sfx: string
  icon: LucideIcon
  tone: 'cyan' | 'magenta' | 'yellow'
}

const CAPABILITIES: Capability[] = [
  {
    no: '#01',
    title: 'Unified API',
    description:
      'One OpenAI-compatible endpoint for chat, images, video and more.',
    sfx: 'ZAP!',
    icon: Plug,
    tone: 'cyan',
  },
  {
    no: '#02',
    title: 'Multi-model Routing',
    description:
      'Send each request to the right model and channel, with automatic failover.',
    sfx: 'WHOOSH!',
    icon: Route,
    tone: 'magenta',
  },
  {
    no: '#03',
    title: 'Pay as You Go',
    description: 'Pay only for what you use, with every call itemized.',
    sfx: 'KA-CHING!',
    icon: Coins,
    tone: 'yellow',
  },
  {
    no: '#04',
    title: 'Key Security',
    description:
      'Scoped keys, quotas and audit logs keep every credential in check.',
    sfx: 'CLANK!',
    icon: KeyRound,
    tone: 'cyan',
  },
  {
    no: '#05',
    title: 'Async Tasks',
    description:
      'Long-running jobs with polling and callbacks, no more timeouts.',
    sfx: 'TICK TOCK',
    icon: Timer,
    tone: 'magenta',
  },
  {
    no: '#06',
    title: 'Multi-protocol',
    description:
      'Speak OpenAI, Claude, Gemini and more protocols out of the box.',
    sfx: 'BOOM!',
    icon: Boxes,
    tone: 'yellow',
  },
]

export function CapabilityWall() {
  const { t } = useTranslation()
  return (
    <section className='hc-section'>
      <div className='hc-container'>
        <SectionHeading
          kicker={t('Core Capabilities')}
          title={t('Six panels of superpowers')}
        />
        <div className='hc-wall'>
          {CAPABILITIES.map((c, i) => {
            const Icon = c.icon
            return (
              <div key={c.no} className={`hc-wall-cell hc-wall-cell-${i + 1}`}>
                <article className={`hc-card hc-tone-${c.tone}`}>
                  <span className='hc-card-no'>{c.no}</span>
                  <span className='hc-card-sfx' aria-hidden='true'>
                    {c.sfx}
                  </span>
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
      </div>
    </section>
  )
}
