import type { LogFileContent } from '~~/app/types'

export default defineEventHandler(async (_event) => {
  const logs: LogFileContent[] = [
    {
      rawLine: `172.21.0.2 - - [08/May/2025:13:53:53 +0000] "GET /assets/GrowthPage-CvJGw4NJ.js HTTP/1.1" 200 9385 "-" "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36"`,
      remoteAddress: '172.21.0.2',
      dateTime: '08/May/2025:13:53:53 +0000',
      method: 'GET',
      path: '/assets/GrowthPage-CvJGw4NJ.js',
      protocole: 'HTTP/1.1',
      statusCode: 200,
      bodyBytesSent: 9385,
      referrer: '-',
      userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36',
      requestTime: 0.000,
      isSuccess: true,
      metaData: {
        isPhp: false,
        isAssets: true,
        isJs: true,
        isHttp2: false,
        isRobotsTxt: false,
        isXml: false,
        isAttemptedLogin: false,
        isWordPress: false,
        isEnv: false,
        isExecutable: false,
        isPowerShell: false,
        isNuxt: false,
        isGponRouter: false,
        isWindowsPath: false,
        isGitHub: false
      }
    }
  ]
  return logs
})
