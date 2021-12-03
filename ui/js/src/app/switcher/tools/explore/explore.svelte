<script>
import {onMount} from 'svelte'
import { navigate, useLocation } from "svelte-navigator";
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
    indicate = $location.pathname == `/explore`
}

function go() {
    navigate(`/explore`, {replace: true})
}
</script>

<div class="explore" on:click={go}>
<div class="item gr-default pv1" 
    bind:this={container}
    class:expanded={expanded}
    class:selected={indicate}
    class:collapsed={collapsed}>
    <div class="ic gr-center gr-default">
        <svg class="gr-center" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16"><path fill-rule="evenodd" d="M1.543 7.25h2.733c.144-2.074.866-3.756 1.58-4.948.12-.197.237-.381.353-.552a6.506 6.506 0 00-4.666 5.5zm2.733 1.5H1.543a6.506 6.506 0 004.666 5.5 11.13 11.13 0 01-.352-.552c-.715-1.192-1.437-2.874-1.581-4.948zm1.504 0h4.44a9.637 9.637 0 01-1.363 4.177c-.306.51-.612.919-.857 1.215a9.978 9.978 0 01-.857-1.215A9.637 9.637 0 015.78 8.75zm4.44-1.5H5.78a9.637 9.637 0 011.363-4.177c.306-.51.612-.919.857-1.215.245.296.55.705.857 1.215A9.638 9.638 0 0110.22 7.25zm1.504 1.5c-.144 2.074-.866 3.756-1.58 4.948-.12.197-.237.381-.353.552a6.506 6.506 0 004.666-5.5h-2.733zm2.733-1.5h-2.733c-.144-2.074-.866-3.756-1.58-4.948a11.738 11.738 0 00-.353-.552 6.506 6.506 0 014.666 5.5zM8 0a8 8 0 100 16A8 8 0 008 0z"></path></svg>
    </div>
    {#if expanded}
        <div class="gr-start-center" class:selected={indicate}>
        Explore Spaces
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
</style>
