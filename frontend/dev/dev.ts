import "../src/app.css";


import Posts from "../src/apps/Posts.svelte";
import Pages from "../src/apps/Pages.svelte";
import Settings from "../src/apps/Settings.svelte";
import MediaPicker from "../src/lib/MediaPicker.svelte";
import { Inity } from "../src/lib/inity";

Inity.register("posts", Posts, {
  onSubmit: (data, done, error) => {
    console.log("Submitted post data:", data);
    done("saved");
    error('Error');
  },
});

Inity.register("pages", Pages, {
  onSubmit: (data, done, error) => {
    console.log("Submitted page data:", data);
    done("saved");
    error('Error');
  },
});

Inity.register("settings", Settings);

Inity.register("media-picker", MediaPicker, {
  mediaProvider: (done, error) => {
    setTimeout(() => {
      done([
        { alt: 'shoes', src: 'https://flowbite.s3.amazonaws.com/docs/gallery/square/image-1.jpg' },
        { alt: 'shoes', src: 'https://flowbite.s3.amazonaws.com/docs/gallery/square/image-2.jpg' },
        { alt: 'shoes', src: 'https://flowbite.s3.amazonaws.com/docs/gallery/square/image-3.jpg' },
        { alt: 'shoes', src: 'https://flowbite.s3.amazonaws.com/docs/gallery/square/image-4.jpg' },
        { alt: 'shoes', src: 'https://flowbite.s3.amazonaws.com/docs/gallery/square/image-5.jpg' },
        { alt: 'shoes', src: 'https://flowbite.s3.amazonaws.com/docs/gallery/square/image-6.jpg' },
        { alt: 'shoes', src: 'https://flowbite.s3.amazonaws.com/docs/gallery/square/image-7.jpg' },
        { alt: 'shoes', src: 'https://flowbite.s3.amazonaws.com/docs/gallery/square/image-8.jpg' },
        { alt: 'shoes', src: 'https://flowbite.s3.amazonaws.com/docs/gallery/square/image-9.jpg' },
      ]);
    }, 1000);
  },
});

Inity.attach();
