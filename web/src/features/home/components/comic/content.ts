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
  LogIn,
  Plug,
  Route,
  Send,
  Timer,
  type LucideIcon,
} from 'lucide-react'

/* Shared section content (i18n keys) for the comic and enterprise homes. */

export type Capability = {
  no: string
  title: string
  description: string
  icon: LucideIcon
  tone: 'cyan' | 'magenta' | 'yellow'
}

export const CAPABILITIES: Capability[] = [
  {
    no: '#01',
    title: 'Unified API',
    description:
      'One OpenAI-compatible endpoint for chat, images, video and more.',
    icon: Plug,
    tone: 'cyan',
  },
  {
    no: '#02',
    title: 'Multi-model Routing',
    description:
      'Send each request to the right model and channel, with automatic failover.',
    icon: Route,
    tone: 'magenta',
  },
  {
    no: '#03',
    title: 'Pay as You Go',
    description: 'Pay only for what you use, with every call itemized.',
    icon: Coins,
    tone: 'yellow',
  },
  {
    no: '#04',
    title: 'Key Security',
    description:
      'Scoped keys, quotas and audit logs keep every credential in check.',
    icon: KeyRound,
    tone: 'cyan',
  },
  {
    no: '#05',
    title: 'Async Tasks',
    description:
      'Long-running jobs with polling and callbacks, no more timeouts.',
    icon: Timer,
    tone: 'magenta',
  },
  {
    no: '#06',
    title: 'Multi-protocol',
    description:
      'Speak OpenAI, Claude, Gemini and more protocols out of the box.',
    icon: Boxes,
    tone: 'yellow',
  },
]

export type Step = {
  title: string
  description: string
  icon: LucideIcon
  /** Optional in-app target; the step's bubble becomes a link. */
  to?: string
}

export const STEPS: Step[] = [
  {
    title: 'Sign In',
    description: 'Sign in to your account to get started.',
    icon: LogIn,
    to: '/sign-in',
  },
  {
    title: 'Get a Key',
    description: 'Generate an API key in the console.',
    icon: KeyRound,
  },
  {
    title: 'Send a Request',
    description: 'Set the gateway as your SDK base URL.',
    icon: Send,
  },
]
