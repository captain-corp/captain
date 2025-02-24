<script>
  let { logoUrl = $bindable(''), onDelete = () => {} } = $props();
</script>

<div class="mb-4 flex justify-center items-center flex-col">
  <div class="w-32 h-32 object-cover rounded">
    <img src={logoUrl} alt="Logo" class="w-32 h-32 object-cover rounded" />
  </div>
  <button
    type="button"
    class="mt-4 inline-flex justify-center px-4 py-2 text-sm font-medium text-red-700 hover:text-red-900"
    onclick={() => onDelete()}
  >
    Delete
  </button>
</div>
