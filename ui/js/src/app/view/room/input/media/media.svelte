<script>
import { store } from '../../../../store/store.js'
import { useLocation } from 'svelte-navigator'
import { fade, fly } from 'svelte/transition'
import { onMount, onDestroy } from 'svelte'
import { add, upload, thread } from '../../../../utils/icons.js'
import { createEventDispatcher } from 'svelte'
import { makeid } from '../../../../utils/utils.js'

const dispatch = createEventDispatcher()

export let room;
export let editorContent;

export let page;

const location =  useLocation()
$: indicate = $location.pathname == page?.pathname


let fileInput;
let uploading = false;
let uploaded = false;
let file;
let url;

let mxc;

function addMedia() {
    kill()
    fileInput.click()
}


let files = [];

let items = [];

let ready = false;

let tooLarge = false

function killTooLarge() {
    tooLarge = false
}

let build = (e) => {

    if(e.target.files.length > 13) {
        alert("That's too many attachments at once.")
        return
    }

    for(let i =0 ; i < e.target.files.length ; i++) {

        const file = e.target.files[i]

        if (file.size > 8388608) {
            files = []
            files = files
            tooLarge = true
            break
        }
        files = [...files, e.target.files[i]]

    }


    for(let i =0 ; i < files.length; i++) {

        var reader = new FileReader();
        const file = files[i]
        reader.readAsDataURL(file);

        reader.onload = e => {

            /*

            uploadFile(file).then(res => {
                console.log(res)
                let item = {
                    url: res.content_uri,
                    info: {
                        size: file.size,
                        mimetype: file.type,
                    },
                    msgtype: 'm.file',
                    room_id: room.room_id,
                }
                if(file.type.includes('image')) {
                    var image = new Image();
                    image.src = URL.createObjectURL(file)
                    image.onload = () => {
                        item.info.h = image.height
                        item.info.w = image.width
                        item.msgtype = 'm.image'
                        dispatch('uploaded', item)
                    }
                } else{
                    dispatch('uploaded', item)
                }
            })
            */

            let item = {
                id: makeid(16),
                info: {
                    size: file.size,
                    mimetype: file.type,
                    filename: file.name,
                },
                msgtype: 'm.file',
                room_id: room.room_id,
                file: file,
            }

            if(file.type.includes('image')) {
                var image = new Image();
                image.src = URL.createObjectURL(file)
                image.onload = () => {
                    item.info.h = image.height
                    item.info.w = image.width
                    item.msgtype = 'm.image'
                    items.push(item)
                    items = items
                }
            } else if(file.type.includes('video')) {
                let video = document.createElement('video');
                video.src = URL.createObjectURL(file)
                video.addEventListener('loadedmetadata', function(e){
                    item.info.h = video.videoHeight
                    item.info.w = video.videoWidth
                    item.msgtype = 'm.video'
                    items.push(item)
                    items = items
                });
            } else if(file.type.includes('audio')) {
                    item.msgtype = 'm.audio'
                    items.push(item)
                    items = items
            } else {
                items.push(item)
                items = items
            }
        }
    }

    ready = true

    files = []
    fileInput.value = ''
}


async function uploadFile(file) {
    let endpoint = `${homeServer}/_matrix/media/r0/upload`

    let resp = await fetch(endpoint, {
    method: 'POST', // or 'PUT'
    body: file,
    headers:{
        'Authorization': `Bearer ${account.matrix_access_token}`,
        'Content-Type': file.type
    }
    })
    const ret = await resp.json()
    return Promise.resolve(ret)
}

$: account = $store.accounts.filter(account => account.user_id == $store.active_account)[0]


function createThread() {
    store.newThread(room, null, editorContent)
    dispatch('createThread', true)
}

onMount(() => {
})

let menu;
let active;

let top;

function toggle() {
    console.log(container.getBoundingClientRect())
    top = container.getBoundingClientRect().top
    active = !active
}

$: if(active) {
    setTimeout(() => {
        document.addEventListener('click', killMe)
    }, 10)
} else {
    document.removeEventListener('click', killMe)
}

onDestroy(() => {
    document.removeEventListener('click', killMe)
})

let container;

let killMe = (e) => {
    if(e.target !== container && !container?.contains(e.target)){
        kill()
    }
}

function kill() {
    active = false
}

$: if(items?.length > 0 && ready) {
    dispatch('files-ready', items)
}

</script>

<div class="unset gr-default relative no-select" 
    bind:this={container} 
    on:dblclick={addMedia}
    on:click={toggle}>
    <div class="stick add gr-start-center gr-default">
        {@html add}
        <input 
            type="file" 
            name="images"
            bind:this={fileInput} 
            on:change={build} 
            hidden 
            multiple
        >

    </div>

    {#if active}
    <div class="media-container flex flex-column" 
        style={`--top: ${top}px;`}
        bind:this={menu}
        class:dis={active}
        in:fly="{{ y: 20, duration: 80 }}">
        <div class="item flex" on:click={addMedia}>
            <div class="gr-center pr2">
                {@html upload}
            </div>
            <div class="gr-center flex-one">
                Upload a file
            </div>
        </div>
        <div class="item flex" on:click={createThread}>
            <div class="gr-center pr2">
                {@html thread}
            </div>
            <div class="gr-center flex-one">
                Create a thread
            </div>
        </div>
    </div>
    {/if}

</div>




{#if tooLarge}
    <div class="mask gr-default" 
        on:click={killTooLarge}
        transition:fade="{{duration: 100}}">
        <div class="modal gr-center" 
            in:fly="{{ y: -200, duration: 100 }}">

            <div class="too-large pa2">
                <div class="flex flex-column tl-i pa3">
                    <div class="gr-center bold">
                        Your files are too large.
                    </div>
                    <div class="gr-center pt2">
                        Max file size is 8 MB.
                    </div>
                </div>
            </div>

        </div>
    </div>

{/if}

<style>
.add {
    padding-top: 12px;
    padding-bottom: 12px;
    width: 54px;
    align-self: start;
    fill: var(--icon);
    cursor: pointer;
}
.add:hover {
    fill: var(--white);
}

.media-container {
    position: fixed;
    top: calc(var(--top) - 108px);
    background-color: var(--m-bg);
    color: var(--primary-text);
    padding: 0.5rem;
    border-radius: 4px;
    box-shadow: 0 20px 40px rgba(0,0,0,.18);
    display: none;
    min-width: 200px;
}

.dis {
    display: block;
}

.item {
    cursor: pointer;
    padding: 0.5rem;
    border-radius: 3px;
}

.item:hover {
    background-color: var(--green);
}
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
    transition: 0.2s;
    box-shadow: 0 30px 60px rgba(0,0,0,.1);
}

.media-item {
    background-color: var(--background-3);
    border-radius: 7px;
    min-wdith: 530;
    max-height: 80vh;
}

.too-large {
    background-color: hsl(359,calc(var(--saturation-factor, 1)*66.7%),54.1%);
    border-radius: 12px;
    color: var(--white);
    width: 310px;
}

.tl-i {
    border-radius: 10px;
    border: 2px dashed var(--text);
}
.no-dis {
  display: none;
}
</style>
