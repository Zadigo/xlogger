import { h } from 'vue'
import type { TableColumn } from '@nuxt/ui'
import type { LogFile, LogFileContent } from '~/types'
import { NuxtLink, UCheckbox, UBadge, UButton } from '#components'
import type { ButtonProps } from '@nuxt/ui'

export function useFilesTable() {
  const columns: TableColumn<LogFile>[] = [
    {
      id: 'select',
      header: ({ table }) =>
        h(UCheckbox, {
          modelValue: table.getIsSomePageRowsSelected() ? 'indeterminate' : table.getIsAllPageRowsSelected(),
          'onUpdate:modelValue': (value: boolean | 'indeterminate') => table.toggleAllPageRowsSelected(!!value),
          'aria-label': 'Select all'
        }),
      cell: ({ row }) =>
        h(UCheckbox, {
          modelValue: row.getIsSelected(),
          'onUpdate:modelValue': (value: boolean | 'indeterminate') => row.toggleSelected(!!value),
          'aria-label': 'Select row'
        })
    },
    {
      accessorKey: 'id',
      header: 'ID',
      cell: ({ row }) => {
        return row.getValue('id')
      }
    },
    {
      accessorKey: 'name',
      header: 'Name',
      cell: ({ row }) => {
        const encodedName = base64Name(row.getValue('name'))
        return h(NuxtLink, { to: `/logfiles/${encodedName}` }, () => row.getValue('name'))
      }
    },
    {
      accessorKey: 'path',
      header: 'Path',
      cell: ({ row }) => {
        return row.getValue('path')
      }
    },
  ]

  return {
    columns,
    base64Name
  }
}

export function useLogsTable<T extends LogFileContent = LogFileContent>(data: MaybeRefOrGetter<T[] | undefined>) {
  const columns: TableColumn<T>[] = [
    {
      id: 'expand',
      cell: ({ row }) =>
        h(UButton, {
          color: 'neutral',
          variant: 'ghost',
          icon: 'i-lucide-chevron-down',
          square: true,
          'aria-label': 'Expand',
          ui: {
            leadingIcon: [
              'transition-transform',
              row.getIsExpanded() ? 'duration-200 rotate-180' : ''
            ]
          },
          onClick: () => row.toggleExpanded()
        })
    },
    {
      accessorKey: 'raw_line',
      header: 'Raw Line',
      cell: ({ row }) => {
        return row.getValue('raw_line')
      }
    },
    {
      accessorKey: 'remote_address',
      header: 'IP Address',
      cell: ({ row }) => {
        return row.getValue('remote_address')
      }
    },
    {
      accessorKey: 'method',
      header: 'Method',
      cell: ({ row }) => {
        const color = {
          GET: 'success' as const,
          POST: 'warning' as const,
          PUT: 'warning' as const,
          PATCH: 'warning' as const,
        }[row.getValue('method') as string]

        return h(UBadge, { class: 'capitalize', variant: 'subtle', color }, () =>
          row.getValue('method')
        )
      }
    },
    {
      accessorKey: 'date_time',
      header: 'Date Time',
      cell: ({ row }) => {
        return row.getValue('date_time')
      }
    },
    {
      accessorKey: 'path',
      header: 'Path',
      cell: ({ row }) => {
        return row.getValue('path')
      }
    },
    {
      accessorKey: 'protocole',
      header: 'Protocole',
      cell: ({ row }) => {
        return row.getValue('protocole')
      }
    },
    {
      accessorKey: 'is_success',
      header: 'Success',
      cell: ({ row }) => {
        const color: ButtonProps['color'] = row.getValue('is_success') === true ? 'success' : 'error'
        return h(UBadge, { class: 'capitalize', variant: 'subtle', color }, () =>
          row.getValue('is_success') ? 'Success' : 'Failure'
        )
      }
    },
  ]

  const columnVisibility = ref<Partial<Record<keyof T, boolean>>>({
    raw_line: false,
    user_agent: false,
    body_bytes_sent: false,
    referrer: false,
    meta_data: false,
  })

  return {
    columns,
    columnVisibility
  }
}
