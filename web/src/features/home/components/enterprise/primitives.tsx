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
import { cn } from '@/lib/utils'

/**
 * Building blocks for the enterprise home / about pages used by every preset
 * except KT Comic: one content width, quiet type scale, 1px borders and a
 * single accent (the preset's --primary).
 */
export function Container(props: { children: ReactNode; className?: string }) {
  return (
    <div
      className={cn(
        'mx-auto w-full max-w-6xl px-4 sm:px-6 lg:px-8',
        props.className
      )}
    >
      {props.children}
    </div>
  )
}

export function Section(props: {
  children: ReactNode
  className?: string
  id?: string
}) {
  return (
    <section id={props.id} className={cn('py-20 sm:py-24', props.className)}>
      <Container>{props.children}</Container>
    </section>
  )
}

export function Kicker(props: { children: ReactNode; className?: string }) {
  return (
    <p
      className={cn(
        'text-primary text-sm font-medium tracking-wide',
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

/** Soft radial glow + a dot matrix that fades out towards the edges (hero only). */
export function HeroBackdrop() {
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

/** Primary / outline call-to-action that handles external URLs. */
export function CtaLink(props: {
  to: string
  children: ReactNode
  variant?: 'default' | 'outline'
  className?: string
}) {
  const className = cn(
    buttonVariants({ variant: props.variant ?? 'default', size: 'lg' }),
    'h-11 rounded-lg px-5 text-sm font-medium',
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
export function Surface(props: { children: ReactNode; className?: string }) {
  return (
    <div className={cn('bg-background/60 rounded-lg border', props.className)}>
      {props.children}
    </div>
  )
}
