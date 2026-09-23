<template>
  <div>
    <u-input v-model="search" class="max-w-sm" placeholder="Search by name, user agent..." />

    <div class="flex gap-2">
      <u-button v-for="type in Object.keys(byTypes)" :key="type" :label="type" variant="subtle" color="info" class="w-25 flex justify-center" />
    </div>

    <u-switch label="Only successful requests" v-model="onlySuccessfulRequests" />
    <u-select v-model="httpMethods" :items="['GET', 'POST', 'PATCH', 'DELETE']" class="w-30" multiple />

    <u-input placeholder="IP Regex" class="w-50" />
    <u-input placeholder="User Agent Regex" class="w-50" />

    <u-select :items="['Any of', 'Contains', 'Is', 'Is not', 'None of']" placeholder="Filter type" class="w-30" />
    <u-select v-model="statusCodes" :items="['200', '201', '202', '301', '404', '500']" placeholder="Status codes" class="w-30" multiple />
    <u-input placeholder="Request Path Regex" class="w-50" />

    <div>
      <p class="font-semibold">Date Range</p>
      <u-input-date v-model="startDate" class="w-30" />
      <u-input-date v-model="endDate" class="w-30" />
    </div>
  </div>
</template>

<script setup lang="ts">
import type { LogFileContent } from '~/types'
import { useFilteringByFileType } from '~/utils/filtering'

const { files = [] } = defineProps<{ files: LogFileContent[] }>()

const _files = computed(() => files)
const byTypes = useFilteringByFileType(_files)

/**
 * Filtering composable
 */

const { onlySuccessfulRequests, statusCodes, httpMethods, startDate, endDate, search } = useFilteringComposable()
</script>
