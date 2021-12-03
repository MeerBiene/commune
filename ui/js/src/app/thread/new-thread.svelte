<script>
import { onMount, onDestroy, tick } from 'svelte'
import { useLocation } from 'svelte-navigator'
import { thread, threadBig } from '../utils/icons'
import { store } from '../store/store.js'
import { closeBig  } from '../utils/icons.js'

import { slide } from 'svelte/transition'

import EventItem from '../view/room/events/event-item.svelte'
import Editor from '../components/editor/editor.svelte'

export let room;

$: events = room?.thread?.events
$: isNew = room?.thread?.new
$: dm = room?.thread?.dm

function kill() {
    editorContent = null
    name = null
    nameInput.value = null
    newEvent = null
    store.discardNewThread(room, dm)
}

const location = useLocation()

$: server = $location.pathname.split("/")[1]

let name;
let nameInput;

async function focusInput() {
    await tick();
    nameInput.focus()
}

onMount(() => {
    focusInput()
    window.focusNewThreadInput = () => {
        focusInput()
    }
})

onDestroy(() => {
    window.focusNewThreadInput = null
})

async function createServer(ev) {
    let endpoint = `/server/create`
    let data = {
        title: name,
        room: true,
        server_id: room.server_id,
        room_id: room.channel_id,
        thread: true,
        thread_events: ev,
        type: "thread",
        expire_thread: expiryPeriod,
    };
    if(dm) {
        data.server_id = room.room_id
    }

    let resp = await fetch(endpoint, {
        method: 'POST', // or 'PUT'
        body: JSON.stringify(data),
        headers:{
            'Authorization': account.access_token,
            'Content-Type': 'application/json'
        }
    })
    const ret = await resp.json()
    return Promise.resolve(ret)
}


let creating = false;

function create() {
    if(nameInput.value.length == 0) {
        focusInput()
        return
    }
    creating = true
    if(isNew) {
        createFreshThread()
        return
    }

    createServer(events).then((res) => {
        console.log(res)
        if(res?.created) {
            let newRoom = {
                alias: res.room.alias,
                room_id: res.room.room_id,
                channel_id: res.room.room_id,
                name: name,
                joined: true,
                pathname: `/${res.room.alias}`,
                server_id: room.server_id,
                room_type: "thread",
                thread_events: events,
                thread_in_room_id: room.channel_id,
                expire_thread: expiryPeriod,
                sender: $store.active_account,
                child: true,
                streams: {
                    "chat": res?.room?.room_id,
                }
            }
            if(dm) {
                store.addRoom(room.alias, newRoom, false, true)
            } else {
                store.addRoom(`/${server}`, newRoom, false, false)
            }
            kill()
            setTimeout(() => {
                window.syncEventsPosition(room?.room_id)
            }, 1000)
        }
    }).then(() => {
        creating = false
    })
}

$: account = $store.accounts.filter(account => account.user_id == $store.active_account)[0]
$: matrix = account?.matrix

let newEvent;

function createFreshThread() {
    let content = {
        'body': '',
        'thread_title': name, 
        'new_thread': true,
        "msgtype": "commune.room.thread.new",
    };
    matrix.sendEvent(room.room_id, "m.room.message", content, "", (err, res) => {
        if(res) {
            console.log(res?.event_id)
            newEvent = res?.event_id
            let evv = {
                event_id: res?.event_id,
                sender: $store.active_account,
                age: 0,
                unsigned: {
                    age: 0,
                }
            }
            createServer(evv).then((res) => {
                console.log(res)
                if(res?.created) {
                    let newRoom = {
                        room_id: res.room.room_id,
                        name: name,
                        pathname: `/${res.room.alias}`,
                        server_id: room.server_id,
                        owner: true,
                        room_type: "thread",
                        events: [],
                        pinged: false,
                        thread_events: evv,
                        thread_in_room_id: room.room_id,
                        expire_thread: expiryPeriod,
                        sender: $store.active_account,
                        created_on: 0,
                    }
                    store.addRoom(`/${server}`, newRoom, true)

                    let con = {
                        'body': editorContent?.plain_text,
                        'formatted': editorContent?.html,
                        "msgtype": "m.room.message",
                    };
                    matrix.sendEvent(res?.room?.room_id, "m.room.message", con, "", (err, res) => {
                        if(res?.event_id) {
                            let ev = $store.temp_events.filter(x => x.event_id =newEvent)[0]
                            store.addEventToRoom(room?.room_id, ev)
                            kill()
                            creating = false
                            store.openThread(room, event, newRoom)
                        }
                    });

                }
            }).then(() => {
            })
        }
    });

}

let editorContent;
function syncContent(e) {
    editorContent = e.detail
}

function onKeyPress(e) {
    if(e.key == 'Enter') {
        create()
    }
}

let periods = [
    {
        text: '1 Hour',
        value: '1 Hour',
    },
    {
        text: '24 Hours',
        value: '24 Hours',
        default: true,
    },
    {
        text: '3 Days',
        value: '3 Days',
    },
    {
        text: '1 Week',
        value: '1 Week',
    },
]

let expiryPeriod = '24 Hours';

function periodSelected(e) {
    console.log(e.detail)
    expiryPeriod = e.detail
}

let editor;

let editorFocused;

function focusState(e) {
    editorFocused = e.detail
}
</script>

<div class="header flex">
    <div class="gr-center ph3">
        {@html thread}
    </div>
    <div class="gr-center flex-one">
        <strong>New Thread</strong>
    </div>
    <div class="gr-center ph3 pointer icon" 
        on:click={kill}
        aria-label="Close"
        data-microtip-position="bottom"
        data-microtip-size="fit"
        role="tooltip">
        {@html closeBig}
    </div>
</div>

<div></div>

<div class="hidov flex flex-column">

    <div class="flex flex-column ph3 pt3">
        <div class="th-ic gr-default">
            {@html threadBig}
        </div>
    </div>

    <div class="flex flex-column ph3">
        <div class="fr-l pt3">
            thread name
        </div>
        <div class="pt2">
            <input 
            bind:this={nameInput}
            bind:value={name}
            maxlength="100"
            on:keypress={onKeyPress}
            placeholder='New thread'/>
        </div>
    </div>

    {#if isNew}
    <div class="flex flex-column ph3">
        <div class="fr-l pt3">
            starter message
        </div>
        <div class="mt2 editor" class:ef={editorFocused}>
            <Editor 
            initial={room?.thread?.content?.html}
            focus={true} 
            show={true}
            room={room}
            on:sync={syncContent}
            on:focused={focusState}
            bind:this={editor}/>
        </div>
    </div>
    {/if}




    <div class="ph3 pt3 ">
    <div class="sep"></div>
    </div>

    {#if events}
        <div class="event-items-container">
            <div class="event-items scrl-s pt3">
            {#each events as event (event?.origin_server_ts)}
                <div class="event-item-container">
                    <EventItem preview={true} room={room} event={event} />
                </div>
            {/each}
            </div>
        </div>
    {/if}

    <div class="pa3 gr-center w-100">
        <button class="ph4 pv3" 
        disabled={creating}
        on:click={create}>Create Thread</button>
    </div>
</div>

<style>
.header {
    border-bottom: 1px solid var(--background-1);
}

button {
    width: 100%;
    background-color: var(--green);
    color: var(--white);
}

.editor {
    font-size: 1.2rem;
    background-color: var(--background-2);
    width: 100%;
    padding: 0 1rem;
    transition: 0.1s;
    border: 1px solid var(--background-1);
    min-height: 120px;
    border-radius: 5px;
    display: grid;
}

.ef {
    border: 1px solid var(--green);
}

input {
    font-size: 1.2rem;
    background-color: var(--background-2);
    width: 100%;
    padding: 0.5rem 1rem;
    transition: 0.1s;
    border: 1px solid var(--background-1);
}

input:focus {
    border: 1px solid var(--green);
}

button {
    padding: 0.5rem 0.5rem;
    font-size: 1.1rem;
}
.fr-l {
    font-size: 0.74rem;
    text-transform: uppercase;
    letter-spacing: 1px;
    color: var(--text-light);
}

.th-ic {
    background-color: var(--background-10);
    width: 64px;
    height: 64px;
    border-radius: 50%;
    fill: var(--text);
}

.th-ic svg{
    fill: var(--text);
}

.sep {
    height: 1px;
    width: 100%;
    background-color: var(--background-10);
}
.event-items-container {
    display: grid;
    overflow: hidden;
}
.event-items {
    overflow: hidden auto;
    height: 100%;
}
</style>
