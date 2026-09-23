import type { LogFileContent } from '~/types'

export function useFilteringByFileType(files: Ref<LogFileContent[]>) {
  return {
    php: computed(() => files.value.filter(file => file.metaData.isPhp)),
    js: computed(() => files.value.filter(file => file.metaData.isJs)),
    env: computed(() => files.value.filter(file => file.metaData.isEnv)),
    github: computed(() => files.value.filter(file => file.metaData.isGitHub)),
    nuxt: computed(() => files.value.filter(file => file.metaData.isNuxt)),
    powershell: computed(() => files.value.filter(file => file.metaData.isPowershell))
  }
}
