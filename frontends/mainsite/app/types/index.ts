import zod from 'zod'


const LogSchema = zod.object({
  id: zod.string(),
  title: zod.string(),
  path: zod.string()
})

export type LogFile = zod.infer<typeof LogSchema>

const LogFileContetntSchema = zod.object({
  userAgent: zod.string(),
})


export type LogFileContent = zod.infer<typeof LogFileContetntSchema>
