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
import type { ReactNode } from 'react'

import { buttonVariants } from '@/components/ui/button'
import { useThemeCustomization } from '@/context/theme-customization-provider'
import { cn } from '@/lib/utils'

import { HeroFx } from './comic-decorations'

/**
 * Shared building blocks for the marketing pages.
 *
 * Home and about are composed from these primitives; the docs page uses
 * Container, Kicker and HeroBackdrop for its outer shell only and keeps its
 * own reference pieces (TOC, tables, endpoint cards, code samples; `dc-*`).
 *
 * One component set serves both presets. The KT Comic look lives in CSS
 * (`hc-*` classes scoped to `[data-home]` / `[data-docs]`); every other
 * preset gets the quiet enterprise layout. Pages never branch on the preset
 * for these primitives — only for composition-level decoration.
 */
export type MarketingTone = 'comic' | 'plain'

export function useMarketingTone(): MarketingTone {
  return useThemeCustomization().customization.preset === 'kt-comic'
    ? 'comic'
    : 'plain'
}

type Base = { children?: ReactNode; className?: string }

export function Container(props: Base) {
  const tone = useMarketingTone()
  return (
    <div
      className={cn(
        tone === 'comic'
          ? 'hc-container'
          : 'mx-auto w-full max-w-6xl px-4 sm:px-6 lg:px-8',
        props.className
      )}
    >
      {props.children}
    </div>
  )
}

export function Section(
  props: Base & { id?: string; containerClassName?: string }
) {
  const tone = useMarketingTone()
  return (
    <section
      id={props.id}
      className={cn(
        tone === 'comic' ? 'hc-section' : 'py-20 sm:py-24',
        props.className
      )}
    >
      <Container className={props.containerClassName}>
        {props.children}
      </Container>
    </section>
  )
}

export function Kicker(props: Base) {
  const tone = useMarketingTone()
  return (
    <p
      className={cn(
        tone === 'comic'
          ? 'hc-kicker'
          : 'text-primary text-sm font-medium tracking-wide',
        props.className
      )}
    >
      {props.children}
    </p>
  )
}

export function SectionHeading(props: {
  kicker: string
  title: string
  description?: string
  align?: 'left' | 'center'
  className?: string
}) {
  const tone = useMarketingTone()
  if (tone === 'comic') {
    return (
      <div className={cn('hc-heading', props.className)}>
        <span className='hc-kicker'>{props.kicker}</span>
        <h2 className='hc-title'>{props.title}</h2>
        {props.description && (
          <p className='hc-heading-desc'>{props.description}</p>
        )}
      </div>
    )
  }
  const center = props.align !== 'left'
  return (
    <div
      className={cn(
        'mb-12 max-w-2xl',
        center && 'mx-auto text-center',
        props.className
      )}
    >
      <Kicker>{props.kicker}</Kicker>
      <h2 className='text-foreground mt-3 text-3xl font-semibold tracking-tight text-balance sm:text-4xl'>
        {props.title}
      </h2>
      {props.description && (
        <p className='text-muted-foreground mt-4 text-base leading-relaxed text-pretty sm:text-lg'>
          {props.description}
        </p>
      )}
    </div>
  )
}

/** Primary / outline call-to-action that handles external URLs. */
export function CtaLink(props: {
  to: string
  children: ReactNode
  variant?: 'primary' | 'outline'
  className?: string
}) {
  const tone = useMarketingTone()
  const variant = props.variant ?? 'primary'
  const className = cn(
    tone === 'comic'
      ? ['hc-btn', `hc-btn-${variant}`]
      : [
          buttonVariants({
            variant: variant === 'primary' ? 'default' : 'outline',
            size: 'lg',
          }),
          'h-11 rounded-lg px-5 text-sm font-medium',
        ],
    props.className
  )
  if (/^https?:\/\//.test(props.to)) {
    return (
      <a
        href={props.to}
        target='_blank'
        rel='noopener noreferrer'
        className={className}
      >
        {props.children}
      </a>
    )
  }
  return (
    <Link to={props.to} className={className}>
      {props.children}
    </Link>
  )
}

/** Small bordered surface used by the bento previews and mocks. */
export function Surface(props: Base) {
  return (
    <div className={cn('bg-background/60 rounded-lg border', props.className)}>
      {props.children}
    </div>
  )
}

/**
 * Hero backdrop. Plain: soft radial glow + fading dot matrix. Comic: the
 * neon / skyline / scanline layer styled by page-backdrop.css.
 */
export function HeroBackdrop() {
  const tone = useMarketingTone()
  if (tone === 'comic') return <HeroFx />
  return (
    <div
      aria-hidden='true'
      className='pointer-events-none absolute inset-0 -z-10 overflow-hidden'
    >
      <div className='absolute top-[-14rem] left-1/2 h-[38rem] w-[min(72rem,140vw)] -translate-x-1/2 rounded-full bg-[radial-gradient(closest-side,color-mix(in_oklch,var(--primary)_18%,transparent),transparent)]' />
      <div className='absolute inset-0 [background-image:radial-gradient(color-mix(in_oklch,var(--foreground)_14%,transparent)_1px,transparent_1.4px)] [mask-image:radial-gradient(ellipse_75%_65%_at_50%_0%,#000_30%,transparent_80%)] [background-size:22px_22px]' />
    </div>
  )
}
