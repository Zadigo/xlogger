import type { LogFileContent } from '~~/app/types'
import { createErrorTemplate } from '~~/app/utils/error'

export default defineEventHandler(async (event) => {
  try {
    const { id } = getRouterParams(event) as { id: string }
    const query = getQuery<{ limit: string, offset: string }>(event)

    return await $fetch<LogFileContent[]>(`/v1/files/${id}`, {
      method: 'GET',
      baseURL: 'http://127.0.0.1:9000',
      query: {
        limit: query.limit || '100',
        offset: query.offset || '0'
      },
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
