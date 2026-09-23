import type { QueryParams } from '#shared/types/query'

export const useFilteringComposable = createSharedComposable(() => {
  const query = useUrlSearchParams() as QueryParams

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

  const search = ref<string>('')

  watch(search, (newValue) => {
    query.search = newValue
  })

  const onlySuccessfulRequests = ref(false)
  watch(onlySuccessfulRequests, (newValue) => {
    query.successful = newValue ? '1' : undefined
  })
  
  const statusCodes = ref<string[]>([])
  watch(statusCodes, (newValue) => {
    query.status = newValue.join(',')
  })

  const httpMethods = ref<string[]>([])
  watch(httpMethods, (newValue) => {
    query.methods = newValue.join(',')
  })

  const startDate = ref<string | undefined>(undefined)
  const endDate = ref<string | undefined>(undefined)
  watch([startDate, endDate], ([newStartDate, newEndDate]) => {
    query.startDate = newStartDate
    query.endDate = newEndDate
  })

  const queryDict = computed(() => ({ ...query }))

  return {
    limitOffset,
    search,
    onlySuccessfulRequests,
    statusCodes,
    httpMethods,
    startDate,
    endDate,
    queryDict
  }
})
