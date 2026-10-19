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
import { zodResolver } from '@hookform/resolvers/zod'
import { useEffect, useMemo, useRef } from 'react'
import { useForm } from 'react-hook-form'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import * as z from 'zod'

import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Separator } from '@/components/ui/separator'
import { Switch } from '@/components/ui/switch'
import { Textarea } from '@/components/ui/textarea'

import {
  SettingsForm,
  SettingsFormGrid,
  SettingsSwitchContent,
  SettingsSwitchItem,
} from '../components/settings-form-layout'
import { SettingsPageFormActions } from '../components/settings-page-context'
import { SettingsSection } from '../components/settings-section'
import { useUpdateOption } from '../hooks/use-update-option'
import { safeNumberFieldProps } from '../utils/numeric-field'
import {
  TASK_CALLBACK_BOOLEAN_FIELDS,
  TASK_CALLBACK_INJECT_MODE_LABEL,
  TASK_CALLBACK_INJECT_MODES,
  TASK_CALLBACK_NUMBER_FIELDS,
  type TaskCallbackSettingValues,
} from './constants'
import { TaskCallbackDeliveries } from './task-callback-deliveries'

const taskCallbackSchema = z.object({
  enabled: z.boolean(),
  extra_headers_enabled: z.boolean(),
  https_only: z.boolean(),
  allow_private_ip: z.boolean(),
  domain_allowlist_mode: z.boolean(),
  reject_unresolvable: z.boolean(),
  extended_retry: z.boolean(),
  signing_enabled: z.boolean(),
  domain_list: z.string(),
  allowed_ports: z.string(),
  max_url_length: z.number().int().min(64).max(8192),
  timeout_seconds: z.number().int().min(1).max(30),
  retry_interval_seconds: z.number().int().min(1).max(60),
  terminal_max_attempts: z.number().int().min(1).max(10),
  workers: z.number().int().min(1).max(128),
  per_host_concurrency: z.number().int().min(1).max(64),
  retention_days: z.number().int().min(1).max(90),
  signing_secret: z.string(),
  upstream_public_base: z
    .string()
    .trim()
    .refine((value) => value === '' || /^https?:\/\/[^\s/]+/i.test(value), {
      message: 'Enter an http(s) URL or leave it empty',
    }),
  upstream_inject_mode: z.enum(TASK_CALLBACK_INJECT_MODES),
})

type TaskCallbackFormValues = z.infer<typeof taskCallbackSchema>

type TaskCallbackSettingsSectionProps = {
  defaultValues: TaskCallbackSettingValues
}

const PREFIX = 'task_callback_setting.'

const splitList = (value: string) =>
  value
    .split(/[\n,]/)
    .map((entry) => entry.trim())
    .filter(Boolean)

function toFormValues(
  defaults: TaskCallbackSettingValues
): TaskCallbackFormValues {
  return {
    enabled: defaults[`${PREFIX}enabled`],
    extra_headers_enabled: defaults[`${PREFIX}extra_headers_enabled`],
    https_only: defaults[`${PREFIX}https_only`],
    allow_private_ip: defaults[`${PREFIX}allow_private_ip`],
    domain_allowlist_mode: defaults[`${PREFIX}domain_allowlist_mode`],
    reject_unresolvable: defaults[`${PREFIX}reject_unresolvable`],
    extended_retry: defaults[`${PREFIX}extended_retry`],
    signing_enabled: defaults[`${PREFIX}signing_enabled`],
    domain_list: defaults[`${PREFIX}domain_list`].join('\n'),
    allowed_ports: defaults[`${PREFIX}allowed_ports`].join(','),
    max_url_length: defaults[`${PREFIX}max_url_length`],
    timeout_seconds: defaults[`${PREFIX}timeout_seconds`],
    retry_interval_seconds: defaults[`${PREFIX}retry_interval_seconds`],
    terminal_max_attempts: defaults[`${PREFIX}terminal_max_attempts`],
    workers: defaults[`${PREFIX}workers`],
    per_host_concurrency: defaults[`${PREFIX}per_host_concurrency`],
    retention_days: defaults[`${PREFIX}retention_days`],
    signing_secret: '',
    upstream_public_base: defaults[`${PREFIX}upstream_public_base`],
    upstream_inject_mode: defaults[`${PREFIX}upstream_inject_mode`],
  }
}

// Option updates (key -> serialized value) that differ from the saved form.
// The signing secret is write-only: it is sent only when a new one is typed.
function collectUpdates(
  values: TaskCallbackFormValues,
  baseline: TaskCallbackFormValues
): Array<{ key: string; value: string | number | boolean }> {
  const updates: Array<{ key: string; value: string | number | boolean }> = []
  for (const field of TASK_CALLBACK_BOOLEAN_FIELDS) {
    if (values[field.name] !== baseline[field.name]) {
      updates.push({ key: PREFIX + field.name, value: values[field.name] })
    }
  }
  for (const field of TASK_CALLBACK_NUMBER_FIELDS) {
    if (values[field.name] !== baseline[field.name]) {
      updates.push({ key: PREFIX + field.name, value: values[field.name] })
    }
  }
  const lists = ['domain_list', 'allowed_ports'] as const
  for (const name of lists) {
    const next = JSON.stringify(splitList(values[name]))
    if (next !== JSON.stringify(splitList(baseline[name]))) {
      updates.push({ key: PREFIX + name, value: next })
    }
  }
  const base = values.upstream_public_base.trim()
  if (base !== baseline.upstream_public_base.trim()) {
    updates.push({ key: `${PREFIX}upstream_public_base`, value: base })
  }
  if (values.upstream_inject_mode !== baseline.upstream_inject_mode) {
    updates.push({
      key: `${PREFIX}upstream_inject_mode`,
      value: values.upstream_inject_mode,
    })
  }
  if (values.signing_secret.trim() !== '') {
    updates.push({
      key: `${PREFIX}signing_secret`,
      value: values.signing_secret.trim(),
    })
  }
  return updates
}

export function TaskCallbackSettingsSection(
  props: TaskCallbackSettingsSectionProps
) {
  const { t } = useTranslation()
  const updateOption = useUpdateOption()
  const formDefaults = useMemo(
    () => toFormValues(props.defaultValues),
    [props.defaultValues]
  )
  const baselineRef = useRef<TaskCallbackFormValues>(formDefaults)
  const form = useForm<TaskCallbackFormValues>({
    resolver: zodResolver(taskCallbackSchema),
    defaultValues: formDefaults,
  })

  useEffect(() => {
    baselineRef.current = formDefaults
    form.reset(formDefaults)
  }, [formDefaults, form])

  const onSubmit = async (values: TaskCallbackFormValues) => {
    const updates = collectUpdates(values, baselineRef.current)
    if (updates.length === 0) {
      toast.info(t('No changes to save'))
      return
    }
    for (const update of updates) {
      await updateOption.mutateAsync(update)
    }
    baselineRef.current = { ...values, signing_secret: '' }
    form.setValue('signing_secret', '')
  }

  return (
    <SettingsSection title={t('Task callbacks')} className='gap-6'>
      <Form {...form}>
        <SettingsForm onSubmit={form.handleSubmit(onSubmit)}>
          <SettingsPageFormActions
            onSave={form.handleSubmit(onSubmit)}
            isSaving={updateOption.isPending}
            saveLabel='Save task callback settings'
          />
          <p className='text-muted-foreground text-sm'>
            {t(
              'Seedance video tasks submitted with callback_url receive a POST with the same body as the task query whenever the task status changes. Polling stays the source of truth.'
            )}
          </p>
          <SettingsFormGrid>
            {TASK_CALLBACK_BOOLEAN_FIELDS.map((item) => (
              <FormField
                key={item.name}
                control={form.control}
                name={item.name}
                render={({ field }) => (
                  <SettingsSwitchItem>
                    <SettingsSwitchContent>
                      <FormLabel>{t(item.label)}</FormLabel>
                      <FormDescription>{t(item.description)}</FormDescription>
                    </SettingsSwitchContent>
                    <FormControl>
                      <Switch
                        checked={field.value}
                        onCheckedChange={field.onChange}
                      />
                    </FormControl>
                  </SettingsSwitchItem>
                )}
              />
            ))}
          </SettingsFormGrid>

          <Separator />

          <SettingsFormGrid>
            <FormField
              control={form.control}
              name='domain_list'
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('Callback domain allowlist')}</FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder='example.com&#10;*.example.com'
                      rows={4}
                      {...field}
                    />
                  </FormControl>
                  <FormDescription>
                    {t(
                      'One domain per line. Used only when the allowlist mode is on.'
                    )}
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='allowed_ports'
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('Allowed callback ports')}</FormLabel>
                  <FormControl>
                    <Input placeholder='80,443,8080,8443' {...field} />
                  </FormControl>
                  <FormDescription>
                    {t(
                      'Comma-separated ports or ranges such as 8000-9000 (empty = all ports)'
                    )}
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            {TASK_CALLBACK_NUMBER_FIELDS.map((item) => (
              <FormField
                key={item.name}
                control={form.control}
                name={item.name}
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t(item.label)}</FormLabel>
                    <FormControl>
                      <Input
                        type='number'
                        min={item.min}
                        max={item.max}
                        step={1}
                        {...safeNumberFieldProps(field)}
                      />
                    </FormControl>
                    <FormDescription>{t(item.description)}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            ))}
            <FormField
              control={form.control}
              name='signing_secret'
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('Signing master secret')}</FormLabel>
                  <FormControl>
                    <Input
                      type='password'
                      autoComplete='new-password'
                      {...field}
                    />
                  </FormControl>
                  <FormDescription>
                    {t(
                      'Write-only. Leave empty to keep the current secret. Each user gets a secret derived from it.'
                    )}
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
          </SettingsFormGrid>

          <Separator />

          <SettingsFormGrid>
            <FormField
              control={form.control}
              name='upstream_public_base'
              render={({ field }) => (
                <FormItem>
                  <FormLabel>
                    {t('Upstream callback public base URL')}
                  </FormLabel>
                  <FormControl>
                    <Input placeholder='https://api.example.com' {...field} />
                  </FormControl>
                  <FormDescription>
                    {t(
                      'Optional fast path: the supplier calls this kt-api address when a task changes. Leave empty to rely on polling only. Channels must also enable it.'
                    )}
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='upstream_inject_mode'
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('Upstream callback mode')}</FormLabel>
                  <Select
                    items={TASK_CALLBACK_INJECT_MODES.map((mode) => ({
                      value: mode,
                      label: t(TASK_CALLBACK_INJECT_MODE_LABEL[mode]),
                    }))}
                    onValueChange={field.onChange}
                    value={field.value}
                  >
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent alignItemWithTrigger={false}>
                      <SelectGroup>
                        {TASK_CALLBACK_INJECT_MODES.map((mode) => (
                          <SelectItem key={mode} value={mode}>
                            {t(TASK_CALLBACK_INJECT_MODE_LABEL[mode])}
                          </SelectItem>
                        ))}
                      </SelectGroup>
                    </SelectContent>
                  </Select>
                  <FormDescription>
                    {t('When to send the fast-path URL to the supplier')}
                  </FormDescription>
                </FormItem>
              )}
            />
          </SettingsFormGrid>
        </SettingsForm>
      </Form>

      <Separator />

      <TaskCallbackDeliveries />
    </SettingsSection>
  )
}
