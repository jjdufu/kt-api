import { useReducedMotion } from 'motion/react'
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
import { useEffect, useState } from 'react'

const MODELS = ['gpt', 'claude', 'deepseek', 'qwen', 'gemini', 'doubao']
const LATENCY: Record<string, number> = {
  gpt: 182,
  claude: 204,
  deepseek: 156,
  qwen: 131,
  gemini: 175,
  doubao: 118,
}

const TYPE_MS = 90
const DELETE_MS = 45
const HOLD_MS = 1600

function apiOrigin(): string {
  if (typeof window === 'undefined') return 'https://api.example.com'
  return window.location.origin
}

/**
 * Compact green-phosphor terminal that keeps retyping the model name, so the
 * hero shows one endpoint serving many models.
 */
export function ComicTerminal() {
  const reduce = useReducedMotion()
  const [modelIndex, setModelIndex] = useState(0)
  const [typed, setTyped] = useState(reduce ? MODELS[0] : '')
  const [phase, setPhase] = useState<'typing' | 'holding' | 'deleting'>(
    reduce ? 'holding' : 'typing'
  )

  const model = MODELS[modelIndex]

  useEffect(() => {
    if (reduce) return
    let timer: number
    if (phase === 'typing') {
      if (typed.length < model.length) {
        timer = window.setTimeout(
          () => setTyped(model.slice(0, typed.length + 1)),
          TYPE_MS
        )
      } else {
        timer = window.setTimeout(() => setPhase('holding'), 0)
      }
    } else if (phase === 'holding') {
      timer = window.setTimeout(() => setPhase('deleting'), HOLD_MS)
    } else if (typed.length > 0) {
      timer = window.setTimeout(() => setTyped(typed.slice(0, -1)), DELETE_MS)
    } else {
      timer = window.setTimeout(() => {
        setModelIndex((i) => (i + 1) % MODELS.length)
        setPhase('typing')
      }, 120)
    }
    return () => window.clearTimeout(timer)
  }, [model, phase, reduce, typed])

  const done = phase === 'holding' && typed === model

  return (
    <div className='hc-terminal' aria-label='curl /v1/chat/completions'>
      <div className='hc-terminal-bar'>
        <span className='hc-terminal-prompt'>&gt;_</span>
        <span className='hc-terminal-dots' aria-hidden='true'>
          <i />
          <i />
          <i />
        </span>
      </div>
      <pre className='hc-terminal-body'>
        <span className='hc-t-cmd'>curl</span>{' '}
        <span className='hc-t-url'>{apiOrigin()}/v1/chat/completions</span>
        {' \\\n  '}
        <span className='hc-t-flag'>-H</span>{' '}
        <span className='hc-t-str'>
          &quot;Authorization: Bearer sk-••••&quot;
        </span>
        {' \\\n  '}
        <span className='hc-t-flag'>-d</span>{' '}
        <span className='hc-t-str'>
          {'\'{"model": "'}
          <span className='hc-t-model'>{typed}</span>
          <span className='hc-t-caret' aria-hidden='true' />
          {'", ...}\''}
        </span>
        {'\n'}
        <span className={done ? 'hc-t-ok hc-t-ok-on' : 'hc-t-ok'}>
          {`← 200 OK · ${LATENCY[model]}ms · model: ${model}`}
        </span>
      </pre>
    </div>
  )
}
