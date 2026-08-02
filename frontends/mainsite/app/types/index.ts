import zod from 'zod'


const LogSchema = zod.object({
  id: zod.string(),
  name: zod.string(),
  path: zod.string()
})

export type LogFile = zod.infer<typeof LogSchema>

const metaDataSchema = zod.object({
  isPhp: zod.boolean(),
  isAssets: zod.boolean(),
  isJs: zod.boolean(),
  isHttp2: zod.boolean(),
  isRobotsTxt: zod.boolean(),
  isXml: zod.boolean(),
  isAttemptedLogin: zod.boolean(),
  isWordPress: zod.boolean(),
  isEnv: zod.boolean(),
  isExecutable: zod.boolean(),
  isPowerShell: zod.boolean(),
  isNuxt: zod.boolean(),
  isGponRouter: zod.boolean(),
  isWindowsPath: zod.boolean(),
  isGitHub: zod.boolean(),
})

const LogFileContetntSchema = zod.object({
  rawLine: zod.string(),
  remoteAddress: zod.string(),
  dateTime: zod.string(),
  method: zod.string(),
  path: zod.string(),
  protocole: zod.string(),
  statusCode: zod.number(),
  bodyBytesSent: zod.number(),
  referrer: zod.string(),
  userAgent: zod.string(),
  requestTime: zod.number(),
  isSuccess: zod.boolean(),
  metaData: metaDataSchema
})


export type LogFileContent = zod.infer<typeof LogFileContetntSchema>
