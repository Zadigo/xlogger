<template>
  <section id="admin">
    <u-card>
      <template #title>
        <h1 class="text-lg font-semibold">Log Files</h1>
      </template>

      <u-input :model-value="table?.tableApi?.getColumn('title')?.getFilterValue() as string" class="max-w-sm" placeholder="Filter files..." @update:model-value="table?.tableApi?.getColumn('title')?.setFilterValue($event)" />

      <u-table ref="table" v-model:column-filters="columnFilters" :data="logs" :loading="!hasLogs" :columns="columns" />
    </u-card>
  </section>
</template>

<script setup lang="ts">
import type { LogFile } from '~/types'

definePageMeta({
  layout: 'admin'
})

const logs = computedAsync(async () => {
  return await $fetch<LogFile[]>('/api/files', {
    method: 'GET'
  })
})

const hasLogs = computed(() => isDefined(logs) && logs.value.length > 0)

const { columns } = useLogsTable()

const table = useTemplateRef('table')

const columnFilters = ref([])
</script>
