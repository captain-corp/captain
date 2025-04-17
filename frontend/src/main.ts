import './app.css';

import Posts from './apps/Posts.svelte';
import Pages from './apps/Pages.svelte';
import Settings from './apps/Settings.svelte';
import MediaPicker from './lib/MediaPicker.svelte';
import { Inity } from './lib/inity';

const Apps = {
  Posts,
  Pages,
  Settings,
  MediaPicker,
};

Object.assign(window, {
  Inity,
  Apps,
  MediaPicker,
});
