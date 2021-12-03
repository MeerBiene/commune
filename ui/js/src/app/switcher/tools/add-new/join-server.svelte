<script>
import {createEventDispatcher} from 'svelte'
const dispatch = createEventDispatcher()
import { onMount } from 'svelte'
import { store } from '../../../store/store.js'
import { makeid } from '../../../utils/utils.js'

let url = '';
let urlInput;

let avatar;


onMount(() => {
    urlInput.focus()
})


$: account = $store.accounts.filter(account => account.user_id == $store.active_account)[0]

async function resolveFullDomain(url) {
    let endpoint = `${url}/resolve_domain`
    let resp = await fetch(endpoint, {
        headers:{
            'Content-Type': 'application/json'
        }
    })
    const ret = await resp.json()
    return Promise.resolve(ret)
}

async function resolveAlias(alias) {
    console.log(alias)
    let endpoint = `${homeServer}/_matrix/client/r0/directory/room/${alias}`
    console.log(endpoint)
    let resp = await fetch(endpoint, {
        headers:{
            'Content-Type': 'application/json'
        }
    })
    const ret = await resp.json()
    return Promise.resolve(ret)
}

let error;

let creating = false;

let showWarning = false;

function isValidFormat(d) {
    if(d.slice(-1) == '/') {
        d = d.substring(0, d.length - 1);
    }
    let s = d.split('/')
    return s.length === 4 &&
        (s[0] == 'http:' || s[0] == 'https:') &&
        s[2]?.length > 0 &&
        s[3]?.length >= 7 
}

function isFederated(d) {
    let s = d.split('/')
    return s[2] != shortlinkDomain
}

function urlOrigin(d) {
    let s = d.split('/')
    s = s.slice(0, -1)
    return s.join('/')
}

function stripShortlink(d) {
    let s = d.split('/')
    return s.slice(-1)
}


$: federated = isFederated(url)
$: shortlink = stripShortlink(url)

function create() {
    if(urlInput.value.length < 4) {
        showWarning = true
        urlInput.focus()
        return
    }
    const isValid = isValidFormat(url)
    if(!isValid) {
        showWarning = true
        urlInput.focus()
        return
    }



    if(federated) {
        joinFederatedServer()
    } else {
        joinServer()
    }


    //creating = true
    
}

function joinFederatedServer() {
    const origin = urlOrigin(url)
    resolveFullDomain(origin).then(res => {
        console.log(res)
        if(res?.server) {
            return res.server
        }
    }).then(server => {
        const alias = `%23${shortlink}:${server}`
        resolveAlias(alias).then(res => {
            console.log(res)
            if(res.error) {
                showWarning = true
            }
            if(res?.room_id) {
                store.joinRoom({
                    room_id: alias,
                    token: account?.matrix_access_token,
                    federated: true,
                })
            }
        })
    })
}

function joinServer() {
    const alias = `%23${shortlink}:${federationDomain}`
    resolveAlias(alias).then(res => {
        console.log(res)
        if(res.error) {
            showWarning = true
        }
        if(res?.room_id) {
            store.joinRoom({
                room_id: res?.room_id,
                token: account?.matrix_access_token,
            })
        }
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

function reset() {
    showWarning = false
}

function onKeyPress(e) {
    if(e.key == 'Enter') {
        create()
    }
}

$: dummyInviteLink = `${location.protocol}//${shortlinkDomain}/${makeid(7)}`
$: dummyInviteLinkAlt = `${location.protocol}//${shortlinkDomain}/astronomy`

</script>

<div class="pa3 flex flex-column">
    <div class="flex">
        <div class="flex-one f3 bold">
            Join a Server
        </div>
        <div class="gr-center icon pointer" on:click={kill}>
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24" height="24"><path fill-rule="evenodd" d="M5.72 5.72a.75.75 0 011.06 0L12 10.94l5.22-5.22a.75.75 0 111.06 1.06L13.06 12l5.22 5.22a.75.75 0 11-1.06 1.06L12 13.06l-5.22 5.22a.75.75 0 01-1.06-1.06L10.94 12 5.72 6.78a.75.75 0 010-1.06z"></path></svg>
        </div>
    </div>
    <div class="pt2">
        Enter an invite below to join an existing server.
    </div>
    <div class="fr-l pt4">
        invite link
        {#if showWarning}
            <span class="warn">
                - Invalid link
            </span>
        {/if}
    </div>
    <div class="pt2">
        <input 
        bind:value={url}
        bind:this={urlInput}
        maxlength="100"
        on:keypress={onKeyPress}
        on:input={reset}
        placeholder={dummyInviteLink}/>
    </div>

    <div class="fr-l pt4">
        invites should look like
    </div>
    <div class="pt2 fl-co">
        <span>{dummyInviteLink}</span>
        <span>{dummyInviteLinkAlt}</span>
    </div>

    <div class="pt3 flex">
        <div class="gr-center icon pointer" on:click={goBack}>
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24" height="24"><path fill-rule="evenodd" d="M10.78 19.03a.75.75 0 01-1.06 0l-6.25-6.25a.75.75 0 010-1.06l6.25-6.25a.75.75 0 111.06 1.06L5.81 11.5h14.44a.75.75 0 010 1.5H5.81l4.97 4.97a.75.75 0 010 1.06z"></path></svg>
        </div>
        <div class="flex-one">
        </div>
        <div class="">
            <button class="" 
            disabled={creating}
            on:click={create}>
                {#if creating}
                    Joining server...
                {:else}
                    Join Server
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
.mute {
    opacity: 0.5;
}
</style>
