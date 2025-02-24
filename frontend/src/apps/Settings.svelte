<script lang="ts">
  import { Input, Label, Select, Toast, Toggle } from 'flowbite-svelte';

  import { type Settings } from '../utils/types/settings';
  import SubmitButton from '../lib/SubmitButton.svelte';
  import type { SavingStates } from '../utils/types/common';
  import FileDropzone from '../lib/FileDropzone.svelte';
  import Logo from '../lib/Logo.svelte';
  import { CheckCircleSolid, CloseCircleSolid } from 'flowbite-svelte-icons';

  let {
    title = '',
    subtitle = '',
    logoUrl = '',
    useLogoAsFavicon = false,
    codeHighlightTheme = 'default',
    codeThemeOptions = [],
    postsPerPage = 10,
    onSubmit = (
      data: any,
      done: (savingState: SavingStates) => void,
      error: (error: any) => void
    ) => {},
    uploadLogoHandler = (
      files: File[],
      uploadStarted: () => void,
      uploadFinished: (error: string | null, url: string | null) => void
    ) => {},
    deleteLogoHandler = () => {},
  }: Settings = $props();

  let error = $state('');
  let message = $state('');
  let savingState = $state('draft');

  function setLogo(url: string) {
    logoUrl = url;
    message = 'Upload done';
  }

  function deleteLogo() {
    logoUrl = '';
    message = 'Logo deleted';
  }

  function handleSubmit(e: Event) {
    e.preventDefault();
    error = '';
    savingState = 'saving';

    onSubmit(
      {
        title,
        subtitle,
        logoUrl,
        useLogoAsFavicon,
        codeHighlightTheme,
        codeThemeOptions,
        postsPerPage,
      },
      (newSavingState: string) => {
        savingState = newSavingState;
      },
      (newError: string) => {
        error = newError;
      }
    );
  }
</script>

<form onsubmit={handleSubmit} class="space-y-6">
  <!-- Site Title -->
  <div>
    <Label for="title" class="block text-sm font-bold text-gray-700 mb-2">Site Title</Label>
    <Input type="text" id="title" name="title" bind:value={title} required />
  </div>

  <!-- Site Subtitle -->
  <div>
    <Label for="subtitle" class="block text-sm font-bold text-gray-700 mb-2">Site Subtitle</Label>
    <Input type="text" id="subtitle" name="subtitle" bind:value={subtitle} required />
  </div>

  <!-- Logo -->
  <div>
    <Label for="logo" class="block text-sm font-bold text-gray-700 mb-2">Logo</Label>
    {#if logoUrl}
      <Logo
        {logoUrl}
        onDelete={() => {
          deleteLogoHandler();
          deleteLogo();
        }}
      />
    {:else}
      <FileDropzone
        accept="image/*"
        onFilesSelectedHandler={(files, uploadStarted, uploadFinished) => {
          uploadLogoHandler(files, uploadStarted, (error, url) => {
            console.log(error);
            console.log(url);
            if (error) {
              error = error;
            } else {
              setLogo(url);
            }
            uploadFinished(error, url);
          });
        }}
      />
    {/if}

    {#if error}
      <Toast color="red" dismissable={true} position="top-right">
        <svelte:fragment slot="icon">
          <CloseCircleSolid class="w-5 h-5" />
          <span class="sr-only">Error icon</span>
        </svelte:fragment>
        {error}
      </Toast>
    {/if}

    {#if message}
      <Toast color="green" dismissable={true} position="top-right">
        <svelte:fragment slot="icon">
          <CheckCircleSolid class="w-5 h-5" />
          <span class="sr-only">Check icon</span>
        </svelte:fragment>
        {message}
      </Toast>
    {/if}
  </div>

  <!-- Use logo as favicon -->
  <div>
    <Label for="use-logo-as-favicon" class="block text-sm font-bold text-gray-700 mb-2"
      >Use logo as favicon</Label
    >
    <Toggle id="use-logo-as-favicon" name="use-logo-as-favicon" bind:checked={useLogoAsFavicon}>
      <svelte:fragment slot="offLabel">No</svelte:fragment>
      <span>Yes</span>
    </Toggle>
  </div>

  <!-- Code highlighting theme -->
  <div>
    <Label for="code-theme" class="block text-sm font-bold text-gray-700 mb-2"
      >Code highlighting theme</Label
    >
    <Select
      id="code-theme"
      name="code-theme"
      bind:value={codeHighlightTheme}
      items={codeThemeOptions}
    ></Select>
  </div>

  <!-- Posts per page -->
  <div>
    <Label for="posts-per-page" class="block text-sm font-bold text-gray-700 mb-2"
      >Posts per page</Label
    >
    <Input
      type="number"
      id="posts-per-page"
      name="posts-per-page"
      bind:value={postsPerPage}
      required
    />
  </div>

  <!-- Save Button -->
  <div class="flex justify-end">
    <SubmitButton {savingState} />
  </div>
</form>
