<script>
import {onMount} from 'svelte'
import tippy from 'tippy.js';
import {store} from '../../../store/store.js'
import View from './view.svelte'

export let expanded;
export let collapsed;

let container;

onMount(() => {
    if(!expanded && container) {
        tip();
    }
    window.addNewServer = () => {
        active = true
    }
})

let tooltip;
let tip = () => {
    if(container) {
    tooltip = tippy(container, {
        content: 'Add New Space',
        theme: 'tooltip-html',
        placement: 'right',
        offset: [0, 12],
        delay: 0,
        duration: 0,
    });
    }
}

$: if(expanded && tooltip) {
    tooltip.destroy()
}

$: if(collapsed && container) {
    tip()
}

let active = false;

function toggle() {
    active = !active
}

</script>

<div class="item gr-default pv1" 
    bind:this={container}
    class:expanded={expanded}
    class:collapsed={collapsed} 
    on:click={toggle}>
    <div class="ic gr-center gr-default">
        <svg class="gr-center" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16"><path fill-rule="evenodd" d="M7.75 2a.75.75 0 01.75.75V7h4.25a.75.75 0 110 1.5H8.5v4.25a.75.75 0 11-1.5 0V8.5H2.75a.75.75 0 010-1.5H7V2.75A.75.75 0 017.75 2z"></path></svg>
    </div>
    {#if expanded}
    <div class="gr-start-center">
        Add New Space
    </div>
    {/if}
</div>

{#if active}
    <View on:kill={toggle}/>
{/if}

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
