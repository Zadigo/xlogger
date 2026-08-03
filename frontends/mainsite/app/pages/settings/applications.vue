<template>
  <section id="account" class="max-w-6xl gap-x-2 mx-auto">
    <u-card>      
      <p class="font-light text-center">
        Applications are third-party services that can access your account. You can manage and revoke access to these applications here.
      </p>

      <div class="flex items-center justify-center w-full p-5">
        <u-button @click="() => { toggleModal() }">
          <icon name="i-lucide-plus" />
          Create New Application
        </u-button>
      </div>
    </u-card>

    <u-modal v-model:open="openModal">
      <template #body>
        <div class="flex flex-col gap-2">
          <u-input v-model="crendentials.name" placeholder="Name" />
          <u-input v-model="crendentials.clientId" placeholder="Client ID" />
          <u-input v-model="crendentials.clientSecret" placeholder="Client Secret" />
          <u-input-menu v-model="crendentials.redirectUris" placeholder="Redirect URI" />
        </div>
      </template>

      <template #footer>
        <div class="flex justify-end gap-2">
          <u-button @click="() => { toggleModal() }">
            Cancel
          </u-button>
          
          <u-button>
            Create
          </u-button>
        </div>
      </template>
    </u-modal>
  </section>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'admin',
  title: 'Account Settings'
})

const { openModal, crendentials, toggleModal } = useApplications()
</script>
