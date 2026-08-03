export type Credentials = {
  name: string
  clientId: string
  clientSecret: string
  redirectUris: string[]
  scopes: string[]
}

export function useApplications() {
  const openModal = ref(false)
  const toggleModal = useToggle(openModal)

  const crendentials = ref<Credentials>({
    name: '',
    clientId: '',
    clientSecret: '',
    redirectUris: [],
    scopes: [],
  })

  return {
    crendentials,
    openModal,
    toggleModal,
  }
}
