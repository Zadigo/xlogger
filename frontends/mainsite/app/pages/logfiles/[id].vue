<template>
  <section id="admin">
    <u-card>
      <template #title>
        <h1 class="text-lg font-semibold">Log Files</h1>
      </template>

      <u-input :model-value="table?.tableApi?.getColumn('title')?.getFilterValue() as string" class="max-w-sm" placeholder="Filter files..." @update:model-value="table?.tableApi?.getColumn('title')?.setFilterValue($event)" />

      <u-table ref="table" v-model:column-filters="columnFilters" :data="fileContents" :loading="!hasLogs" />
    </u-card>
  </section>
</template>

<script setup lang="ts">
import type { LogFileContent } from '~/types'

definePageMeta({
  layout: 'admin'
})

const encodedName = useRoute().params.id as string

const fileContents = computedAsync(async () => {
  return await $fetch<LogFileContent[]>(`/api/files/${encodedName}`, {
    method: 'GET'
  })
})

const hasLogs = computed(() => isDefined(fileContents) && fileContents.value.length > 0)

// const { columns } = useLogsTable()

const table = useTemplateRef('table')

const columnFilters = ref([])
</script>
