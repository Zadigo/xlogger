import { h } from 'vue'
import type { TableColumn } from '@nuxt/ui'
import type { LogFile } from '~/types'
import { NuxtLink, UCheckbox } from '#components'

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

export function useLogsTable() {

}
