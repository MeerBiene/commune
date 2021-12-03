<script>
import {createEventDispatcher} from 'svelte'
const dispatch = createEventDispatcher()
import { onMount } from 'svelte'
import Avatar from '../../../components/avatar/avatar.svelte'
import { store } from '../../../store/store.js'

let name = '';
let nameInput;

let avatar;

let showSlug = false;

onMount(() => {
    nameInput.focus()
})

$: slug = slugify(name)

$: noSlug = name.length == 0

function slugify(name) {
    if(name) {
        let x = name.split(" ")
        x = x.join("-")
        x = x.toLowerCase()
        let slug = x.replace(/[^a-z0-9-]/gi,'');
        return slug
    } else {
        return `my-cool-server`
    }
}

$: account = $store.accounts.filter(account => account.user_id == $store.active_account)[0]

async function createServer() {
    let endpoint = `/server/create`
    let data = {
        title: name,
    };
    if(avatar?.length > 0) {
        data.avatar = avatar
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
            let server = {
                room_id: res.room.room_id,
                name: name,
                alias: `${res.room.alias}`,
                short_alias: `${res.room.short_alias}`,
                owner: true,
                pathname: `/${res.room.alias}`,
                joined: true,
                default_room: `/${res.room.alias}/${res.room.default_chatroom.alias}`,
                active_room: `/${res.room.alias}/${res.room.default_chatroom.alias}`,
                rooms: [
                    {
                        room_id: res.room.default_chatroom.streams['chat'],
                        channel_id: res.room.room_id,
                        server_id: res.room.room_id,
                        server_alias: `/${res.room.alias}`,
                        name: 'general',
                        alias: `${res.room.default_chatroom.alias}`,
                        pathname: `/${res.room.default_chatroom.alias}`,
                        joined: true,
                        room_type: 'chat',
                        streams: res.room.default_chatroom.streams,
                    }
                ]
            }

            if(avatar?.length > 0) {
                server.avatar = avatar;
            }
            store.addServer(server)
            kill()
        }
    }).then(() => {
        creating = false
    })
}

function kill() {
    dispatch('kill', true)
}

function goBack() {
    dispatch('go-back', true)
}


function uploaded(e) {
    console.log(e.detail)
    avatar = e.detail
}

function onKeyPress(e) {
    if(e.key == 'Enter') {
        create()
    }
}

</script>

<div class="pa3 flex flex-column">
    <div class="flex">
        <div class="flex-one f3 bold">
            Customize your server
        </div>
        <div class="gr-center icon pointer" on:click={kill}>
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24" height="24"><path fill-rule="evenodd" d="M5.72 5.72a.75.75 0 011.06 0L12 10.94l5.22-5.22a.75.75 0 111.06 1.06L13.06 12l5.22 5.22a.75.75 0 11-1.06 1.06L12 13.06l-5.22 5.22a.75.75 0 01-1.06-1.06L10.94 12 5.72 6.78a.75.75 0 010-1.06z"></path></svg>
        </div>
    </div>
    <div class="pt4 gr-center">
        <Avatar on:uploaded={uploaded}/>
    </div>
    <div class="fr-l pt4">
        name
    </div>
    <div class="pt2">
        <input 
        bind:value={name}
        bind:this={nameInput}
        maxlength="100"
        on:keypress={onKeyPress}
        placeholder="My Cool Server"/>
    </div>
    {#if showSlug}
    <div class="pt2 slug">
        https://commune.chat/<span class:mute={noSlug}>{slug}</span>
    </div>
    {/if}
    <div class="pt3 flex">
        <div class="gr-center icon pointer" on:click={goBack}>
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24" height="24"><path fill-rule="evenodd" d="M10.78 19.03a.75.75 0 01-1.06 0l-6.25-6.25a.75.75 0 010-1.06l6.25-6.25a.75.75 0 111.06 1.06L5.81 11.5h14.44a.75.75 0 010 1.5H5.81l4.97 4.97a.75.75 0 010 1.06z"></path></svg>
        </div>
        <div class="flex-one">
        </div>
        {#if creating}
        <div class="spinner-s gr-center mr3">
        </div>
        {/if}
        <div class="">
            <button class="" 
            disabled={creating}
            on:click={create}>
                {#if creating}
                    Creating server...
                {:else}
                    Create Server
                {/if}
            </button>
        </div>
    </div>


</div>

<style>
input {
    font-size: 1.2rem;
    background-color: var(--background-3);
    width: 100%;
    padding: 1rem;
    border: 1px solid transparent;
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
.slug {
    color: var(--text-light);
    word-break: break-all;
}
.mute {
    opacity: 0.5;
}
</style>
