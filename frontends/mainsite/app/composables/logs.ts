import { h, resolveComponent } from 'vue'
import type { TableColumn } from '@nuxt/ui'
import type { LogFile } from '~/types'
import { NuxtLink, UCheckbox, UBadge } from '#components'

export function useLogsTable() {
  // const UCheckbox = resolveComponent('UCheckbox')
  // const UBadge = resolveComponent('UBadge')
  // const NuxtLink = resolveComponent('NuxtLink')
  
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
      accessorKey: 'title',
      header: 'Title',
      cell: ({ row }) => {
        return h(NuxtLink, { to: `/logfiles/${row.getValue('id')}` }, () => row.getValue('title'))
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
    columns
  }
}
