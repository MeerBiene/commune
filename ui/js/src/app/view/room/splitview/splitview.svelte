<script>
import { onMount, createEventDispatcher } from 'svelte'
import { store } from '../../../store/store.js'
import { close, hourglass } from '../../../utils/icons.js'
import { useLocation } from 'svelte-navigator'
import EventItem from '../events/event-item.svelte'

const location = useLocation()

$: isTopicEvent = $location.pathname.split("/")[3] == `topic`
$: isMessageEvent = $location.pathname.split("/")[3] == `message`

$: noDis = isTopicEvent || isMessageEvent

const dispatch = createEventDispatcher();

export let room;

$: account = $store.accounts.filter(account => account.user_id == $store.active_account)[0]



$: chatroom = room?.room_type == 'chat'
$: topics = room?.room_type == 'topics'
$: thread = room?.room_type == 'thread'
$: dm = room?.room_type == 'dm'

$: streamText = chatroom ? `topic` : `chat`

$: streamType = chatroom ? `topics` : `chat`

function killSplitView() {
    dispatch('killSplitView')
}

$: events = $store.events[room?.streams[streamType]]?.events?.filter(x => x.type
== `m.room.message`)

let ready = false;

onMount(() => {
    ready = true
})

</script>

<div class="split-c" class:no-dis={noDis}>
    <div class="room-view">
        <div class="split-content scrl">
            <div class="hed gr-default">
                <div class="ml3 flex gr-center w-100">
                    <div class="header-t">
                        {streamText} STREAM
                    </div>
                    <div class="flex-one"></div>
                    <div class="kill-sv mr3 pointer gr-default"
                        on:click={killSplitView}>
                        {@html close}
                    </div>
                </div>
            </div>
            <div class="c-con fl-co pv2 ">
                {#if ready}
                    {#if events?.length > 0}
                        <div class="">
                            {#each events as event, i (event?.origin_server_ts)}
                                <EventItem
                                splitView={true}
                                event={event}
                                index={i}
                                room={room}
                                />
                            {/each}

                        </div>
                    {:else if events?.length == 0}
                        <div class="h-100 gr-default">
                            <div class="gr-center">
                                Nothing in this stream.
                            </div>
                        </div>
                    {/if}
                {:else}
                    <div class="h-100 gr-default">
                        <div class="hourglass gr-center">
                            {@html hourglass}
                        </div>
                    </div>
                {/if}
            </div>
        </div>
    </div>
</div>

<style>

.split-c {
    overflow: hidden;
    display: grid;
    grid-column: 2/3;
    grid-row: 2;
}

.kill-sv {
    width: 20px;
    fill: var(--text);
    opacity: 0;
}

.kill-sv:hover {
    fill: var(--white);
}

@media screen and (max-width: 1150px) {
    .split-c {
        width: 432px;
        position: fixed;
        top:0;
        right: 0;
        bottom: 0;
        -webkit-filter: drop-shadow(0 8px 40px rgba(0,0,0,.32));
        filter: drop-shadow(0 8px 40px hsla(0,0%,0%,.32));
    }
    .hed {
        border-bottom: 1px solid var(--background-1);
        background-color: var(--background-3);
    }
    .kill-sv {
        opacity: 1;
    }
}



.room-view {
    border-left: 1px solid var(--background-1);
    background-color: var(--background-3);
    overflow: hidden;
    display: grid;
}

.hed {
    position: sticky;
    top: 0;
    height: 48px;
}

.header-t {
    background-color: var(--background-1);
    border-radius: 500px;
    padding: 0.25rem 0.5rem;
    border: 1px solid transparent;
    transition: 0.1s;
    text-transform: uppercase;
    font-size: 0.72rem;
    letter-spacing: 1px;
    font-weight: bold;
    color: var(--text-light);
}

.split-content {
    display: grid;
    overflow: hidden auto;
    grid-template-rows: auto 1fr;
}

.c-con {
}

.scrl  {
    overflow-y: auto;
    scrollbar-width: thin;
    scrollbar-color: transparent transparent;
    scroll-behavior: smooth;
}

.scrl:hover {
    scrollbar-color: var(--background-1) transparent;
}

.scrl::-webkit-scrollbar {
  width: 6px;
    border-radius: 1px;
}
.scrl::-webkit-scrollbar-track {
    background: transparent;
}
.scrl::-webkit-scrollbar-thumb {
    background-color: transparent;
}
.scrl:hover::-webkit-scrollbar-thumb {
  background-color: var(--background-1);
}

.hourglass {
    fill: var(--text-muted);
    animation-duration: 1s;
    animation-name: rotate;
    animation-iteration-count: infinite;

}
@keyframes rotate {
    25% {
        transform: rotate(180deg);
    }
    75% {
        transform: rotate(180deg);
    }

    to {
        transform: rotate(360deg);
    }
}
.no-dis {
    display: none;
}
</style>
