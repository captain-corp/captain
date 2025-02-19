import { type SavingStates } from './common';

import type { SelectOptionType } from 'flowbite-svelte';

export interface Settings {
  title?: string;
  subtitle?: string;
  logoID?: number;
  useLogoAsFavicon?: boolean;
  codeHighlightTheme?: string;
  codeThemeOptions: SelectOptionType<String>[];
  postsPerPage?: number;
  onSubmit?: (
    data: any,
    done: (savingState: SavingStates) => void,
    error: (error: any) => void
  ) => void;
}
