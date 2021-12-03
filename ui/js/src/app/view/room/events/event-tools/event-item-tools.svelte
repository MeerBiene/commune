<script>
import { store } from '../../../../store/store.js'
import { copy } from '../../../../utils/icons.js'
import { onMount, createEventDispatcher } from 'svelte'
const dispatch = createEventDispatcher()
import { makeid } from '../../../../utils/utils.js'
import ToolsMenu from './tools-menu.svelte'

export let room;
export let event;
export let owner;

export let thread;

export let hasCode;
export let userReactions;


function edit() {
    dispatch('edit', true)
}

function reply() {
    dispatch('reply', true)
}

let container;
let Emoji;
let emojiLoaded;
let content;
let menu;

onMount(() => {
})

function createThread() {
    store.newThread(room, event)
}

let reacting = false;

function react() {
    if(reacting) {
        killed()
        return
    }
    reacting = true
    let opts = {
        room: room,
        event: event,
        container: container,
        userReactions: userReactions,
    }
    window.toggleEmojiPicker(opts)
}

$: if(emojiActive) {
    dispatch('highlight', true)
} else {
    if(!showingMenu) {
        dispatch('highlight', false)
    }
}


$: emojiActive = $store.emoji?.active && $store.emoji?.event == event.event_id

$: if(!emojiActive) {
    killed()
}

function killed() {
    reacting = false
    //dispatch('highlight', emojiActive)
}

$: matrix = $store.accounts.filter(account => account.user_id == $store.active_account)[0]?.matrix

function reacted(e) {
    let tempid = makeid(16)

    let content = {
        'm.relates_to': {
            event_id: event.event_id,
            key: e.detail,
            rel_type: 'm.annotation',
            "client_ts": tempid,
        }
    }


    let newEvent= {
        "type": "m.reaction",
        "room_id": room.room_id,
        "sender": $store.active_account,
        "content": content,
        "origin_server_ts": new Date(),
        "unsigned": {
            "age": 0
        },
        "event_id": tempid,
        "user_id": $store.active_account,
        "delivered": false,
        "client_ts": tempid,
    }
    console.log(newEvent)

    //store.updateReactions(room.room_id, newEvent)
    store.addEventToRoom(room.room_id, newEvent)

    matrix.sendEvent(room.room_id, "m.reaction", content, "", (err, res) => {
        console.log(res)
    });
    dispatch('reacted', e.detail)

}

let copied = false;

$: tip = copied ? `Copied` : 'Copy'

function copyText() {
    copied = true
    setTimeout(() => {
        copied = false
    }, 2000)
}

let showingMenu = false;

function showMenu() {
    if(!showingMenu) {
        dispatch('highlight', true)
    }
    showingMenu = true
}

function hideMenu() {
    if(!reacting && showingMenu) {
        dispatch('highlight', false)
    }
    showingMenu = false
}

</script>



{#if hasCode}
    <div class="gr-center et-ic icon"
        on:click={copyText}
        aria-label={tip}
        data-microtip-position="top"
        data-microtip-size="fit"
        role="tooltip">
        {@html copy}
    </div>
{/if}

    <div class="gr-center et-ic"
        bind:this={container}
        on:click={react}
        aria-label="Add Reaction"
        data-microtip-position="top"
        data-microtip-size="fit"
        role="tooltip">
        <svg class="icon-3Gkjwa" aria-hidden="false" width="24" height="24" viewBox="0 0 24 24"><path fill="currentColor" fill-rule="evenodd" clip-rule="evenodd" d="M12.2512 2.00309C12.1677 2.00104 12.084 2 12 2C6.477 2 2 6.477 2 12C2 17.522 6.477 22 12 22C17.523 22 22 17.522 22 12C22 11.916 21.999 11.8323 21.9969 11.7488C21.3586 11.9128 20.6895 12 20 12C15.5817 12 12 8.41828 12 4C12 3.31052 12.0872 2.6414 12.2512 2.00309ZM10 8C10 6.896 9.104 6 8 6C6.896 6 6 6.896 6 8C6 9.105 6.896 10 8 10C9.104 10 10 9.105 10 8ZM12 19C15.14 19 18 16.617 18 14V13H6V14C6 16.617 8.86 19 12 19Z"></path><path d="M21 3V0H19V3H16V5H19V8H21V5H24V3H21Z" fill="currentColor"></path></svg>
    </div>

    {#if owner}
        <div class="gr-center et-ic"
            on:click={edit}
            aria-label="Edit"
            data-microtip-position="top"
            data-microtip-size="fit"
            role="tooltip">
            <svg class="icon-3Gkjwa" aria-hidden="false" width="16" height="16" viewBox="0 0 24 24"><path fill-rule="evenodd" clip-rule="evenodd" d="M19.2929 9.8299L19.9409 9.18278C21.353 7.77064 21.353 5.47197 19.9409 4.05892C18.5287 2.64678 16.2292 2.64678 14.817 4.05892L14.1699 4.70694L19.2929 9.8299ZM12.8962 5.97688L5.18469 13.6906L10.3085 18.813L18.0201 11.0992L12.8962 5.97688ZM4.11851 20.9704L8.75906 19.8112L4.18692 15.239L3.02678 19.8796C2.95028 20.1856 3.04028 20.5105 3.26349 20.7337C3.48669 20.9569 3.8116 21.046 4.11851 20.9704Z" fill="currentColor"></path></svg>
        </div>
    {:else}
        <div class="gr-center et-ic"
            on:click={reply}
            aria-label="Reply"
            data-microtip-position="top"
            data-microtip-size="fit"
            role="tooltip">
            <svg class="icon-3Gkjwa" width="24" height="24" viewBox="0 0 24 24"><path d="M10 8.26667V4L3 11.4667L10 18.9333V14.56C15 14.56 18.5 16.2667 21 20C20 14.6667 17 9.33333 10 8.26667Z" fill="currentColor"></path></svg>
        </div>
    {/if}

    {#if !thread}
    <div class="gr-center et-ic"
        aria-label="Create Thread"
        data-microtip-position="top"
        data-microtip-size="fit"
        on:click={createThread}
        role="tooltip">
        <svg class="icon-3Gkjwa" aria-hidden="false" width="24" height="24" viewBox="0 0 24 24" fill="none"><path fill="currentColor" d="M5.43309 21C5.35842 21 5.30189 20.9325 5.31494 20.859L5.99991 17H2.14274C2.06819 17 2.01168 16.9327 2.02453 16.8593L2.33253 15.0993C2.34258 15.0419 2.39244 15 2.45074 15H6.34991L7.40991 9H3.55274C3.47819 9 3.42168 8.93274 3.43453 8.85931L3.74253 7.09931C3.75258 7.04189 3.80244 7 3.86074 7H7.75991L8.45234 3.09903C8.46251 3.04174 8.51231 3 8.57049 3H10.3267C10.4014 3 10.4579 3.06746 10.4449 3.14097L9.75991 7H15.7599L16.4523 3.09903C16.4625 3.04174 16.5123 3 16.5705 3H18.3267C18.4014 3 18.4579 3.06746 18.4449 3.14097L17.7599 7H21.6171C21.6916 7 21.7481 7.06725 21.7353 7.14069L21.4273 8.90069C21.4172 8.95811 21.3674 9 21.3091 9H17.4099L17.0495 11.04H15.05L15.4104 9H9.41035L8.35035 15H10.5599V17H7.99991L7.30749 20.901C7.29732 20.9583 7.24752 21 7.18934 21H5.43309Z"></path><path fill="currentColor" d="M13.4399 12.96C12.9097 12.96 12.4799 13.3898 12.4799 13.92V20.2213C12.4799 20.7515 12.9097 21.1813 13.4399 21.1813H14.3999C14.5325 21.1813 14.6399 21.2887 14.6399 21.4213V23.4597C14.6399 23.6677 14.8865 23.7773 15.0408 23.6378L17.4858 21.4289C17.6622 21.2695 17.8916 21.1813 18.1294 21.1813H22.5599C23.0901 21.1813 23.5199 20.7515 23.5199 20.2213V13.92C23.5199 13.3898 23.0901 12.96 22.5599 12.96H13.4399Z"></path></svg>
    </div>
    {/if}



    <ToolsMenu
    event={event}
    owner={owner}
    on:showMenu={showMenu}
    on:hideMenu={hideMenu}
    />


<style>
.et-ic {
    padding: 0.25rem;
    cursor: pointer;
}
.et-ic:hover {
    fill: var(--white);
}
</style>
