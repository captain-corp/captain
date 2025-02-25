<script lang="ts">
  import { Spinner } from 'flowbite-svelte';
  import { type FileDropzone } from '../utils/types/filedropzone';

  let {
    accept = '*/*',
    multiple = false,
    maxSize = 10_000_000,
    onFilesSelectedHandler = () => {},
  }: FileDropzone = $props();
  let isDragging = $state(false);
  let isUploading = $state(false);
  let error = $state<string | null>(null);
  let files = $state<File[]>([]);

  $effect(() => {
    if (files.length > 0) {
      const validFiles = files.filter((file) => {
        if (file.size > maxSize) {
          error = `File ${file.name} exceeds ${maxSize / 1_000_000}MB limit`;
          return false;
        }
        return true;
      });

      if (validFiles.length > 0) {
        onFilesSelectedHandler(validFiles, uploadStarted, uploadFinished);
      }
    }
  });

  function uploadStarted(): void {
    isUploading = true;
    error = null;
    files = [];
  }

  function uploadFinished(): void {
    isUploading = false;
    files = [];
  }

  function preventDefaults(e: Event): void {
    e.preventDefault();
    e.stopPropagation();
  }

  function handleDragEnter(): void {
    isDragging = true;
  }

  function handleDragLeave(): void {
    isDragging = false;
  }

  function handleDrop(e: DragEvent): void {
    preventDefaults(e);
    isDragging = false;
    files = Array.from(e.dataTransfer?.files ?? []);
  }

  function handleInputChange(e: Event): void {
    files = Array.from((e.target as HTMLInputElement).files ?? []);
    (e.target as HTMLInputElement).value = '';
  }
</script>

<svelte:window on:dragover={preventDefaults} on:drop={preventDefaults} />

<!-- svelte-ignore a11y_no_static_element_interactions -->
<div
  class="group relative cursor-pointer border-2 border-dashed rounded-xl p-8 text-center transition-colors"
  class:border-indigo-500={isDragging}
  class:border-red-500={error}
  ondragover={handleDragEnter}
  ondragleave={handleDragLeave}
  ondrop={handleDrop}
>
  {#if error}
    <p class="text-sm text-destructive mt-2">{error}</p>
  {/if}

  {#if isUploading}
    <div class="absolute inset-0 w-full h-full flex items-center justify-center">
      <Spinner class="size-8" />
    </div>
  {:else}
    <input
      type="file"
      {accept}
      {multiple}
      onchange={handleInputChange}
      class="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
      aria-label="File upload"
    />

    <div class="flex flex-col items-center justify-center gap-3">
      <div
        class="h-12 w-12 bg-muted rounded-full flex items-center justify-center transition-colors group-hover:bg-muted/80"
      >
        <svg
          class="w-6 h-6 text-muted-foreground"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
          />
        </svg>
      </div>

      <p class="text-sm text-muted-foreground font-medium">Drag files here or click to upload</p>
    </div>
  {/if}
</div>
