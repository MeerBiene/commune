<script>
import { navigate, useLocation } from "svelte-navigator";
import {onMount} from 'svelte'
import tippy from 'tippy.js';
import { store } from '../../store/store.js'

$: expanded = $store.settings.switcher.mode == 'expanded'
$: collapsed = $store.settings.switcher.mode == 'collapsed'

let container;

export let activeServer;
export let server;
export let room;
export let killTooltip;

let go = () => {
    navigate(activeServer.active_room)
    if(room && server) {
        store.updateActiveRooms(`/${server}`, `/${server}/${room}`)
        //store.updateActiveRoom(`/${server}`, `/${server}/${room}`)
        setTimeout(() => {
            store.loadIdleRoomsEvents(serverpath)
        }, 1000)
    }
}

onMount(() => {
    tip()
})

let tooltip;
let tip = () => {
    tooltip = tippy(container, {
        content: tipcontent.innerHTML,
        theme: 'tooltip-html',
        placement: 'right',
        allowHTML: true,
        offset: [0, 12],
        delay: 0,
        duration: 0,
    });
}


$: hasAvatar = activeServer?.avatar?.length > 0

//$: avatar = server?.avatar
$: avatar = `${homeServer}/_matrix/media/r0/download/${activeServer?.avatar?.substring(6)}`


const location = useLocation()

let indicate = false;

$: serverAlias = $location.pathname.split("/")[1]

$: {
    indicate = activeServer.pathname == `/${serverAlias}`
}

$: locationPathname = $location.pathname.split("/")[1]

$: serverpath = `/${locationPathname}`

let idleRoomsLoaded;

$: if(indicate) {
    store.updateActiveRooms(`/${server}`, `/${server}/${room}`)
}

$: activeRoom = $store.allRooms.filter(x => x.alias == server)[0]

/*
$: if(indicate && container) {
        container.scrollIntoView()
}
*/


$: if(killTooltip) {
    if(tooltip) {
        tooltip.hide()
    }
}

function tit(name) {
    let newtitle = '';
    const splut = name?.split(" ")
    splut?.forEach(x => {
        newtitle = newtitle + x.substring(0,1)
    })
    return newtitle.toLowerCase().substring(0,4)
}

$: title = tit(activeServer.name)

let tipcontent;

function rightClick(e) {
    e.preventDefault()
    console.log(activeServer)
}

</script>

<template bind:this={tipcontent}>
    <div class="fl-co">
        <div class="bold">
            {activeServer.name}
        </div>
    </div>
</template>

<div class:switcher-item-root={expanded}
    class:selected-expanded={indicate && expanded}
    bind:this={container} 
    draggable="true"
    on:click={go}
    on:contextmenu={rightClick}>

<div class="item-root">
<div class="switcher-item pv1"
    class:selected={indicate}
    class:switcher-root-expanded={expanded}>
        <div class="switcher-item-container gr-center relative"
            class:switcher-item-active={indicate}
            class:expanded={expanded}
            class:collapsed={collapsed}
            class:has-avatar={hasAvatar}
            class:no-avatar={!hasAvatar}>

            {#if hasAvatar}
                <div class="bg-img avatar" 
                    style="background-image: url({avatar});">
                </div>
            {:else}
                <div class="title gr-center smaller bold" 
                    class:smallestest={(collapsed && title.length > 3) || (expanded && title.length > 3)}
                    class:smallest={collapsed || (expanded && title.length > 2)}>
                    {title}
                </div>
            {/if}

        <div class="switcher-indicator" 
            class:collapsed-indicator={collapsed}
            class:expanded-indicator={expanded}
            class:show-indicator={indicate}>
        </div>

        </div>

{#if expanded}
    <div class="name gr-start-center fl-co lh-copy pr2 clmp-1"
        class:name-active={indicate}>
        {activeServer.name}
    </div>
{/if}


</div>

</div>

{#if expanded}
    <div class="item-tools gr-default">
        <div class="gr-start-center mr2 pointer" >
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16"><path d="M8 9a1.5 1.5 0 100-3 1.5 1.5 0 000 3zM1.5 9a1.5 1.5 0 100-3 1.5 1.5 0 000 3zm13 0a1.5 1.5 0 100-3 1.5 1.5 0 000 3z"></path></svg>
        </div>
    </div>
{/if}
</div>

<style>

.switcher-item-root {
    display: grid;
    grid-template-columns: 1fr auto;
}

.switcher-item-root:hover {
}

.switcher-root-expanded {
    grid-template-columns: [icon] 52px [content] 1fr;
}

.selected-expanded {
}


.switcher-expanded:hover {
    background-color: var(--background-5);
}

.collapsed {
    height: 32px;
    width: 32px;
}

.switcher-item:hover .collapsed {
}

.expanded {
    height: 32px;
    width: 32px;
}

.collapsed-indicator{
    margin-left: -10px;
    height: 20px;
    top: 6;
}

.switcher-item:hover .collapsed-indicator {
    margin-left: -10px;
    height: 20px;
    top: 6;
}

.expanded-indicator{
    margin-left: -12px;
    height: 20px;
    top: 6;
}

.switcher-item:hover .expanded-indicator {
    margin-left: -12px;
    height: 20px;
    top: 6;
}

.small {
    font-size: 0.7rem;
}

.title {
    text-transform: uppercase;
}

.name {
    color: var(--text-light);
}

.name-active {
    color: var(--text);
    color: var(--white);
}


a {
    color: var(--text-muted);
}

.switcher-item:hover .name {
    opacity: 1;
}

.avatar {
    border-radius: 50%;
    background-color: var(--background-3)
}

.selected .avatar {
    border-radius: 30%;
}

.switcher-item-active{
    border-radius: 30%;
    background-color: var(--primary);
    color: black;
}

.has-avatar {
    background-color: transparent;
}

.smallestest {
    font-size: 0.46rem;
}


svg {
    fill: var(--primary-muted);
}

.item-tools {
    opacity: 0;
}

.switcher-item-root:hover .item-tools {
    opacity: 1;
}

</style>
