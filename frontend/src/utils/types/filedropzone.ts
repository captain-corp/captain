export interface FileDropzone {
  accept?: string;
  multiple?: boolean;
  maxSize?: number;
  onFilesSelectedHandler?: (
    files: File[],
    uploadStarted: () => void,
    uploadFinished: (error: string | null, url: string | null) => void
  ) => void;
}
