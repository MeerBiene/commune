<script>
import {store} from '../../store/store.js'
import { navigate } from "svelte-navigator";
import {user} from '../../utils/icons.js'
import tippy from 'tippy.js';
import {onMount} from 'svelte'

export let room_id;
export let count;


let container;
$: expanded = $store.settings.switcher.mode == 'expanded'
$: collapsed = $store.settings.switcher.mode == 'collapsed'

$: account = $store.accounts.filter(x => x.user_id == $store.active_account)[0]

$: room = account.direct_messages.filter(x => x.room_id == room_id)[0]

$: name = room?.name

$: hasAvatar = room?.avatar_url?.length > 0
$: avatar = `${homeServer}/_matrix/media/r0/download/${room?.avatar_url?.substring(6)}`

let go = () => {
    navigate(`/messages/${room.alias}`, { replace: true })
    store.updateActiveHomePage(`/messages/${room.alias}`)
    store.updateActiveDirectMessages(`/messages/${room.alias}`)
    setTimeout(() => {
        store.resetDMnotification(room_id)
    }, 1000)
}

onMount(() => {
    tip()
})

let tooltip;
let tipcontent;
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


$: notMe = Object.fromEntries(
    Object.entries(room?.members).filter(([key, _]) => key != $store.active_account) )

$: others = Object.entries(notMe)

$: multiple = others?.length > 1

</script>

<template bind:this={tipcontent}>
    <div class="fl-co">
        {#each Object.entries(notMe) as [user_id, member] (user_id)}
        <div class="bold">
            {user_id}
        </div>
        {/each}
    </div>
</template>

<div class:switcher-item-root={expanded}
    bind:this={container} 
    on:click={go}>

<div class="item-root">
<div class="switcher-item pv1"
    class:switcher-root-expanded={expanded}>
        <div class="switcher-item-container gr-center relative"
            class:expanded={expanded}
            class:collapsed={collapsed}
            class:has-avatar={hasAvatar}
            class:no-avatar={!hasAvatar}>

            {#if hasAvatar}
                <div class="profile-avatar ncol bg-img"
                    class:pa-d={!expanded && !collapsed}
                    style="background-image: url({avatar});">
                </div>
            {:else}
                <div class="profile-avatar gr-default"
                    class:pa-d={!expanded && !collapsed}>
                    <div class="log gr-center">
                        {@html user}
                    </div>
                </div>
            {/if}


<div class="dm-count gr-default"
    class:dm-c={collapsed || expanded}>
    <div class="gr-center">
        {count}
    </div>
</div>

        <div class="indicator" 
            class:collapsed-indicator={collapsed}
            class:normal-indicator={!collapsed && !expanded}
            class:expanded-indicator={expanded}>
        </div>

        </div>

{#if expanded}
    {#if multiple}
    <div class="flex flex-column">
        <div class="name fl-co pr2 clmp-1 gr-center-start lh">
            {name}
        </div>
            <div class="mem mute lh">
                {others.length} members
            </div>
    </div>
    {:else}
        <div class="name fl-co pr2 clmp-1 gr-start-center lh">
            {name}
        </div>
    {/if}
{/if}


</div>


</div>

</div>



<style>

.dm-count {
    width: 22px;
    height: 22px;
    background: red;
    border-radius: 50%;
    border: 3px solid var(--background-1);
    position :absolute;
    right: -3px;
    bottom: -3px;
    font-size: 0.7rem;
    font-weight: bold;
    color: var(--white);
    line-height: 0.7rem;
}

.dm-c {
    width: 18px;
    height: 18px;
    right: -5px;
    bottom: -5px;
    font-size: 0.6rem;
    line-height: 0.6rem;
}

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

.indicator {
    display: grid;
    width: 5px;
    height: 0px;
    border-radius: 0 3px 3px 0;
    position: absolute;
    top: 20px;
    margin-left: -33px;
    background-color: var(--white);
}


.collapsed-indicator{
    margin-left: -10px;
    height: 10px;
    top: 12;
}

.switcher-item:hover .collapsed-indicator {
    margin-left: -10px;
    height: 20px;
    top: 6;
}

.normal-indicator{
    margin-left: -12px;
    height: 10px;
    top: 20;
}

.switcher-item:hover .normal-indicator {
    margin-left: -12px;
    height: 20px;
    top: 14;
}


.expanded-indicator{
    margin-left: -12px;
    height: 10px;
    top: 12;
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

.profile-avatar {
    width: 32px;
    height: 32px;
    background-color: var(--avatar);
    border-radius: 50%;
}
.ncol {
    background-color: transparent;
}

.pa-d {
    width: 48px;
    height: 48px;
}

.log {
    fill: var(--text);
    width: 16px;
    height: 16px;
}
.mem {
    margin-top: 0.25rem;
    font-size: 0.8rem;
}
.lh {
    line-height: 1;
}
</style>
