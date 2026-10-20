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
import { useEffect, useRef } from 'react'
import { useTranslation } from 'react-i18next'

import { getDefaultStats } from '../../constants'
import { SectionHeading } from './primitives'

function CountUp({ end, suffix }: { end: number; suffix: string }) {
  const ref = useRef<HTMLSpanElement>(null)

  useEffect(() => {
    const el = ref.current
    if (!el) return
    const finalText = `${end}${suffix}`
    if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
      el.textContent = finalText
      return
    }
    let frame = 0
    const run = () => {
      const start = performance.now()
      const step = (now: number) => {
        const progress = Math.min((now - start) / 1400, 1)
        const eased = 1 - Math.pow(1 - progress, 3)
        el.textContent = `${Math.round(eased * end)}${suffix}`
        if (progress < 1) frame = requestAnimationFrame(step)
      }
      frame = requestAnimationFrame(step)
    }
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          observer.disconnect()
          run()
        }
      },
      { threshold: 0.4 }
    )
    observer.observe(el)
    return () => {
      observer.disconnect()
      cancelAnimationFrame(frame)
    }
  }, [end, suffix])

  return (
    <span ref={ref} className='tabular-nums'>
      {`0${suffix}`}
    </span>
  )
}

const TONES = ['cyan', 'magenta', 'yellow', 'cyan'] as const

export function NeonStats() {
  const { t } = useTranslation()
  const stats = getDefaultStats(t)
  return (
    <section className='hc-section hc-section-night'>
      <div className='hc-container'>
        <SectionHeading
          kicker={t('By the Numbers')}
          title={t('The gateway in neon')}
        />
        <div className='hc-neon-grid'>
          {stats.map((s, i) => (
            <div key={s.description}>
              <div
                className={`hc-neon-sign hc-neon-${TONES[i % TONES.length]}`}
              >
                <span className='hc-neon-value'>
                  <CountUp end={Number(s.value)} suffix={s.suffix} />
                </span>
                <span className='hc-neon-label'>{s.description}</span>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}
