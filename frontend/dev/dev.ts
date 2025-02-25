import "../src/app.css";


import Posts from "../src/apps/Posts.svelte";
import Pages from "../src/apps/Pages.svelte";
import Settings from "../src/apps/Settings.svelte";
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

Inity.register("settings", Settings, {
  onSubmit: (data, done, error) => {
    console.log("Submitted settings data:", data);
    done("saved");
    error('Error');
  },
});


Inity.attach();
