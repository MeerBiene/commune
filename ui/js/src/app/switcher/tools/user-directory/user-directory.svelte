<script>
import {onMount} from 'svelte'
import { navigate, useLocation } from "svelte-navigator";
import { user } from '../../../utils/icons.js'
import tippy from 'tippy.js';

export let expanded;
export let collapsed;

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
        content: 'Explore Spaces',
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

const location = useLocation()

let indicate = false;

$: {
    indicate = $location.pathname.substring(0,2) == `/@` || 
        $location.pathname == `/users`
}

function go() {
    navigate(`/users`, {replace: true})
}
</script>

<div class="explore" on:click={go}>
<div class="item gr-default pv1" 
    bind:this={container}
    class:expanded={expanded}
    class:selected={indicate}
    class:collapsed={collapsed}>
    <div class="ic gr-center gr-default">
        <div class="u-icon gr-center">
            {@html user}
        </div>
    </div>
    {#if expanded}
        <div class="gr-start-center" class:selected={indicate}>
            User Directory
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
.selected {
    color: var(--white);
}
.u-icon {
    width: 20px;
}
</style>
