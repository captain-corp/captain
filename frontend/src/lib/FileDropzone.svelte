<script lang="ts">
  let { accept = '*/*', multiple = false, maxSize = 10_000_000, onFiles = () => {} } = $props();
  let isDragging = $state(false);
  let error = $state<string | null>(null);
  let files = $state<File[]>([]);

  $effect(() => {
    if (files.length > 0) {
      const validFiles = files.filter(file => {
        if (file.size > maxSize) {
          error = `File ${file.name} exceeds ${maxSize / 1_000_000}MB limit`;
          return false;
        }
        return true;
      });

      if (validFiles.length > 0) {
        onFiles(validFiles);
      }
    }
  });

  function preventDefaults(e: Event) {
    e.preventDefault();
    e.stopPropagation();
  }

  function handleDragEnter() {
    isDragging = true;
    console.log('drag enter');
  }

  function handleDragLeave() {
    isDragging = false;
    console.log('drag leave');
  }

  function handleDrop(e: DragEvent) {
    isDragging = false;
    files = Array.from(e.dataTransfer?.files ?? []);
  }

  function handleInputChange(e: Event) {
    files = Array.from((e.target as HTMLInputElement).files ?? []);
    (e.target as HTMLInputElement).value = '';
  }
</script>

<svelte:window
  on:dragover={preventDefaults}
  on:drop={preventDefaults}
/>

<div
  class="group relative cursor-pointer border-2 border-dashed rounded-xl p-8 text-center transition-colors"
  class:border-primary={isDragging}
  class:border-destructive={error}

  on:dragover={handleDragEnter}
  on:dragleave={handleDragLeave}
  on:drop|preventDefault={handleDrop}
>
  <input
    type="file"
    {accept}
    {multiple}
    on:change={handleInputChange}
    class="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
    aria-label="File upload"
  />

  <div class="flex flex-col items-center justify-center gap-3">
    <div class="h-12 w-12 bg-muted rounded-full flex items-center justify-center transition-colors group-hover:bg-muted/80">
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

    <p class="text-sm text-muted-foreground font-medium">
      Drag files here or click to upload
    </p>

    {#if error}
      <p class="text-sm text-destructive mt-2">{error}</p>
    {/if}
  </div>
</div>
