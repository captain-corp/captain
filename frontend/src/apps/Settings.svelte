<script lang="ts">
  import { Input, Label, Select, Toggle } from 'flowbite-svelte';

  import { type Settings } from '../utils/types/settings';
  import SubmitButton from '../lib/SubmitButton.svelte';
  import type { SavingStates } from '../utils/types/common';
  import FileDropzone from '../lib/FileDropzone.svelte';

  let {
    title = '',
    subtitle = '',
    logoID = 0,
    useLogoAsFavicon = false,
    codeHighlightTheme = 'default',
    codeThemeOptions = [],
    postsPerPage = 10,
    onSubmit = (data: any, done: (savingState: SavingStates) => void, error: (error: any) => void) => {},
  }: Settings = $props();

  let error = $state('');
  let savingState = $state('draft');

  function handleSubmit(e: Event) {
    e.preventDefault();
    error = '';
    savingState = 'saving';

    onSubmit(
      {
        title,
        subtitle,
        logoID,
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
    <FileDropzone accept="image/*" onFiles={(files) => {}} />
  </div>

  <!-- Use logo as favicon -->
  <div>
    <Label for="use-logo-as-favicon" class="block text-sm font-bold text-gray-700 mb-2">Use logo as favicon</Label>
    <Toggle id="use-logo-as-favicon" name="use-logo-as-favicon" bind:checked={useLogoAsFavicon}>
      <svelte:fragment slot="offLabel">No</svelte:fragment>
      <span>Yes</span>
    </Toggle>
  </div>

  <!-- Code highlighting theme -->
  <div>
    <Label for="code-theme" class="block text-sm font-bold text-gray-700 mb-2">Code highlighting theme</Label>
    <Select id="code-theme" name="code-theme" bind:value={codeHighlightTheme} items={codeThemeOptions}></Select>
  </div>

  <!-- Posts per page -->
  <div>
    <Label for="posts-per-page" class="block text-sm font-bold text-gray-700 mb-2">Posts per page</Label>
    <Input type="number" id="posts-per-page" name="posts-per-page" bind:value={postsPerPage} required />
  </div>

  <!-- Save Button -->
  <div class="flex justify-end">
    <SubmitButton {savingState} />
  </div>
</form>

