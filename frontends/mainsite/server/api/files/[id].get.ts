import type { LogFileContent } from '~~/app/types'
import type { QueryParams } from '#shared/types/query'
import { createErrorTemplate } from '~~/app/utils/error'

export default defineEventHandler(async (event) => {
  try {
    const { id } = getRouterParams(event) as { id: string }
    const query = getQuery<QueryParams>(event)

    return await $fetch<LogFileContent[]>(`/v1/files/${id}`, {
      method: 'GET',
      baseURL: 'http://127.0.0.1:9000',
      query,
      headers: {
        'Origin': 'http://localhost:3000',
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      }
    })
  } catch (error) {
    const template = createErrorTemplate(error)
    return createError(template)
  }
})
