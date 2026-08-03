<template>
  <section id="account" class="max-w-6xl gap-x-2 mx-auto">
    <div class="grid grid-cols-2 gap-x-2">
      <u-card>
        <template #title>
          Password
        </template>
        
        lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
      </u-card>
      
      <u-card>
        <div class="flex flex-col gap-4">
          <u-input type="password" label="Current Password" placeholder="Current Password" />
          <u-input type="password" label="New Password" placeholder="New Password" />
          <u-input type="password" label="Confirm New Password" placeholder="Confirm New Password" />
          <div>
            <u-button color="primary">
              Change Password
            </u-button>
          </div>
        </div>
      </u-card>
    </div>
  </section>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'admin',
  title: 'Account Settings'
})
</script>
