<script>
import { store } from '../../../store/store.js'
import { gearSmall } from '../../../utils/icons.js'
import {onMount} from 'svelte'
import tippy from 'tippy.js';

$: expanded = $store.settings.switcher.mode == 'expanded'
$: collapsed = $store.settings.switcher.mode == 'collapsed'

let container;

onMount(() => {
    if(!expanded) {
        tip();
    }
})

let tooltip;
let tip = () => {
    if(!container) {
        return
    }
    tooltip = tippy(container, {
        content: 'Settings',
        theme: 'tooltip-html',
        placement: 'right',
        offset: [0, 12],
        delay: 0,
        duration: 0,
    });
}

$: if(expanded && tooltip) {
    tooltip.destroy()
}

$: if(collapsed) {
    tip()
}

function toggleSettings(e) {
    e.preventDefault()
    window.toggleSettings()
}
</script>

<div class="settings mt1" on:click={toggleSettings}>
<div class="item gr-default pv1" 
    bind:this={container}
    class:expanded={expanded}
    class:collapsed={collapsed}>
    <div class="ic gr-center gr-default">
        {@html gearSmall}
    </div>
    {#if expanded}
        <div class="gr-start-center">
            Settings
    </div>
    {/if}
</div>
</div>

<style>
.item {
    display: grid;
    grid-template-columns: 72px auto;
    grid-template-rows: auto;
    cursor: pointer;
}

.collapsed {
    grid-template-columns: 48px auto;
}
.expanded {
    grid-template-columns: 52px auto;
}
.ic {
    background-color: var(--background-2);
    border-radius: 50%;
    width: 48px;
    height: 48px;
    fill: var(--text);
}

.item:hover .ic {
    border-radius: 30%;
    background-color: var(--secondary);
    fill: var(--white);
}

.collapsed .ic {
    width: 32px;
    height: 32px;
}
.expanded .ic {
    width: 32px;
    height: 32px;
}
</style>
