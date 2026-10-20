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
import { useId, type ReactNode } from 'react'

import { cn } from '@/lib/utils'

/**
 * Decorative SVG/CSS pieces for the comic home page. Everything here is
 * aria-hidden and purely visual.
 */

type BurstProps = {
  label: string
  className?: string
  points?: number
}

/** Spiky comic burst ("ZAP!") with a halftone fill and ink outline. */
export function Burst({ label, className, points = 14 }: BurstProps) {
  const id = useId().replaceAll(':', '')
  const outer = 50
  const inner = 33
  const coords: string[] = []
  for (let i = 0; i < points * 2; i += 1) {
    const r = i % 2 === 0 ? outer : inner + ((i * 7) % 5)
    const a = (Math.PI * i) / points - Math.PI / 2
    coords.push(
      `${(60 + r * Math.cos(a)).toFixed(1)},${(55 + r * Math.sin(a) * 0.9).toFixed(1)}`
    )
  }
  return (
    <svg
      viewBox='0 0 120 110'
      className={cn('hc-burst', className)}
      aria-hidden='true'
    >
      <defs>
        <pattern
          id={`dots-${id}`}
          width='5'
          height='5'
          patternUnits='userSpaceOnUse'
        >
          <circle cx='2.5' cy='2.5' r='1.1' fill='rgba(0,0,0,0.28)' />
        </pattern>
      </defs>
      <polygon
        points={coords.join(' ')}
        fill='var(--hc-yellow)'
        stroke='#000'
        strokeWidth='4'
        strokeLinejoin='round'
      />
      <polygon points={coords.join(' ')} fill={`url(#dots-${id})`} />
      <text
        x='60'
        y='66'
        textAnchor='middle'
        className='hc-burst-text'
        transform='rotate(-10 60 55)'
      >
        {label}
      </text>
    </svg>
  )
}

type Building = {
  x: number
  w: number
  h: number
  tone: 'cyan' | 'magenta'
  windows: Array<{ x: number; y: number; on: boolean }>
}

function buildSkyline(): Building[] {
  // Deterministic pseudo-random so SSR/hydration and screenshots match.
  let seed = 7
  const rand = () => {
    seed = (seed * 9301 + 49297) % 233280
    return seed / 233280
  }
  const buildings: Building[] = []
  let x = 0
  while (x < 1200) {
    const w = 28 + Math.floor(rand() * 46)
    const h = 70 + Math.floor(rand() * 230)
    const tone = rand() > 0.5 ? 'cyan' : 'magenta'
    const windows: Building['windows'] = []
    for (let wy = 300 - h + 14; wy < 290; wy += 14) {
      for (let wx = x + 6; wx < x + w - 6; wx += 9) {
        windows.push({ x: wx, y: wy, on: rand() > 0.62 })
      }
    }
    buildings.push({ x, w, h, tone, windows })
    x += w + 4 + Math.floor(rand() * 10)
  }
  return buildings
}

const SKYLINE = buildSkyline()

/** Neon city skyline silhouette (pure SVG). */
export function Skyline({ className }: { className?: string }) {
  return (
    <svg
      viewBox='0 0 1200 300'
      preserveAspectRatio='xMidYMax slice'
      className={cn('hc-skyline', className)}
      aria-hidden='true'
    >
      {SKYLINE.map((b) => (
        <g key={`b-${b.x}`}>
          <rect
            x={b.x}
            y={300 - b.h}
            width={b.w}
            height={b.h}
            className='hc-skyline-body'
          />
          <rect
            x={b.x}
            y={300 - b.h}
            width={b.w}
            height='2'
            className={b.tone === 'cyan' ? 'hc-neon-cyan' : 'hc-neon-magenta'}
          />
          {b.windows
            .filter((w) => w.on)
            .map((w) => (
              <rect
                key={`w-${w.x}-${w.y}`}
                x={w.x}
                y={w.y}
                width='4'
                height='6'
                className={
                  b.tone === 'cyan' ? 'hc-window-cyan' : 'hc-window-magenta'
                }
              />
            ))}
        </g>
      ))}
    </svg>
  )
}

/** Halftone dots + radial speed lines background layer. */
export function ComicBackdrop({ className }: { className?: string }) {
  return (
    <div className={cn('hc-backdrop', className)} aria-hidden='true'>
      <div className='hc-speedlines' />
      <div className='hc-halftone' />
    </div>
  )
}

type BubbleProps = {
  children: ReactNode
  className?: string
  tail?: 'left' | 'right'
}

/** Comic speech bubble with an ink outline and tail. */
export function SpeechBubble({
  children,
  className,
  tail = 'left',
}: BubbleProps) {
  return (
    <div className={cn('hc-bubble', `hc-bubble-${tail}`, className)}>
      {children}
    </div>
  )
}
