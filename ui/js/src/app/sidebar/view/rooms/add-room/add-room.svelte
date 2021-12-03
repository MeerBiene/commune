<script>
import { fade, fly } from 'svelte/transition'
import { onMount, tick, createEventDispatcher } from 'svelte'
import { store } from '../../../../store/store.js'
import { hash, topics, lock } from '../../../../utils/icons.js'

import Toggle from '../../../../components/ui/toggle/toggle.svelte'

const dispatch = createEventDispatcher()

export let server;
export let serverID;


let active = false;

onMount(() => {
    active = true
    focus()
})

async function focus() {
    await tick();
    nameInput.focus()
}

function kill() {
    dispatch('kill', true)
}

let name;
let nameInput;

$: account = $store.accounts.filter(account => account.user_id == $store.active_account)[0]

async function createServer() {
    let endpoint = `/server/create`
    let data = {
        title: name.toLowerCase(),
        room: true,
        server_id: serverID,
        type: selectedStream,
    };
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

let error;

let creating = false;

function create() {
    if(nameInput.value.length < 4) {
        alert("That name is too short.")
        nameInput.focus()
        return
    }
    creating = true
    createServer().then((res) => {
        console.log(res)
        if(res?.created) {
            let newRoom = {
                alias: `${res.room.alias}`,
                channel_id: res.room.room_id,
                room_id: res.room.streams[selectedStream],
                name: name.toLowerCase(),
                pathname: `/${res.room.alias}`,
                server_id: serverID,
                joined: true,
                room_type: selectedStream,
                streams: res.room.streams,
            }
            store.addRoom(`/${server}`, newRoom)
            kill()
        }
    }).then(() => {
        creating = false
    })
}

function updateTitle(e) {
  const letters = /^[0-9a-zA-Z-]+$/;
  if(!e.key.match(letters)){
    e.preventDefault()
  }
    if(e.key == 'Enter') {
        create()
    }
}

let selectedStream = 'chat';

$: isChat = selectedStream == 'chat'
$: isTopics = selectedStream == 'topics'

function selectStream(type) {
    selectedStream = type
}

let streams = [
    {
        text: 'Topics',
        value: 'topics',
    },
    {
        text: 'Chat',
        value: 'chat',
        default: true,
    },
]

let focused;

function focusInput() {
    focused = true
}
function blurInput() {
    focused = false
}

let isPrivate = false;

function togglePrivate() {
    isPrivate = !isPrivate
}

</script>

{#if active}
<div class="mask gr-default" 
    on:click|self={kill}
    transition:fade="{{duration: 100}}">
    <div class="modal gr-center flex flex-column pa3" 
        in:fly="{{ y: -200, duration: 100 }}">
        <div class="flex">
            <div class="flex-one f4 bold">
                Create Channel
            </div>
            <div class="gr-center icon pointer" on:click={kill}>
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24" height="24"><path fill-rule="evenodd" d="M5.72 5.72a.75.75 0 011.06 0L12 10.94l5.22-5.22a.75.75 0 111.06 1.06L13.06 12l5.22 5.22a.75.75 0 11-1.06 1.06L12 13.06l-5.22 5.22a.75.75 0 01-1.06-1.06L10.94 12 5.72 6.78a.75.75 0 010-1.06z"></path></svg>
            </div>
        </div>


        <div class="fr-l pt3">
            channel name
        </div>

        <div class="input-holder mt2 flex" 
            class:focused={focused}
            on:click={focus}>
            <div class="ic gr-center">
                {#if isChat}
                    {@html hash}
                {:else}
                    {@html topics}
                {/if}
            </div>
            <div class="gr-center flex-one ml1">
                <input 
                bind:value={name}
                bind:this={nameInput}
                on:focus={focusInput}
                on:blur={blurInput}
                maxlength="100"
                on:keypress={updateTitle}
                placeholder={"new-channel"}/>
            </div>
        </div>

        <div class="fr-l mt3">
            default stream
        </div>

        <div class="pt2 fl-co">

            <div class="s-co flex"
                on:click={() => selectStream("chat")}
                class:selected={isChat}>
                <div class="gr-center">
                    <div class="rad gr-default">
                        {#if isChat}
                            <div class="radi gr-center"></div>
                        {/if}
                    </div>
                </div>
                <div class="ic gr-center ml3">
                    {@html hash}
                </div>
                <div class="gr-center flex-one ml3 fl-co">
                    <div class="">
                        Chat
                    </div>
                    <div class="sm mt1">
                        Post chat messages, attachments, GIFs
                    </div>
                </div>
            </div>


            <div class="s-co flex mt2"
                on:click={() => selectStream("topics")}
                class:selected={isTopics}>
                <div class="gr-center">
                    <div class="rad gr-default">
                        {#if isTopics}
                            <div class="radi gr-center"></div>
                        {/if}
                    </div>
                </div>
                <div class="ic gr-center ml3">
                    {@html topics}
                </div>
                <div class="gr-center flex-one ml3 fl-co">
                    <div class="">
                        Topics
                    </div>
                    <div class="sm mt1">
                        Longform discussion
                    </div>
                </div>
            </div>




        </div>

        <div class="pt3 flex pointer" on:click={togglePrivate}>

            <div class="ic gr-center lock">
                {@html lock}
            </div>

            <div class="gr-center flex-one ml2">
                Private Channel
            </div>
            <div class="gr-center">
                <Toggle bind:value={isPrivate}/>
            </div>

        </div>


        <div class="pt3 flex">
            <div class="flex-one">
            </div>
            <div class="gr-center">
                <span class="pointer can pa2" on:click={kill}>
                    Cancel
                </span>
            </div>
            <div class="ml3">
                <button class="" 
                disabled={creating}
                on:click={create}>
                    {#if creating}
                        Creating Channel...
                    {:else}
                        Create Channel
                    {/if}
                </button>
            </div>
        </div>
    </div>
</div>
{/if}

<style>
.mask {
    transition: 0.3s;
    position: fixed;
    z-index: 1000;
    width: 100%;
    height: 100%;
    background-color: var(--mask);
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
}

.modal {
    background-color: var(--background-2);
    width: 440px;
    border-radius: 7px;
    transition: 0.2s;
    box-shadow: 0 30px 60px rgba(0,0,0,.1);
}

.input-holder {
    cursor: text;
    background-color: var(--background-1);
    padding: 0.5rem;
    transition: 0.1s;
    border-radius: 2px;
    border: 1px solid transparent;
}

.focused {
    border: 1px solid var(--green);
}

input {
    font-size: 1.2rem;
    background-color: var(--background-1);
    width: 100%;
    border: none;
}

button {
    padding: 0.5rem 0.5rem;
    font-size: 1rem;
}
.fr-l {
    font-size: 0.74rem;
    text-transform: uppercase;
    letter-spacing: 1px;
    color: var(--text-light);
}

.s-co {
    border-radius: 1px;
    padding: 0.5rem 1rem;
    background-color: var(--background-3);
    cursor: pointer;
    transition: 0.1s;
}
.s-co:hover {
    background-color: var(--background-3);
}

.rad {
    border-radius: 50%;
    border: 2px solid var(--white);
    width: 20px;
    height: 20px;
}
.radi {
    border-radius: 50%;
    border: 3px solid var(--background-2);
    background-color: var(--green);
    width: 100%;
    height: 100%;
}

.selected {
    background-color: var(--background-1);
}
.selected:hover {
    background-color: var(--background-1);
}
.ic {
    fill: var(--text);
}
.wic {
    fill: var(--white);
}
.sm {
    font-size: 0.8rem;
    color: var(--text);
}
.lock {
    width: 16px;
    height: 16px;
}
.can:hover {
    text-decoration: underline;
}
</style>

