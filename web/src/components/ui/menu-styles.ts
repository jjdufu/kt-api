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
/**
 * Shared class hooks for every dropdown-style popup (DropdownMenu, Select,
 * Combobox, Command, Popover). The look itself lives in
 * src/styles/ui-bridge.css and is driven by the --surface-popup-* and
 * --item-* theme tokens (src/styles/tokens), so these upstream shadcn files
 * only carry one extra class each and stay easy to sync with new-api.
 */
export const menuPopupClass = 'ui-menu-popup'
export const menuItemClass = 'ui-menu-item'
export const menuLabelClass = 'ui-menu-label'
export const menuSeparatorClass = 'ui-menu-separator'
