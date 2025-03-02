<script lang="ts">
  import { onMount } from 'svelte';
  import { Alert, Input, Label, Toggle, Select } from 'flowbite-svelte';

  import { type Pages } from '../utils/types/pages';
  import { slugify } from '../utils/text';
  import SubmitButton from '../lib/SubmitButton.svelte';
  import Editor from '../lib/Editor.svelte';

  let protectedSlug = $state(false);
  let protectedSlugViolation = $state(false);
  let error = $state('');

  let {
    title = '',
    content = '',
    visible = false,
    slug = '',
    savingState = 'draft',
    contentType = '',
    onSubmit = (data: any) => {},
  }: Pages = $props();

  onMount(() => {
    if (slug !== '') {
      protectedSlug = true;
    }
  });

  const contentTypeOptions = [
    { value: 'html', name: 'HTML' },
    { value: 'markdown', name: 'Markdown' },
  ];

  function updateSlug() {
    if (protectedSlug) {
      return;
    }
    slug = slugify(title);
  }

  function onSlugChange(e: Event) {
    const target = e.target as HTMLInputElement;
    protectedSlugViolation = target.value !== slug && protectedSlug;
  }

  function handleSubmit(e: Event) {
    e.preventDefault();
    error = '';

    onSubmit(
      {
        contentType,
        title,
        content,
        visible,
        slug,
      },
      (savingState) => {
        savingState = savingState;
      },
      (newError) => {
        error = newError;
      }
    );
  }
</script>

<div>
  {#if error !== ''}
    <Alert color="red" class="my-2">
      <span class="font-medium">{error}</span>
    </Alert>
  {/if}

  <form onsubmit={handleSubmit} class="space-y-6">
    <!-- Title -->
    <div>
      <Label for="title" class="block text-sm font-bold text-gray-700 mb-2">Title</Label>
      <Input type="text" id="title" name="title" bind:value={title} onkeyup={updateSlug} required />
    </div>

    <!-- Slug -->
    <div>
      <Label for="slug" class="block text-sm font-bold text-gray-700 mb-2">Slug</Label>
      <Input type="text" id="slug" name="slug" value={slug} onkeyup={onSlugChange} required />
      {#if protectedSlugViolation}
        <Alert color="red" class="mt-2">
          <span class="font-medium">Changing the slug may break links.</span>
        </Alert>
      {/if}
    </div>

    <!-- Content Type -->
    <div>
      <Label for="content-type" class="block text-sm font-bold text-gray-700 mb-2"
        >Content Type</Label
      >
      <Select id="content-type" bind:value={contentType} items={contentTypeOptions}></Select>
    </div>

    <!-- Content -->
    <div>
      <Label for="content" class="block text-sm font-bold text-gray-700 mb-2">Content</Label>
      <Editor name="content" bind:value={content} rows={10} />
    </div>

    <div class="grid gap-4 sm:grid-cols-2 sm:gap-6">
      <!-- Visibility Toggle -->
      <div>
        <Label for="visible" class="block text-sm font-bold text-gray-700 mb-4">Visibility</Label>
        <Toggle id="visible" name="visible" bind:checked={visible}>
          <svelte:fragment slot="offLabel">Hidden</svelte:fragment>
          <span>Visible</span>
        </Toggle>
      </div>
    </div>

    <!-- Submit Button -->
    <div class="flex justify-end">
      <SubmitButton {savingState} />
    </div>
  </form>
</div>
