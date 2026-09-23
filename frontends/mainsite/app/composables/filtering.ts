export const useFilteringComposable = createSharedComposable(() => {
  const query = useUrlSearchParams() as {
    successful?: string
    status?: string
    methods?: string
    startDate?: string
    endDate?: string
    search?: string
  }

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

  return {
    search,
    onlySuccessfulRequests,
    statusCodes,
    httpMethods,
    startDate,
    endDate
  }
})
