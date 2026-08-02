import type { LogFileContent } from '~~/app/types'

export default defineEventHandler(async (_event) => {
  const logs: LogFileContent[] = [
    {
      userAgent: 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36'
    }
  ]
  return logs
})
