import type { LogFile } from  '~~/app/types'
import { createErrorTemplate } from '~~/app/utils/error'

export default defineEventHandler(async (_event) => {
  // const logs: LogFile[] = [
  //   {
  //     id: '1',
  //     title: 'Application Log',
  //     path: '/logs/application.log'
  //   }
  // ]

  try {
    return await $fetch<LogFile[]>('/v1/files', {
      method: 'GET',
      baseURL: 'http://127.0.0.1:9000',
      headers: {
        'Origin': 'http://localhost:3000',
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      }
    })
  } catch(error) {
    const template = createErrorTemplate(error)
    return createError(template)
  }
})
