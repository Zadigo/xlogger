<template>
  <section id="admin">
    <u-card>
      <template #title>
        <h1 class="text-lg font-semibold">Log Files</h1>
      </template>

      <div>
        <u-input :model-value="table?.tableApi?.getColumn('title')?.getFilterValue() as string" class="max-w-sm" placeholder="Filter files..." @update:model-value="table?.tableApi?.getColumn('title')?.setFilterValue($event)" />
        <u-button color="info" variant="soft">
          <icon name="i-lucide-refresh-ccw" />
        </u-button>

        <u-dropdown-menu :items="dropdown()" :content="{ align: 'end' }">
          <u-button label="Columns" color="neutral" variant="outline" trailing-icon="i-lucide-chevron-down" />
        </u-dropdown-menu>
      </div>

      <u-table 
        ref="table" 
        v-model:column-filters="columnFilters" 
        v-model:column-visibility="columnVisibility"
        :data="fileContents" 
        :loading="!hasLogs"
        :columns="columns"
      >
        <template #expanded="{ row }">
          <u-card>
            <pre>{{ row.getValue<LogFileContent['metaData']>('meta_data') }}</pre>
          </u-card>
        </template>
      </u-table>
    </u-card>
  </section>
</template>

<script setup lang="ts">
import type { LogFileContent } from '~/types'

definePageMeta({
  layout: 'admin'
})

const table = useTemplateRef('table')

/**
 * Data
 */

const encodedName = useRoute().params.id as string

const limitOffset = computed({
  set: (value) => {
    const params = useUrlSearchParams() as { limit?: string; offset?: string }
    params.limit = value.limit
    params.offset = value.offset
  },
  get: () => {
    const params = useUrlSearchParams() as { limit?: string; offset?: string }
    return {
      limit: params.limit ?? '100',
      offset: params.offset ?? '0'
    }
  }
})

const fileContents = computedAsync(async () => {
  return await $fetch<LogFileContent[]>(`/api/files/${encodedName}`, {
    method: 'GET',
    query: {
      limit: limitOffset.value.limit,
      offset: limitOffset.value.offset
    }
  })
})

const hasLogs = computed(() => isDefined(fileContents) && fileContents.value.length > 0)

// onMounted(() => {
//   useInfiniteScroll(
//     table.value?.$el,
//     () => {
//       // if (status.value === 'pending') return
//       if (!hasLogs.value) return

//       const offset = parseInt(limitOffset.value.offset)
//       const limit = parseInt(limitOffset.value.limit)

//       limitOffset.value = {
//         limit: limit.toString(),
//         offset: (offset + limit).toString()
//       }
//     },
//     {
//       distance: 200,
//       canLoadMore: () => {
//         return true
//       }
//     }
//   )
// })

/**
 * Table
 */

const { columns, columnVisibility } = useLogsTable(fileContents)

/**
 * Other
 */

const columnFilters = ref([])

/**
 * Dropdown menu
 */

const dropdown = () =>{
  return table.value?.tableApi?.getAllColumns().filter((column) => column.getCanHide())
  .map((column) => ({
    label: column.id,
    type: 'checkbox' as const,
    checked: column.getIsVisible(),
    onUpdateChecked(checked: boolean) {
      table.value?.tableApi?.getColumn(column.id)?.toggleVisibility(!!checked)
    },
    onSelect(e: Event) {
      e.preventDefault()
    }
  }))
}
</script>
