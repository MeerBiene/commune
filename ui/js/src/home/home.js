let root = document.querySelector(`.auth`)
if(root) {
  import('./auth/auth.svelte').then(res => {
      new res.default({
          target: root,
          hydrate: true
      });
  })
}

