<script>
import { onMount, createEventDispatcher } from 'svelte'
import { store } from '../../../store/store.js'
import {reply as replyIcon, addReaction, user, edit as editIcon} from '../../../utils/icons.js'
import Reaction from '../events/reaction/reaction.svelte'

import { timeAgo, formatTime } from '../../../utils/time.js'

export let OP;
export let event;
export let slug;
export let isReply;
export let roomID;

const dispatch = createEventDispatcher()


$: htmlExists = event?.content?.formatted_body ? event.content.formatted_body : event?.content?.body ? event.content.body : ``

function formatBody(content) {
    let div = document.createElement('div')
    div.innerHTML = content
    let spans = div.querySelectorAll('span')
    if(spans?.length > 0) {
        spans.forEach(span => {
            if(span.dataset.userid == $store.active_account) {
                span.classList.add('mentioned')
            }
        })
    }
    return div.innerHTML
}
$: formatted = formatBody(htmlExists)

function reply() {
    dispatch('reply', event)
}


let reactContainer;

let reacting;

function printEvent(e) {
    e.preventDefault()
    console.log(event)
}

function react() {
    reacting = true
    let opts = {
        room_id: event?.room_id,
        event: event,
        container: reactContainer,
        topics: true,
        topicReply: isReply,
        slug: slug,
        //userReactions: userReactions,
    }
    window.toggleEmojiPicker(opts)
}


$: member = $store.allMembers[event?.sender]

$: displayNameExists = member?.display_name?.length > 0
$: avatarExists = member?.avatar_url?.length > 0
$: name = displayNameExists ? member?.display_name : member?.username

$: avatar = `${homeServer}/_matrix/media/r0/download/${member?.avatar_url?.substring(6)}`

$: when = timeAgo(event?.unsigned?.age)
$: time = formatTime(event?.unsigned?.age)

$: delivered = event?.delivered === true || event?.delivered === undefined


let eventContainer;

onMount(() => {
    if(event?.fresh) {
        eventContainer.scrollIntoView()
    }
})


$: opReactions = $store.events?.[event.room_id]?.events?.filter(x =>
    x?.type == 'm.reaction')?.filter(y =>
        y?.content?.['m.relates_to']?.event_id == event?.event_id)

$: replyReactions = $store.allReplies[event.room_id]?.filter(x =>
    x?.type == 'm.reaction')?.filter(y =>
        y?.content?.['m.relates_to']?.event_id == event?.event_id)

$: reactions = OP ? opReactions : replyReactions

$: distinctReactions = distinct(reactions)

let initReactions = 0

let lockReactionCount = false;

function distinct(r) {
    let dis = []
    r?.forEach(reaction => {
        let ind = dis.filter(x => x?.key == reaction?.content?.['m.relates_to']?.key)[0]
        if(!ind) {
            dis.push({
                key: reaction?.content?.['m.relates_to']?.key,
                count: 0,
            })
        } else {
            ind.count = ind.count + 1
        }
    })
    if(!lockReactionCount) {
        initReactions = dis?.length
        lockReactionCount = true
    }
    return dis?.sort((a, b) => (a.count < b.count) ? 1 : -1)
    /*
    let dis = []
    r.forEach(reaction => {
        let ind = dis.findIndex(x => x == reaction?.content?.['m.relates_to']?.key)
        if(ind == -1) {
            dis.push(reaction?.content?.['m.relates_to']?.key)
        }
    })
    return dis
    */
}


$: eventRoom = {
    room_id: roomID,
}


let adc;

function reactToEvent() {
    let opts = {
        room_id: event?.room_id,
        room: eventRoom,
        event: event,
        container: adc,
        inline: true,
        topics: true,
        topicReply: isReply,
        slug: slug,
    }
    window.toggleEmojiPicker(opts)
}

$: if($store.emoji?.active && $store.emoji?.event == event?.event_id) {
    reacting = true;
} else {
    reacting = false;
}


$: owner = event?.sender == $store.active_account

function edit() {
    console.log("editing")
    dispatch('edit', event)
}

$: relations = event?.unsigned?.['m.relations']
$: edited = relations?.['m.replace'] || event.edited

</script>



{#if !isReply}
<div class="f4 bold">
    {event?.content?.title}
</div>
{/if}

<div class="topic-post-item pt3 ph3 pb2 fl-co mt3"
    bind:this={eventContainer}
    class:fresh={event?.fresh}
    class:o-40={!delivered}
    on:contextmenu={printEvent}>
    <div class="flex">

        <div class="mr3">
        {#if avatarExists}
            <div class="profile-avatar ncol bg-img"
                style="background-image: url({avatar});">
            </div>
        {:else}
            <div class="profile-avatar gr-default">
                <div class="log gr-center">
                    {@html user}
                </div>
            </div>
        {/if}
        </div>
        <div class="fl-co flex-one">
            <div class="flex">
                <div class="gr-center">
                    <span class="username">{name} </span>
                </div>
                <div class="flex-one">
                </div>
                <div class="date gr-center mute">
                    {when} {time}
                </div>
            </div>
            <div class="content-body mt2">
                {@html formatted}
                {#if edited}
                    <span class="edited">(edited)</span>
                {/if}
            </div>
        </div>
    </div>

    <div class="topic-item-footer flex mt2">
        <div class="gr-center flex-one">
            {#if distinctReactions}
                <div class="event-reactions flex" 
                class:pv1={distinctReactions.length > 0}>
                    {#each distinctReactions as reaction, i (i)}
                        <Reaction 
                        key={reaction.key} 
                        event={event}
                        isTopic={true}
                        OP={OP}
                        room={eventRoom} />
                    {/each}
                    {#if distinctReactions?.length > 0}
                    <div class="r-ico-h gr-default">
                        <div class="ml1 r-ico gr-default gr-center"
                            class:r-ac={reacting}
                            bind:this={adc}
                            on:click={reactToEvent}>
                            {@html addReaction}
                        </div>
                    </div>
                    {/if}
                </div>
            {/if}
        </div>
        <div class="gr-center">
        </div>

        {#if owner}
        <div class="pointer ico gr-center mr2 o-50 hov-op"
            on:click={edit}>
            {@html editIcon}
        </div>
        {/if}

        <div class="pointer ico gr-center mr2 o-50 hov-op"
            on:click={react}
            bind:this={reactContainer}>
            {@html addReaction}
        </div>

        <div class="pointer ico gr-center flex o-50 hov-op"
        on:click={reply}>
            {@html replyIcon}
        </div>
    </div>

</div>

<style>
.topic-post-item {
    background-color: var(--room-event-hover);
    border: 1px solid transparent;
    border-radius: 7px;
    transition: 0.1s;
    display: grid;
    grid-template-rows: 1fr auto;
}

.fresh {
    animation-name: fade;
    animation-duration: 60s;
}

@keyframes fade {
  from {border: 1px solid var(--green);}
  to {border: 1px solid transparent;}
}

.content-body {
    line-height: 1.4rem;
}
.ico {
    width: 24px;
    fill: var(--text);
}
.topic-item-footer {
    transition: 0.1s;
}

.profile-avatar {
    width: 32px;
    height: 32px;
    background-color: var(--avatar);
    border-radius: 50%;
    transition: 0.1s;
}
.profile-avatar:hover {
}
.ncol {
    background-color: transparent;
}
.username {
    color: var(--white);
    cursor: pointer;
}


.date {
    font-size: 0.8rem;
}

.bold {
    font-weight: bold;
}

.username:hover {
    text-decoration: underline;
}

.r-ico {
    width: 24px;
    height: 24px;
    fill: var(--text);
    cursor: pointer;
    padding: 0.125rem;
    opacity: 0;
    transition: 0.1s;
}

.r-ico:hover {
    fill: var(--white);
}

.r-ac {
    opacity: 1;
}

.event-reactions:hover .r-ico {
    opacity: 1;
}
.edited {
     font-size: 0.8rem;
     color: var(--text-muted);
}
</style>
