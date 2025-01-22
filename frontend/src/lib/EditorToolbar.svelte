<script lang="ts">
  import { Dropdown, DropdownItem } from 'flowbite-svelte';
  import { LinkOutline, ImageSolid, FileCodeOutline, QuoteSolid } from 'flowbite-svelte-icons';
  let { target } = $props();

  function surroundSelection(before: string, after: string, text: string): void {
    let selectiondEnd = target.selectionEnd + after.length;
    const commandResult = `${before}${text}${after}`;

    target.setRangeText(commandResult, target.selectionStart, target.selectionEnd);
    target.focus();

    if (text === '') {
      selectiondEnd = target.selectionStart + after.length;
    }

    target.setSelectionRange(target.selectionStart + before.length, selectiondEnd);
  }

  function surroundLink(text: string, isImage: boolean = false): void {
    let selectionStart = target.selectionStart;
    const textPrefix = isImage ? '![' : '[';
    const textSuffix = '](';
    const urlPlaceholder = 'url';
    const textInjected = `${textPrefix}${text}${textSuffix}${urlPlaceholder})`;

    target.setRangeText(textInjected);
    target.focus();

    let cursorPosition = selectionStart + textPrefix.length + text.length;

    if (text === '') {
      target.setSelectionRange(cursorPosition, cursorPosition);
    } else {
      selectionStart = cursorPosition + textSuffix.length;
      target.setSelectionRange(selectionStart, selectionStart + urlPlaceholder.length);
    }
  }

  function injectBeforeCursor(before: string, text: string): void {
    const cursorPosition = target.selectionStart;
    const newCursorPosition = cursorPosition + before.length + text.length;
    target.setRangeText(before + ' ' + text);
    target.setSelectionRange(newCursorPosition, newCursorPosition + 1);
    target.focus();
  }

  const commands = {
    Bold: (text: string): void => {
      surroundSelection('**', '**', text);
    },
    Italic: (text: string): void => {
      surroundSelection('*', '*', text);
    },
    Link: (text: string): void => {
      surroundLink(text, false);
    },
    Image: (text: string): void => {
      surroundLink(text, true);
    },
    Code: (text: string): void => {
      surroundSelection('`', '`', text);
    },
    Quote: (text: string): void => {
      injectBeforeCursor('> ', text);
    },
    H1: (text: string): void => {
      injectBeforeCursor('# ', text);
    },
    H2: (text: string): void => {
      injectBeforeCursor('## ', text);
    },
    H3: (text: string): void => {
      injectBeforeCursor('### ', text);
    },
    H4: (text: string): void => {
      injectBeforeCursor('#### ', text);
    },
    H5: (text: string): void => {
      injectBeforeCursor('##### ', text);
    },
    H6: (text: string): void => {
      injectBeforeCursor('###### ', text);
    },
  };

  function getSelectedText() {
    const start = target.selectionStart;
    const end = target.selectionEnd;
    const selectedText = target.value.substring(start, end);
    return selectedText;
  }

  function executeCommand(event: Event) {
    event.preventDefault();
    const selectedText = getSelectedText();

    const command = (event.target as HTMLButtonElement).dataset['command'];
    commands[command](selectedText);
  }
</script>

<div>
  <button
    data-command="Bold"
    onclick={executeCommand}
    class="font-bold border border-black dark:border-gray-400 text-black dark:text-white">B</button
  >
  <button
    data-command="Italic"
    onclick={executeCommand}
    class="font-italic border border-black dark:border-gray-400 text-black dark:text-white"
    >I</button
  >
  <button class="border border-black dark:border-gray-400 text-black dark:text-white">H</button>
  <Dropdown>
    <DropdownItem data-command="H1" on:click={executeCommand}>Heading 1</DropdownItem>
    <DropdownItem data-command="H2" on:click={executeCommand}>Heading 2</DropdownItem>
    <DropdownItem data-command="H3" on:click={executeCommand}>Heading 3</DropdownItem>
    <DropdownItem data-command="H4" on:click={executeCommand}>Heading 4</DropdownItem>
    <DropdownItem data-command="H5" on:click={executeCommand}>Heading 5</DropdownItem>
    <DropdownItem data-command="H6" on:click={executeCommand}>Heading 6</DropdownItem>
  </Dropdown>
  <button
    data-command="Link"
    onclick={executeCommand}
    class="border border-black dark:border-gray-400 text-black dark:text-white"
  >
    <LinkOutline class="w-4 h-6" />
  </button>
  <button
    data-command="Image"
    onclick={executeCommand}
    class="border border-black dark:border-gray-400 text-black dark:text-white"
  >
    <ImageSolid class="w-4 h-6" />
  </button>
  <button
    data-command="Code"
    onclick={executeCommand}
    class="border border-black dark:border-gray-400 text-black dark:text-white"
  >
    <FileCodeOutline class="w-4 h-6" />
  </button>
  <button
    data-command="Quote"
    onclick={executeCommand}
    class="border border-black dark:border-gray-400 text-black dark:text-white"
  >
    <QuoteSolid class="w-4 h-6" />
  </button>
</div>

<style>
  button {
    margin-right: 5px;
    padding: 5px 10px;
    border-radius: 4px;
    vertical-align: middle;
  }
</style>
