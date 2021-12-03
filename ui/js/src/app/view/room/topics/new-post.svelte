<script>
import {createEventDispatcher, tick} from 'svelte'
import { store } from '../../../store/store.js'
import { makeid } from '../../../utils/utils.js'
import {reply as replyIcon, down, edit as editIcon} from '../../../utils/icons.js'
const dispatch = createEventDispatcher()

export let room;
export let reply;
export let replyingTo;
export let replyingToTopic
export let editingEvent;
export let opEvent;

$: editingOP = editingEvent?.event_id == opEvent?.event_id

$: if(opEvent && editingOP && titleInput) {
    titleInput.value = editingEvent?.content?.title
}

$: placeholder = `Post Title`

$: roomID = room?.room_id

import Editor from '../../../components/editor/editor.svelte'

let active = false;

export function activate() {
    active = true
    if(!reply) {
        focusTitle()
    }
}

$: if(reply && editor) {
    editor.focusEditor()
}

function kill() {
    title = ''
    if(titleInput) {
        titleInput.value = ''
    }
    content = null
    active = false
    dispatch('killed')
}


let title;
let titleInput;


async function focusTitle() {
    await tick();
    titleInput.focus()
}


$: buttonText = editing ? `Save Edit` : reply ? `Reply` : `Create Topic`

let content;
function updateContent(e) {
    content = e.detail
}

$: sender = $store.accounts.filter(x => x.user_id == $store.active_account)[0].user_id

$: account = $store.accounts.filter(account => account.user_id == $store.active_account)[0]

$: matrix = account?.matrix

const slugify = text =>
  text
    .toString()
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '')
    .toLowerCase()
    .trim()
    .replace(/\s+/g, '-')
    .replace(/[^\w-]+/g, '')
    .replace(/--+/g, '-')

let randomNumber = () => {
    return Math.random().toString().slice(2,6);
}

let getSlug = () => `${slugify(title)}-${Date.now()}`

function create() {
    if(titleInput?.value?.length == 0 && !reply && !editing) {
        alert("Post title cannot be empty.")
        focusTitle()
        return
    }
    if(titleInput?.value?.length == 0 && editingOP) {
        alert("Post title cannot be empty.")
        focusTitle()
        return
    }

    if(editingOP) {
        if(!content) {
            content = {
                plain_text: editingEvent?.content?.body,
                html: editingEvent?.content?.formatted_body,
                length: editingEvent?.content?.body?.length
            }
        }
    } else {
        if(content?.length == 0 || content?.plain_text == `` || !content) {
            alert("Post body cannot be empty.")
            editor.focusEditor()
            return
        }
    }


    let tempid = makeid(16)

    let {plain_text, html, length} = content


    if(reply && !editing) {
        dispatch('replied', content)
        kill()
        return
    }

    if(editing) {
        dispatch('edited', {
            title: title,
            content: content,
        })
        kill()
        return
    }

    let con = {
        "body": plain_text,
        "formatted_body": html,
        "title": title,
        "msgtype": "m.text",
        //"slug": getSlug(),
        //"name": getSlug(),
    };

    if(!reply) {
        con.topic = true
    } else {
        con.topic_reply = true
    }


    store.createTopicRoom().then(res => {
        if(res?.created && res?.room_id) {
            con.topic_room_id = res?.room_id

            let ts = new Date()
            store.addEventToRoom(roomID, {
                "type": "m.room.message",
                "room_id": roomID,
                "sender": sender,
                "content": con,
                "origin_server_ts": ts,
                "unsigned": {
                    "age": 0
                },
                "age": 0,
                "transaction_id": tempid,
                "delivered": false,
                "fresh": true,
            }, true)
            kill()
        }
    })

}

let editor;
let editorFocused;
function focusState(e) {
    editorFocused = e.detail
}

$: replyToMember = $store.allMembers[replyingTo?.sender]

$: displayNameExists = replyToMember?.display_name?.length > 0
$: avatarExists = replyToMember?.avatar_url?.length > 0
$: name = displayNameExists ? replyToMember?.display_name : replyToMember?.username


let collapsed = false;

function toggleCollapse() {
    collapsed = !collapsed
    if(!reply && !collapsed) {
        focusTitle()
    } else if(!collapsed) {
        editor.focusEditor()
    }
}

$: editing = editingEvent != null && editingEvent != undefined
</script>

{#if active}
<div class="topics-post-container">
    <div class="new-topics-post flex flex-column">
    {#if replyingTo}
        <div class="topic-editor-header ph3 pv2">
        <div class="replying-to">
                <div class="flex">
                    <div class="rico gr-center">
                        {@html replyIcon}
                    </div>
                    <div class="gr-center ml2 clmp-1 flex-one">
                        {#if replyingToTopic}
                            <span class="username">{replyingTo?.content?.title}</span>
                        {:else}
                            <span class="username">@{name}</span>
                        {/if}
                    </div>
                    <div class="down-ic gr-center"
                        class:dv={collapsed}
                        on:click={toggleCollapse}>
                        {@html down}
                    </div>
                </div>
        </div>
        </div>
    {/if}
    {#if editing}
        <div class="topic-editor-header ph3 pv2">
        <div class="replying-to">
                <div class="flex">
                    <div class="rico gr-center">
                        {@html editIcon}
                    </div>
                    <div class="gr-center ml2 clmp-1 flex-one">
                        Editing
                    </div>
                    <div class="down-ic gr-center"
                        class:dv={collapsed}
                        on:click={toggleCollapse}>
                        {@html down}
                    </div>
                </div>
        </div>
        </div>
    {/if}

            <div class="flex flex-column">
                {#if (!reply && !editing) || editingOP}
                    <div class="input-header flex">
                        <div class="flex-one gr-center">
                            {#if !collapsed}
                            <input 
                            bind:this={titleInput}
                            bind:value={title}
                            placeholder={placeholder}/>
                            {:else}
                                <div class="gr-center ml3">
                                    Topic Draft
                                </div>
                            {/if}
                        </div>
                        <div class="down-ic gr-center ph3"
                            class:dv={collapsed}
                            on:click={toggleCollapse}>
                            {@html down}
                        </div>
                    </div>
                {/if}

                <div class="new-post-container gr-default"
                class:no-dis={collapsed}
                class:ef={editorFocused}>
                    <Editor 
                    on:focused={focusState}
                    bind:this={editor}
                    placeholder={"Your post content goes here."}
                    on:sync={updateContent}
                    editing={editing}
                    initial={editingEvent ? editingEvent?.content?.formatted_body : null}
                    room={room}
                    topics={true}/>
                </div>
                <div class="topic-footer flex pa3"
                    class:no-dis={collapsed}>
                    <div class="">
                    </div>
                    <div class="flex-one">
                    </div>
                    <div class="gr-center">
                        <button class="" on:click={create}>{buttonText}</button>
                    </div>
                    <div class="gr-center ml3">
                        <span class="pointer" on:click={kill}>Cancel</span>
                    </div>
                </div>

            </div>

    </div>
</div>
{/if}


<style>

.topics-post-container {
    display: grid;
}

.topics-reply-container {
    display: grid;
}

.new-topics-post {
    max-width: 1060px;
    width: 100%;
    justify-self: center;
    background-color: var(--chat-input);
    box-shadow: 0 0px 20px rgba(0,0,0,.18);
    border-top: 7px solid var(--background-1);
}

@media screen and (max-width: 1280px) {
    .new-topics-post {
        max-width: 100%;
    }
}
:root {
    --chat-input: #40444b;
}

input {
    background-color: var(--chat-input);
    width: 100%;
    padding: 1rem;
    border: 1px solid transparent;
    transition: 0.1s;
    font-size: 1.2rem;
    font-weight: bold;
}

input:focus {
}

.new-post-container {
    background-color: var(--chat-input);
    width: 100%;
    height: 100%;
    min-height: 200px;
    transition: 0.1s;
    border: 1px solid transparent;
}


.ef {
}

.topic-footer {
    border-top: 1px solid var(--background-3);
}

.username {
    color: var(--white);
    cursor: pointer;
}
.username:hover {
    text-decoration: underline;
}
.rico {
    width: 24px;
    fill: var(--text);
    -webkit-transform: scaleX(-1);
  transform: scaleX(-1);
}
.down-ic {
    fill: var(--text);
    cursor: pointer;
}
.down-ic:hover {
    fill: var(--white);
}

.dv {
  transform: scaleY(-1);
}

.no-dis {
    display: none;
}
    .input-header {
        min-height: 56px;
    }
</style>
