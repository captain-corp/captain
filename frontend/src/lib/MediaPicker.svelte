<script lang="ts">
  import { Button, P, Spinner, Hr } from 'flowbite-svelte';
  import FileDropzone from './FileDropzone.svelte';

  // Define the media item interface to match Gallery component expectations
  interface MediaItem {
    alt: string;
    src: string;
    id?: number; // Optional ID for database reference
  }

  const {
    mediaProvider,
    onMediaSelect = (media: MediaItem) => {
      console.log('default', media);
    }, // Callback when media is selected
    uploadHandler = (
      files: File[],
      uploadStarted: () => void,
      uploadFinished: (error: string | null, media: MediaItem | null) => void
    ) => {},
  } = $props();

  // Convert callback-based mediaProvider to Promise-based with proper typing
  const mediaPromise = new Promise<MediaItem[]>((resolve, reject) => {
    mediaProvider(resolve, reject);
  });

  // Handle media selection
  function handleMediaSelect(media: MediaItem) {
    onMediaSelect(media);
  }

  // Handle file upload
  function handleFileUpload(files: File[]) {
    if (!files.length) return;

    uploadHandler(
      files,
      () => {
        // Upload started - could show a loading state here
      },
      (error, media) => {
        if (error) {
          console.error(error);
        } else if (media) {
          // Automatically select the newly uploaded media
          onMediaSelect(media);
        }
      }
    );
  }
</script>

<P class="mb-4 text-center">
  Select an image from your media library
</P>

<div class="my-6">
  <FileDropzone
    accept="image/*"
    onFilesSelectedHandler={(files, uploadStarted, uploadFinished) => {
      uploadHandler(files, uploadStarted, (error, media) => {
        uploadFinished(error, media ? media.src : null);
      });
    }}
  />
</div>

<Hr />

{#await mediaPromise}
  <div class="flex justify-center items-center py-10">
    <Spinner size="12" color="primary" />
  </div>
{:then media}
  <div class="grid gap-4 grid-cols-2 md:grid-cols-3">
    {#each media as item}
      <div
        class="relative overflow-hidden rounded-lg cursor-pointer hover:opacity-90 transition-opacity"
        onclick={() => handleMediaSelect(item)}
        onkeydown={(e) => e.key === 'Enter' && handleMediaSelect(item)}
        tabindex="0"
        role="button"
        aria-label="Select {item.alt}"
      >
        <img src={item.src} alt={item.alt} class="w-full h-auto object-cover" />
      </div>
    {/each}
  </div>
{:catch error}
  <P class="text-center text-red-500">Error loading media: {error}</P>
{/await}

<Hr />

<P class="mt-4 text-center">
  <Button
    onclick={() => handleFileUpload([new File([], 'new-image.jpg')])}
    class="bg-indigo-600 py-2 px-4 text-sm font-bold text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
  >
    Upload Image
  </Button>
</P>
