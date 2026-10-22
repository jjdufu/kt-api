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
import { Activity, Gauge, KeyRound, LayoutDashboard } from 'lucide-react'
import { useTranslation } from 'react-i18next'

import { cn } from '@/lib/utils'

// Decorative sample data for the product mock (not live numbers).
const BARS = [38, 52, 44, 61, 58, 72, 66, 80, 74, 88, 83, 95].map((h, i) => ({
  id: `h${i}`,
  h,
  opacity: 0.45 + i * 0.045,
}))
const ROWS = [
  { model: 'gpt-4.1', latency: '412 ms' },
  { model: 'claude-sonnet-4', latency: '538 ms' },
  { model: 'gemini-2.5-pro', latency: '366 ms' },
  { model: 'deepseek-v3', latency: '291 ms' },
]

/**
 * Browser-framed product shot for the hero: a mini console (stat tiles,
 * request chart, recent calls) with a terminal card overlapping its corner.
 * Pure markup, themed by the preset tokens.
 */
export function ProductMock() {
  const { t } = useTranslation()
  const origin = typeof window === 'undefined' ? '' : window.location.host
  return (
    <div aria-hidden='true' className='relative select-none'>
      <div className='bg-card overflow-hidden rounded-xl border shadow-[0_1px_2px_rgb(0_0_0/0.04),0_24px_60px_-24px_rgb(0_0_0/0.28)]'>
        {/* window chrome */}
        <div className='bg-muted/40 flex h-10 items-center gap-3 border-b px-4'>
          <div className='flex gap-1.5'>
            <span className='size-2.5 rounded-full bg-(--window-dot-close)' />
            <span className='size-2.5 rounded-full bg-(--window-dot-min)' />
            <span className='size-2.5 rounded-full bg-(--window-dot-max)' />
          </div>
          <div className='bg-background text-muted-foreground mx-auto flex h-6 w-full max-w-64 items-center justify-center rounded-md border px-3 font-mono text-[11px]'>
            {origin || 'localhost'}/dashboard
          </div>
          <div className='w-10' />
        </div>
        <div className='flex'>
          {/* mini sidebar */}
          <div className='bg-muted/20 hidden w-40 shrink-0 flex-col gap-1 border-r p-3 text-xs sm:flex'>
            {[
              { icon: LayoutDashboard, label: t('Overview'), active: true },
              { icon: KeyRound, label: t('API Keys') },
              { icon: Activity, label: t('Usage') },
              { icon: Gauge, label: t('Latency') },
            ].map((item) => (
              <div
                key={item.label}
                className={cn(
                  'flex items-center gap-2 rounded-md px-2 py-1.5',
                  item.active
                    ? 'bg-background text-foreground border font-medium shadow-xs'
                    : 'text-muted-foreground'
                )}
              >
                <item.icon className='size-3.5' />
                <span className='truncate'>{item.label}</span>
              </div>
            ))}
          </div>
          {/* main */}
          <div className='min-w-0 flex-1 space-y-3 p-4'>
            <div className='grid grid-cols-3 gap-3'>
              {[
                { label: t('Requests'), value: '128.4k', delta: '+12.5%' },
                { label: t('Success'), value: '99.98%', delta: '+0.02%' },
                { label: t('Latency'), value: '312 ms', delta: '-8.1%' },
              ].map((s) => (
                <div key={s.label} className='rounded-lg border p-3'>
                  <div className='text-muted-foreground truncate text-[11px]'>
                    {s.label}
                  </div>
                  <div className='text-foreground mt-1 text-lg font-semibold tabular-nums'>
                    {s.value}
                  </div>
                  <div className='text-success mt-0.5 text-[10px] tabular-nums'>
                    {s.delta}
                  </div>
                </div>
              ))}
            </div>
            <div className='rounded-lg border p-3'>
              <div className='text-muted-foreground mb-2 flex items-center justify-between text-[11px]'>
                <span>{t('Requests')}</span>
                <span>{t('Today')}</span>
              </div>
              <div className='flex h-24 items-end gap-1.5'>
                {BARS.map(({ id, h, opacity }) => (
                  <div
                    key={id}
                    className='flex-1 rounded-sm bg-[linear-gradient(to_top,color-mix(in_oklch,var(--primary)_55%,transparent),var(--primary))]'
                    style={{ height: `${h}%`, opacity }}
                  />
                ))}
              </div>
            </div>
            <div className='overflow-hidden rounded-lg border text-[11px]'>
              <div className='bg-muted/40 text-muted-foreground grid grid-cols-[1fr_auto_auto] gap-4 border-b px-3 py-1.5'>
                <span>{t('Model')}</span>
                <span>{t('Latency')}</span>
                <span>{t('Status')}</span>
              </div>
              {ROWS.map((r) => (
                <div
                  key={r.model}
                  className='grid grid-cols-[1fr_auto_auto] items-center gap-4 border-b px-3 py-1.5 last:border-b-0'
                >
                  <span className='text-foreground truncate font-mono'>
                    {r.model}
                  </span>
                  <span className='text-muted-foreground font-mono tabular-nums'>
                    {r.latency}
                  </span>
                  <span className='bg-success/10 text-success rounded px-1.5 py-0.5 font-mono'>
                    200
                  </span>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
      {/* terminal card */}
      <div className='absolute -bottom-10 -left-6 hidden w-[22rem] overflow-hidden rounded-lg border border-white/10 bg-(--code-surface) font-mono text-[11px] leading-relaxed text-zinc-300 shadow-2xl md:block lg:-left-10'>
        <div className='flex items-center gap-1.5 border-b border-white/10 px-3 py-2'>
          <span className='size-2 rounded-full bg-white/20' />
          <span className='size-2 rounded-full bg-white/20' />
          <span className='size-2 rounded-full bg-white/20' />
          <span className='ms-2 text-zinc-500'>zsh</span>
        </div>
        <pre className='overflow-hidden px-3 py-2.5 whitespace-pre'>
          <span className='text-zinc-500'>$ </span>
          <span className='text-sky-300'>curl</span> {origin || 'localhost'}
          /v1/chat/completions \{'\n'} -H{' '}
          <span className='text-emerald-300'>
            "Authorization: Bearer sk-•••"
          </span>{' '}
          \{'\n'} -d{' '}
          <span className='text-emerald-300'>
            {'\'{"model":"gpt-4.1",...}\''}
          </span>
          {'\n'}
          <span className='text-zinc-500'>
            {'{'} "choices": [{'{'} "message": … {'}'}] {'}'}
          </span>
        </pre>
      </div>
    </div>
  )
}
