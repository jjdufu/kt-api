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
import { toast } from 'sonner'

/** Success toasts are confirmations; they leave quickly. */
export const SUCCESS_TOAST_DURATION = 2000

let installed = false

/**
 * Give every `toast.success(...)` a short default duration. Other types
 * keep the <Toaster> default; an explicit `duration` always wins.
 */
export function installToastDefaults() {
  if (installed) return
  installed = true
  const success = toast.success
  toast.success = (message, data) =>
    success(message, { duration: SUCCESS_TOAST_DURATION, ...data })
}

installToastDefaults()
