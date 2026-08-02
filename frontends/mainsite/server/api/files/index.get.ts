import type { LogFile } from  '~~/app/types'

export default defineEventHandler(async (_event) => {
  const logs: LogFile[] = [
    {
      id: '1',
      title: 'Application Log',
      path: '/logs/application.log'
    }
  ]
  return logs
})
