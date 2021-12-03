let root = document.querySelector(`.root`)
if(root) {
  import('./app.svelte').then(res => {
      new res.default({
          target: root,
          hydrate: true
      });
  })
}

