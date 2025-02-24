import { type SavingStates } from './common';

import type { SelectOptionType } from 'flowbite-svelte';

export interface Settings {
  title?: string;
  subtitle?: string;
  logoUrl?: string;
  useLogoAsFavicon?: boolean;
  codeHighlightTheme?: string;
  codeThemeOptions: SelectOptionType<String>[];
  postsPerPage?: number;
  deleteLogoHandler?: () => void;
  uploadLogoHandler?: (
    files: File[],
    uploadStarted: () => void,
    uploadFinished: (error: string | null, url: string | null) => void
  ) => void;
  onSubmit?: (
    data: any,
    done: (savingState: SavingStates) => void,
    error: (error: any) => void
  ) => void;
}
