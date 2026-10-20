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
import { cn } from '@/lib/utils'

export type SiteHeaderVariant = 'console' | 'public'

type SiteHeaderProps = React.HTMLAttributes<HTMLElement> & {
  variant: SiteHeaderVariant
}

/**
 * The one top bar used by every page (console and public). Same height
 * (--app-header-height, 3.5rem), gutters, surface and border everywhere;
 * themes style it via [data-slot='site-header'] (variant in data-variant).
 * Sticky, so pages never need their own top offset.
 */
export function SiteHeader({
  variant,
  className,
  children,
  ...props
}: SiteHeaderProps) {
  return (
    <header
      data-slot='site-header'
      data-variant={variant}
      className={cn(
        'bg-background/85 supports-[backdrop-filter]:bg-background/70 sticky top-0 z-50 h-(--app-header-height) w-full shrink-0 border-b backdrop-blur-xl',
        className
      )}
      {...props}
    >
      <div className='flex h-full items-center gap-2 px-3 sm:px-4'>
        {children}
      </div>
    </header>
  )
}

/**
 * Shared active rule for header nav links: exact match for `/`, otherwise
 * the path itself or any sub-path.
 */
export function isNavLinkActive(pathname: string, href: string) {
  if (/^https?:\/\//.test(href)) return false
  const path = href.split(/[?#]/)[0] || '/'
  if (path === '/') return pathname === '/'
  return pathname === path || pathname.startsWith(`${path}/`)
}
