<script>
import { store } from '../../store/store.js'
import tippy from 'tippy.js';
import {onMount, tick} from 'svelte'

$: flip = $store.settings.switcher.mode == 'expanded'

async function toggle() {
    store.toggleSwitcher()
    await tick()
    tip()
}

$: normal = $store.settings.switcher.mode == 'normal'
$: expanded = $store.settings.switcher.mode == 'expanded'
$: collapsed = $store.settings.switcher.mode == 'collapsed'

$: content = expanded ? `Collapse` : collapsed ? `Expand` : `Expand`

let container;
onMount(() => {
    tip()
})

let tooltip;

let tip = () => {
    if(tooltip) {
        tooltip.destroy()
    }
    tooltip = tippy(container, {
        content: content,
        theme: 'tooltip',
        placement: 'right',
        offset: [0, 20],
    });
}

</script>

<div class="switcher-footer" 
    class:collapsed={collapsed}
    class:expanded={expanded}>
    <div class="toggle-switcher gr-center gr-default">
        <div class="gr-center gr-default pointer toggler" 
            class:toggler-normal={normal}
            class:toggler-expanded={expanded}
            bind:this={container}
            on:click={toggle}>
            <svg class="gr-center" class:flip-h={flip} xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16"><path fill-rule="evenodd" d="M3.72 3.72a.75.75 0 011.06 1.06L2.56 7h10.88l-2.22-2.22a.75.75 0 011.06-1.06l3.5 3.5a.75.75 0 010 1.06l-3.5 3.5a.75.75 0 11-1.06-1.06l2.22-2.22H2.56l2.22 2.22a.75.75 0 11-1.06 1.06l-3.5-3.5a.75.75 0 010-1.06l3.5-3.5z"></path></svg>
        </div>
    </div>
</div>

<style>

.switcher-footer {
    display: grid;
    grid-template-rows: [content] 52px;
    grid-template-columns: 100%;
}

.toggle-switcher {
    background-color: var(--background-1);
    width: 100%;
    height: 100%;
}


.toggler {
    background-color: var(--background-2);
    border-radius: 50%;
    width: 32px;
    height: 32px;
}

.toggler:hover {
    background-color: var(--background-3);
}

.toggler-normal {
    border-radius: 500px;
    width: 48px;
}

.toggler-expanded {
    border-radius: 500px;
    width: 200px;
    height: 18px;
}


.switcher-footer {
}

.expanded {
}

.collapsed {
}

svg {
    fill: var(--background-4);
}

svg:hover {
    fill: var(--primary-muted);
}

.toggler:hover svg {
    fill: var(--primary-muted);
}
</style>
