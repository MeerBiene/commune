<script>
import { navigate, useLocation } from "svelte-navigator";
import {onMount} from 'svelte'
import tippy from 'tippy.js';
import {store} from '../../store/store.js'

import { logo, caption } from './logo.js'

const location = useLocation()

$: first = $location.pathname.split("/")[1]
$: second = $location.pathname.split("/")[2]

let indicate = false;

$: {
    indicate = $location.pathname == '/' || first == 'messages'
}


let container;

onMount(() => {
})

$: if(container) {
    tip()
}

let tooltip;

let tip = () => {
    if(container != undefined) {
    tooltip = tippy(container, {
        content: 'Home',
        theme: 'tooltip',
        placement: 'right',
        offset: [0, 24],
        delay: 0,
        duration: 0,
    });
    }
}

$: expanded = $store.settings.switcher.mode == 'expanded'
$: collapsed = $store.settings.switcher.mode == 'collapsed'

$: if(expanded && tooltip) {
    tooltip.destroy()
}

$: if(collapsed) {
    tip()
}

$: width = collapsed ? `16` : expanded ? `18` : `24`



$: home = $store.accounts.filter(account => account.user_id == $store.active_account)[0]?.home

function go() {
    navigate(`${home?.active_page}`)
}
</script>

<div class="switcher-root">


    <div class="switcher-home pointer" on:click={go} 
    class:switcher-home-expanded={expanded}>
    <div class="switcher-home-container gr-center relative"
    class:collapsed={collapsed}
    class:expanded={expanded}
    class:mute-home-container={!indicate}
    bind:this={container}>
        <div class="switcher-home-logo gr-center">
            <div class="logo" style="width:{width};">
                {@html logo}
            </div>
        </div>

        <div class="switcher-indicator" 
            class:expanded-indicator={expanded}
            class:collapsed-indicator={collapsed}
            class:show-indicator={indicate}>
        </div>


    </div>

    {#if expanded}
        <div class="gr-start-center">
            <div class="caption" class:inactive={!indicate}>
                {@html caption}
            </div>
        </div>


    {/if}



</div>

</div>


<style>
.switcher-root {
    display:grid;
    grid-template-columns: 1fr auto;
}

.switcher-home {
    display: grid;
    padding-top: 0.5rem;
    padding-bottom: 0.5rem;
    border-radius: 50%;
}

.switcher-home-container {
    background-color: var(--dark);
    display: grid;
    border-radius: 30%;
    height: 48px;
    width: 48px;
}

.switcher-home:hover .switcher-indicator {
    opacity: 1;
    height: 28px;
    top: 10px;
    margin-left: -13px;
}

.switcher-home-expanded {
    grid-template-columns: [icon] 52px [content] 1fr [settings] auto;
}

.logo {
    fill: var(--primary);
}

.mute-home-container .logo {
    fill: var(--primary-muted);
}

.collapsed {
    height: 32px;
    width: 32px;
}

.collapsed-indicator{
    margin-left: -10px;
    height: 20px;
    top: 6;
}

.switcher-home:hover .collapsed-indicator {
    margin-left: -10px;
    height: 20px;
    top: 6;
}

.expanded {
    height: 32px;
    width: 32px;
}

.expanded-indicator{
    margin-left: -12px;
    height: 20px;
    top: 6;
}

.switcher-home:hover .expanded-indicator {
    margin-left: -12px;
    height: 20px;
    top: 6;
}

.caption {
    margin-top: 2px;
    width: 75px;
    fill: var(--white);
}

.inactive {
    opacity: 0.5;
}

.mute-home-container {
    background-color: var(--background-3);
}

.mute-home-container:hover {
}

.switcher-home:hover .logo{
}

.switcher-item:hover .mute-home-container {
    background-color: #0f0f0f;
}

svg {
    fill: var(--background-4);
}

svg:hover {
    fill: var(--primary-muted);
}

.icon {
    width: 32px;
    height: 32px;
}

.icon:hover {
    background-color: var(--background-2);
    border-radius: 7px;
}
.icon:hover {
    fill: var(--white);
}
</style>
