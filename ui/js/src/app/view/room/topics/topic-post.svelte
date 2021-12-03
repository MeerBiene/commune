<script>
import { useLocation } from 'svelte-navigator'
import {onMount} from 'svelte'
import { store } from '../../../store/store.js'
import { makeid } from '../../../utils/utils.js'
import {formatTimeAgo} from '../../../utils/time.js'

import LoadingEvents from './placeholder/loading-events.svelte'

import TopicPostItem from './topic-post-item.svelte'
import NewPost from './new-post.svelte'

const location = useLocation()
$: page = $location.pathname

$: server = page.split("/")[1]
$: roomAlias = page.split("/")[2]
$: isTopicEvent = page.split("/")[3] == 'topic'
$: topicID = page.split("/")[4]


$: account = $store.accounts.filter(account => account.user_id == $store.active_account)[0]

$: rooms = account?.servers?.filter(s => s.pathname == `/${server}`)[0]?.rooms
$: room = rooms?.filter(r => r.pathname == `/${roomAlias}`)[0]
$: roomID = room?.streams['topics']

async function fetchTopicPost() {
    let endpoint = `${homeServer}/_matrix/client/r0/rooms/${roomID}/event/${topicID}`

    let resp = await fetch(endpoint, {
        headers: {
            'Authorization': `Bearer ${account.matrix_access_token}`
        }
    })
    const ret = await resp.json()
    return Promise.resolve(ret)
}

let event;

$: getReplies = $store.allReplies[event?.content?.topic_room_id]?.filter(x =>
    (x?.type == `m.room.message`) &&
    !x?.content?.[`m.new_content`])

$: replies = getReplies?.reverse()

let root;


$: if(event?.content?.topic_room_id) {
    store.fetchRoomEvents({
        room_id: event?.content?.topic_room_id,
    }).then(res => {
        console.log(res)
        if(res?.chunk) {
            store.updateTopicReplies(event.content.topic_room_id, res?.chunk)
            ready = true
            joinRepliesRoom()
        }
    })
}

function joinRepliesRoom() {
    matrix.joinRoom(event.content.topic_room_id,{}, (err, res) => {
        console.log(res)
    })
}

let ready = false;

function replied(e) {
    let {plain_text, html, length} = e.detail

    let msgType = `m.room.message`
    let content = {
        "body": plain_text,
        "formatted_body": html,
        "msgtype": msgType,
        "topic_reply": true,
    };

    let tempid = makeid(16)

    let eventType = `m.room.message`


    matrix.sendEvent(event.content.topic_room_id, eventType, content, tempid, (err, res) => {
        if(res) {
            //increment()
            store.updateTopicReplies(event?.content?.topic_room_id,[{
                event_id: res?.event_id,
                age: 0,
                unsigned: {
                    age: 0,
                },
                origin_server_ts: Date.now(),
                type: `m.room.message`,
                delivered: true,
                fresh: true,
                sender: $store.active_account,
                content: content,
            }])
            store.updateTopicRepliesCount(server, roomAlias, res?.event_id)
        }
    });
}

function increment() {
    let rel_type = 'commune.room.topics.post.reply'
    let tempid = makeid(32)
    let content = {
        'm.relates_to': {
            event_id: event.event_id,
            key: tempid,
            slug: slug,
            rel_type: "m.annotation",
        }
    }

    matrix.sendEvent(roomID, rel_type, content, "", (err, res) => {
        console.log(res)
    });
}

function edited(e) {

    let title = e.detail.title
    let {plain_text, html, length} = e.detail.content

    if(length == 0) {
        return
    }

    if(plain_text == event.content.body && title == event.content.title) {
        return
    }

    let isOP = editingEvent.event_id == event.event_id

    let content = {
        "body": ` * ${plain_text}`,
        "msgtype": "m.text",
        "m.new_content": {
            "body": plain_text,
            "formatted_body": html,
            "msgtype": "m.text",
        },
        "m.relates_to": {
            "event_id": editingEvent?.event_id,
            "rel_type": "m.replace",
        }
    };

    if(isOP) {
        content['m.new_content']['title'] = title
        content['m.new_content']['topic_room_id'] = event.content.topic_room_id
        content['m.new_content']['topic'] = true
    } else {
        content['m.new_content']['topic_reply'] = true
    }

    let txnID = makeid(32)

    let eventType = `m.room.message`

    let room_id;

    if(isOP) {
        room_id = event.room_id
    } else {
        room_id = event.content.topic_room_id
    }

    console.log(room_id)

    matrix.sendEvent(room_id, eventType, content, txnID, (err, res) => {
        if(res) {
            console.log(res)
            editingEvent = null

            if(isOP) {
                event.content.body = plain_text
                event.content.title = e.detail.title
                event.content.formatted_body = html
                event.edited = true
            }
        }
    });
}

$: matrix = account?.matrix


let replyEditor;

let replyingTo;

let replyingToTopic = false;

function replyToPost(e) {
    replyingTo = e.detail
    console.log(e.detail)
    if(e.detail.event_id == event.event_id) {
        replyingToTopic = true
    } else {
        replyingToTopic = false
    }
    replyEditor.activate()
}

let editingEvent;

function editPost(e) {
    editingEvent = e.detail
    replyEditor.activate()
}

function newReply() {
    replyingTo = event
    replyingToTopic = true
    replyEditor.activate()
}

function killed() {
    replyingTo = null
    replyingToTopic = false;
    editingEvent = null
}

$: slug = `${event?.content?.slug}`

$: firstEventTime = formatTimeAgo(event?.unsigned?.age)
$: lastEventTime = formatTimeAgo(replies?.[replies?.length - 1]?.unsigned?.age)

let tickTop = 0

onMount(() => {
    fetchTopicPost().then(res => {
        console.log(res)
        if(res) {
            event = res
        }
    })
    root.addEventListener('scroll', scrollTicker)
    store.toggleRoomStreamTo(room?.server_id, room?.room_id, 'topics')
})

let lockScroll = false;

function scrollTicker(e) {
    if(lockScroll) {
        return
    }
    let per = ((root.scrollTop) / (root.scrollHeight - root.clientHeight))
    if(!lockScroll) {
        tickTop = per * 250
    }

    let div  = 250 / (replies?.length + 1)

    let newPosition = Math.round(tickTop/div)
    if(newPosition == 0) {
        newPosition = 1
    }
    currentPosition = newPosition
}


function seekTimeline(e) {
    if(lockScroll) {
        return
    }
    let per = e.offsetY / 300
    let div  = 250 / (replies?.length + 1)
    if(e.offsetY < div) {
        root.scrollTop = 0
        return
    }
    if(e.offsetY + 50 >= 300) {
        root.scrollTop = 10000
        return
    }
    root.scrollTop = per * (root.scrollHeight - root.clientHeight)
}


let tick;


let pos1;
let pos2;

let moving = false;

function dragTick(e) {
    moving = true;
    pos1 = pos2 - e.clientY;
    pos2 = e.clientY;
    if(tickTop >= 0 && tickTop <= 250) {
        let x = tickTop - pos1
        if(x < 0) {
            tickTop = 0
            return
        }
        if(x > 250) {
            tickTop = 250
            return
        }
        tickTop = (tickTop - pos1)

        let div  = 250 / (replies?.length + 1)

        let newPosition = Math.round(tickTop/div)
        if(newPosition == 0) {
            newPosition = 1
        }
        currentPosition = newPosition
    }
}

function startDragTick(e) {
    lockScroll = true
    e.preventDefault()
    pos2 = e.clientY;
    root.onmouseup = closeDragElement;
    root.onmousemove = dragTick;
}

let currentPosition = 1

function closeDragElement() {
    root.onmouseup = null;
    root.onmousemove = null;

    let per = tickTop / 250
    root.scrollTop = per * (root.scrollHeight - root.clientHeight)

    moving = false;

    setTimeout(() => {
        lockScroll = false
    },10)
}


$: overflowing = root?.scrollHeight > root?.clientHeight 

</script>

<div class="topic-post-root scrl" bind:this={root}>

    <div class="topic-post-container pa3">
        {#if !ready}
            <LoadingEvents />
        {/if}

        <div class="topic-post fl-co" class:gr-center={!ready}>
        {#if ready}
            <div class="topic-posts-container">
                <div class="mt4">
                    <TopicPostItem 
                    OP={true}
                    on:reply={replyToPost} 
                    on:edit={editPost} 
                    slug={slug}
                    roomID={event?.room_id}
                    event={event} />
                </div>
                {#if replies?.length > 0}
                    {#each replies as item, i (item.origin_server_ts)}
                        <TopicPostItem 
                        last={i == replies.length - 1}
                        on:reply={replyToPost} 
                        on:edit={editPost} 
                        slug={slug}
                        roomID={item?.room_id}
                        event={item} 
                        isReply={true}/>
                    {/each}
                {/if}
            </div>
        {/if}

        {#if ready}
            <div class="flex mt4">
                <div class="gr-center">
                    <button class="" on:click={newReply}>Reply</button>
                </div>
            </div>
        {/if}

        </div>


        {#if ready && replies?.length > 0 && root?.scrollHeight > root?.clientHeight}
            <div class="topic-timeline ml5 mt5 no-select">
                <div class="timeline-container fl-co">
                    <div class="">
                        {firstEventTime}
                    </div>
                    <div class="timeline ml1 mv2" 
                        on:click={seekTimeline}>
                        <div class="timeline-tick gr-default"
                            class:moving={moving}
                            bind:this={tick}
                            on:mousedown={startDragTick}
                            style="--top:{tickTop};">
                            <div class="gr-start-center ml3">
                                {currentPosition} / {replies?.length + 1}
                            </div>
                        </div>
                    </div>
                    <div class="">
                        {lastEventTime}
                    </div>
                </div>
            </div>
        {/if}




    </div>

</div>


<div class="topic-editor mt4">
    <NewPost 
    opEvent={event}
    bind:this={replyEditor} 
    room={room} 
    reply={true} 
    replyingTo={replyingTo}
    editingEvent={editingEvent}
    replyingToTopic={replyingToTopic}
    on:killed={killed}
    on:edited={edited}
    on:replied={replied}/>
</div>

<style>
.topic-post-root {
    display: grid;
    overflow: hidden scroll;
    position: relative;
}
.topic-editor {
    position: absolute;
    bottom: 0;
    right: 1rem;
    left: 1rem;
}

.topic-post-container {
    max-width: 1060px;
    width: 100%;
    justify-self: center;
    display: grid;
    grid-template-columns: 1fr auto;
}

.topic-post {
}

@media screen and (max-width: 1280px) {
    .topic-post {
        max-width: 100%;
    }
}

.no-dis {
    display: none;
}
.timeline-container {
    position: sticky;
    top: 5rem;
}
.timeline {
    border-left: 1px solid var(--background-10);
    height: 300px;
    width: 120px;
    cursor: pointer;
    position: relative;
}
.timeline-tick {
    position: absolute;
    top: var(--top);
    border-left: 6px solid var(--background-4);
    width: 100%;
    right: 3px;
    height: 50px;
    cursor: ns-resize;
    transition: background-color 2s;
}
.moving {
    border-left: 6px solid var(--green);
    transition: background-color 2s;
}
</style>
