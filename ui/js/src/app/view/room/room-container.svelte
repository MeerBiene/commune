<script>
import { store } from '../../store/store.js'
import { useLocation } from 'svelte-navigator'
import {onMount} from 'svelte'

import Room from './room.svelte'
import Input from './input/input.svelte'
import Header from './header/header.svelte'
import Members from './members/members.svelte'

import SplitView from './splitview/splitview.svelte'

import Thread from '../../thread/thread.svelte'

import TopicPost from './topics/topic-post.svelte'

const location = useLocation()


$: page = $location.pathname

$: server = page.split("/")[1]
$: roomAlias = page.split("/")[2]
$: isTopicEvent = page.split("/")[3] == 'topic'
$: isMessageEvent = page.split("/")[3] == 'message'
$: topicPost = page.split("/")[4]


$: hideSplitView = isTopicEvent || isMessageEvent

onMount(() => {
    if(topicPost?.length > 0 && isTopicEvent) {
        store.updateActiveRoom(``, page)
        store.updateActiveRooms(``, page)
    }
})

$: isDirectMessage = server == 'messages'


$: account = $store.accounts.filter(account => account.user_id == $store.active_account)[0]
$: servers = account?.servers

$: rooms = servers?.filter(s => s.pathname == `/${server}`)[0]?.rooms

$: dms = $store.accounts.filter(account => account.user_id == $store.active_account)[0]?.direct_messages

$: dmRoom = dms?.filter(x => x.alias == roomAlias) [0]

$: room = isDirectMessage ? dmRoom : rooms?.filter(r => r.pathname == `/${roomAlias}`)[0]

$: activeDMs = $store.accounts.filter(account => account.user_id == $store.active_account)[0]?.active_direct_messages

$: activeRooms = isDirectMessage ? activeDMs : account?.active_rooms

$: chatroom = room?.room_type == 'chat'
$: topics = room?.room_type == 'topics'
$: thread = room?.room_type == 'thread'
$: dm = room?.room_type == 'dm'

$: showMembers = membersMode()

function membersMode() {
    let mode = localStorage.getItem("members-list");
    if(mode && mode == "visible") {
        return true
    }
    return false
}

function toggleShowMembers() {
    if(showSplitView) {
        showSplitView = false
        localStorage.removeItem("split-view")
    }
    showMembers = !showMembers
    if(showMembers) {
        localStorage.setItem("members-list", "visible")
    } else {
        localStorage.removeItem("members-list")
    }
}


$: showSplitView = splitViewMode()

function splitViewMode() {
    let mode = localStorage.getItem("split-view");
    if(mode && mode == "visible") {
        return true
    }
    return false
}

function toggleSplitView() {
    if(showMembers) {
        showMembers = false
        localStorage.removeItem("members-list")
    }
    showSplitView = !showSplitView
    if(showSplitView) {
        localStorage.setItem("split-view", "visible")
        store.splitView = true
    } else {
        localStorage.removeItem("split-view")
        store.splitView = false
    }
}

$: if(store.splitView == false) {
    showSplitView = false
    localStorage.removeItem("split-view")
}


let replying = false;
let replyEvent;

function reply(e) {
    replyEvent = e.detail.event
    replying = true
}

function discardReply() {
    replyEvent = null
    replying = false
}


$: roomExists = room != null && room != undefined


$: threadActive = room?.thread?.active


$: if(threadActive && showMembers) {
    toggleShowMembers()
}

$: if(threadActive && showSplitView) {
    toggleSplitView()
}

$: hide = $location.pathname == '/' ||
            $location.pathname == '/explore' ||
            $location.pathname == '/users' ||
            $location.pathname.substring(1,2) == `@`



</script>

<div class="room-root"
    class:no-dis={hide}
    class:thread-active={threadActive}>
    <div class="room-container" 
        class:show-sidebar={showMembers || (showSplitView && !hideSplitView)}
        class:sd-1={showMembers}
        class:sd-2={showSplitView && !hideSplitView}
        class:r-br={threadActive}>

    <Header 
        page={page} 
        roomAlias={roomAlias} 
        server={server}
        servers={servers}
        rooms={rooms}
        showMembers={showMembers}
        showSplitView={showSplitView}
        on:toggleShowMembers={toggleShowMembers}
        on:toggleSplitView={toggleSplitView}
        room={room} />

        <div class="room-division relative" 
            class:no-dis={topicPost && isTopicEvent}
            class:sm={showMembers || (showSplitView && !hideSplitView)}
            class:chatroom={chatroom || thread || dm}>

            <div class="room-items">
                {#if activeRooms}
                {#each activeRooms as p, i (p.pathname)}
                    <Room 
                        page={p} 
                        on:replying={reply}
                        replyEvent={replyEvent}
                        showMembers={showMembers}/>
                {/each}
                {/if}
            </div>

            {#if chatroom || thread || topics || dm}
            <Input 
                roomAlias={roomAlias} 
                server={server}
                replying={replying}
                replyEvent={replyEvent}
                on:discardReply={discardReply} />
            {/if}
        </div>

        {#if topicPost && isTopicEvent}
            <TopicPost />
        {/if}

    {#if showMembers}
        <Members room={room} />
    {/if}

    {#if showSplitView && !hideSplitView}
        <SplitView room={room} 
        on:killSplitView={toggleSplitView}/>
    {/if}


</div>

{#if threadActive}
    <div></div>
    <Thread room={room}/>
{/if}

</div>

<style>
.room-root {
    display: grid;
    width: 100%;
    height: 100%;
    grid-template-columns: 100%;
    grid-template-rows: 100%;
}
.thread-active {
    grid-template-columns: 1fr 8px 400px;
}

@media screen and (max-width: 1150px) {
    .thread-active {
        grid-template-columns: 100%;
    }
}

.room-container {
    display: grid;
    width: 100%;
    height: 100%;
    grid-template-columns: 100%;
    grid-template-rows: 48px 1fr;
    position: relative;
    background-color: var(--background-3);
}

.r-br {
    border-radius: 0 8px 8px 0;
}

.show-sidebar {
    grid-template-areas:
      "header header"
      "division sidebar"
}

.sd-1 {
    grid-template-columns: 1fr 232px;
}

.sd-2 {
    grid-template-columns: 1fr 232px;
}

.room-items {
    display: grid;
    width: 100%;
    height: 100%;
    overflow: hidden;
}
.room-division {
    display: grid;
    width: 100%;
    height: 100%;
    grid-template-columns: 100%;
    grid-template-rows: 100%;
    overflow: hidden;
}
.sm {
    grid-column: 1/2;
    grid-row: 2;
}
.chatroom {
    grid-template-rows: [chat] 1fr [input] minmax(68px, auto);
}
.no-dis {
     display: none;
}
</style>
