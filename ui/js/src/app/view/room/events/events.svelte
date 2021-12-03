<script>
import { store } from '../../../store/store.js'
import { createEventDispatcher, onMount } from 'svelte'
import { threadBiggest } from '../../../utils/icons.js'
const dispatch = createEventDispatcher()
import { useLocation } from 'svelte-navigator'
import EventItem from './event-item.svelte'


export let threadView;
export let inactive;
export let topicsView;

const location = useLocation()

let indicate = false;

$: {
    indicate = $location.pathname == page.pathname
}

$: server = page.pathname.split("/")[1]
$: roomAlias = page.pathname.split("/")[2]

$: eventAlias = $location.pathname.split("/")[4]



$: mentionedEvent = messages.filter(x => x.event_id == eventAlias)

$: mentionedEventExists = eventAlias?.length > 0 && mentionedEvent != null && mentionedEvent != undefined



$: isDirectMessage = server == 'messages'

$: account = $store.accounts.filter(account => account.user_id == $store.active_account)[0]
$: servers = account.servers

$: rooms = servers?.filter(s => s.pathname == `/${server}`)[0]?.rooms

$: dms = $store.accounts.filter(account => account.user_id == $store.active_account)[0]?.direct_messages

$: dmRoom = dms?.filter(x => x.alias == roomAlias) [0]

$: room = isDirectMessage ? dmRoom : rooms?.filter(r => r.pathname == `/${roomAlias}`)[0]

$: threadRoom = isDirectMessage ? dms?.filter(r => r?.channel_id ==
    room.thread_in_room_id)[0] : rooms?.filter(r => r?.channel_id == room.thread_in_room_id)[0]

$: threadRoomEvents = events?.filter(x => x.type == 'commune.room.thread.initial')


$: roomID = room?.room_id

export let page;
export let visible;

$: roomID = room?.room_id


$: events = $store.events[roomID]?.events
$: from = $store.events[roomID]?.end


let init = 0;

let intersecting;
let tintersecting;

let obs;
let observer;
let options;

let top;
let tobserver;
let toptions;

let tobs;
let obst;
let optt;

let fobs;
let obsf;
let optf;


onMount(() => {
    dispatch('ready', true)
    if((chatroom || thread) && obs && !topicsView) {
        setTimeout(() => {
            options = {
              root: document,
              rootMargin: '0px',
              threshold: 1.0
            }
            observer = new IntersectionObserver(callback, options);
            if(obs) {
                observer.observe(obs);
            }
        }, 100)
    }
    if((chatroom || thread) && top && !topicsView) {
        setTimeout(() => {
            toptions = {
              root: document,
              rootMargin: '0px',
              threshold: 0
            }
            tobserver = new IntersectionObserver(tcallback, toptions);
            if(top) {
                tobserver.observe(top);
            }
        }, 100)
    }
    if((topics) && tobs && topicsView) {
        setTimeout(() => {
            toptions = {
              root: document,
              rootMargin: '0px',
              threshold: 0
            }
            obst = new IntersectionObserver(tocallback, optt);
            if(tobs) {
                obst.observe(tobs);
            }
        }, 100)
    }
    if((mentionedEventExists) && fobs) {
        setTimeout(() => {
            toptions = {
              root: document,
              rootMargin: '0px',
              threshold: 0
            }
            obsf = new IntersectionObserver(focallback, optf);
            if(fobs) {
                obsf.observe(fobs);
            }
        }, 100)
    }
    window.syncEventsPosition = (id) => {
        if(id == roomID) {
            dispatch('scrollDown', true)
            /*
            if(intersecting) {
                dispatch('scrollDown', true)
            }
            */
        }
    }
    window.syncReactionPosition = (room_id, event_id) => {
        console.log("we adddddddd")
        if(room_id == roomID) {
        console.log("same room")
            if(messages[messages.length-1]?.event_id == event_id) {
            console.log("is last event")
                dispatch('scrollDown', true)
            }
        }
    }
})

let callback = (entries, observer) => {
  entries.forEach(entry => {
      if(entry?.isIntersecting) {
          intersecting = true
      } else {
          intersecting = false
      }
  });
};

let tcallback = (entries, observer) => {
  entries.forEach(entry => {
      if(entry?.isIntersecting) {
          tintersecting = true
      }
  });
};

let tbintersecting = false;

let tocallback = (entries, observer) => {
  entries.forEach(entry => {
      if(entry?.isIntersecting) {
          console.log(entry)
          tbintersecting = true
      }
  });
};

let fintersecting = false;
let focallback = (entries, observer) => {
  entries.forEach(entry => {
      if(entry?.isIntersecting) {
          console.log(entry)
          fintersecting = true
      }
  });
};

$: if(fintersecting) {
    dispatch('load-forward', {
        room_id: roomID,
        from: from,
    })
    fintersecting = false
}


export let noMore = false;


$: if(tintersecting && (events?.length >=50) && !noMore) {
    console.log("boink")
    dispatch('load-more', {
        room_id: roomID,
        from: from,
    })
    tintersecting = false
}

$: if(tbintersecting && (events?.length >=7) && !noMore) {
    console.log("boink")
    dispatch('load-more', {
        room_id: roomID,
        from: from,
    })
    tintersecting = false
}

$: msgType = `m.room.message`

$: messages = events?.filter(x => x.type == msgType &&
    !x.content?.[`m.new_content`])


$: count = messages?.length

$: if(count != init) {
    if(!indicate && dm) {
    }
    if(!indicate) {
        store.updateUnreadCount(roomID)
    }
    init = count
    if(intersecting) {
        if(!mentionedEventExists) {
            dispatch('scrollDown', true)
        }
    }
}


function scroll(e) {
    if(indicate) {
        if(!mentionedEventExists) {
            dispatch('scrollDown', true)
        }
    }
}

function forceScroll(e) {
    if(!mentionedEventExists) {
    }
}

$: chatroom = room?.room_type == 'chat'
$: topics = room?.room_type == 'topics'
$: thread = room?.room_type == 'thread'
$: dm = room?.room_type == 'dm'

export let replyEvent;

function showEvent(event) {
    let newContent = event?.content['m.new_content']
    if(newContent) {
        return false
    }
    return true
}

let threadUser;

function setThreadUser(e) {
    threadUser = e.detail
}


let eventsContainer;

function mounted(e) {
    if(!mentionedEventExists) {
    }
    ready = true

    let event = messages[e.detail]
    if(event.sender != $store.active_account) {
        store.updateReadMarker(roomID, event.event_id)
    }

    /*
    setTimeout(() => {
        console.log("scrolling down again")
        dispatch('scrollDown', true)
        ready = true
    }, 50)
    setTimeout(() => {
        dispatch('scrollDown', true)
    }, 100)
    */
}

let ready = false;

$: name = room?.name?.includes(',') ? room?.name : `@${room?.name}`


let highlightedEvent;

function highlighted(e) {
    highlightedEvent = e.detail
}


let eventEditing;

function editingEvent(e) {
    eventEditing = e.detail
}

function eventMentioned(e) {
    dispatch('scrollToEvent', e.detail)
}

</script>

<div class="room-events-container" class:no-vis={!ready && !dm && messages?.length > 0}>
<div class="room-events" class:topics={topics} bind:this={eventsContainer}>


    {#if threadRoom && thread}
        <div class="flex flex-column pa3">
            <div class="flex">
                <div class="th-ic gr-default">
                    {@html threadBiggest}
                </div>
                <div class="flex-one"></div>
            </div>
            <div class="pt3 f3 bold">
                {room?.name}
            </div>
        </div>
    {/if}

    {#if noMore || messages?.length < 50 || !messages}
        {#if chatroom}
            <div class="start pa3 mb3">
                <div class="f3">
                    <strong>Welcome to #{room?.name}</strong>
                </div>
                <div class="pt3 pb3 light brdr">
                    This is the start of the #{room?.name} {room.room_type}.
                </div>
            </div>
        {/if}
        {#if dm}
            <div class="start pa3 mb3">
                <div class="f3">
                    <strong>{name}</strong>
                </div>
                <div class="pt3 pb3 light brdr">
                    This is the start of your direct message with {name}.
                </div>
            </div>
        {/if}
        {#if thread}
            <div class="start mb3 flex flex-column ph3">
                <div class="light ">
                    Started by <span class="thread-user">{threadUser}</span>
                </div>
                <div class="pb3 light brdr">
                </div>
            </div>
        {/if}
    {/if}

    {#if thread && threadRoom && threadRoomEvents}
        {#each threadRoomEvents as item, i}
            <EventItem on:thread-user={setThreadUser} 
            isThread={true} 
            event={item.content.event} 
            room={room} />
        {/each}
    {/if}

{#if events && count > 0}
    {#if chatroom || thread || dm} 
        <div class="top mb3" bind:this={top}>
            {#if tintersecting && events?.length >= 50 && !noMore}
                <div class="gr-default">
                    <div class="loading gr-center pa1 mt3">
                        Loading more...
                    </div>
                </div>
            {/if}
        </div>
    {/if}

    {#each messages as event, i (event.origin_server_ts)}

                
        <EventItem 
        on:highlighted={highlighted}
        highlightedEvent={highlightedEvent}
        on:mounted={mounted} 
        on:event-mentioned={eventMentioned} 
        on:forceScroll={forceScroll} 
        on:editing={editingEvent} 
        eventEditing={eventEditing}
        rooms={rooms} 
        thread={thread} 
        replyEvent={replyEvent} 
        on:replying on:scroll={scroll} 
        room={room} 
        event={event} 
        inactive={inactive}
        last={i == messages?.length - 1}
        index={i}  />

    {/each}
    {#if chatroom || thread || dm}
        <div class="obs mb3" bind:this={obs}></div>
    {/if}

    {#if topics}
        <div class="tobs mb3" bind:this={tobs}></div>
    {/if}

    {#if mentionedEventExists}
        <div class="fobs mb3" bind:this={fobs}></div>
    {/if}

{/if}
</div>

{#if topics && messages?.length == 0}
    <div class="gr-default h-100 pv5">
        <div class="gr-center">
            No posts!
        </div>
    </div>
{/if}
</div>

<style>
.room-events-container {
    display: grid;
}

.topics {
    max-width: 960px;
    width: 100%;
    justify-self: center;
}

@media screen and (max-width: 1160px) {
    .topics {
        max-width: 100%;
    }
}
.loading {
    background-color: var(--background-1);
    border-radius: 4px;
}
.brdr {
    border-bottom: 1px solid hsla(0,0%,100%,0.06);
}
.thread-user {
    font-weight: bold;
    color: var(--text);
}
.th-ic {
    background-color: var(--background-10);
    width: 64px;
    height: 64px;
    border-radius: 50%;
}

.expire {
    color: var(--text);
}

.no-dis {
    display: none;
}
.no-vis {
    visibility: hidden;
    pointer-events: none!important;
}
</style>
