type _QueryParams = {
  limit: string
  offset: string
  successful: string
  status: string
  methods: string
  startDate: string
  endDate: string
  search: string
}

export type QueryParams = Partial<_QueryParams>
