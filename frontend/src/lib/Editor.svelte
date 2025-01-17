<script lang="ts">
  import { Textarea } from 'flowbite-svelte';
  import EditorToolbar from './EditorToolbar.svelte';

  let { value = $bindable(''), rows = 10, name = '' } = $props();
  let element = $state(null);
</script>

<div>
  <EditorToolbar target={element} />
  <textarea
    bind:this={element}
    bind:value
    {name}
    {rows}
    class="h-48 w-full my-2 p-2 border border-gray-300 focus:border-indigo-500 focus:ring-indigo-500 rounded-md shadow-sm bg-gray-50 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-indigo-500 dark:focus:border-indigo-500"
  ></textarea>
</div>
