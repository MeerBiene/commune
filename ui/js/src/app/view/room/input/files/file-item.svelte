<script>

import { store } from '../../../../store/store.js'
import { onMount, createEventDispatcher } from 'svelte'
import { close } from '../../../../utils/icons.js'
import { formatBytes } from '../../../../utils/utils.js'

const dispatch = createEventDispatcher()

export let file;
export let last;
export let room;


let uploaded = false;

onMount(() => {
    uploadFile()
    if(last) {
        window.syncEventsPosition(room?.room_id)
    }
})

$: fileSize = formatBytes(file?.info?.size)
$: type = file?.info?.mimetype

$: account = $store.accounts.filter(account => account.user_id == $store.active_account)[0]

let content_uri;

function uploadFile() {

    let xhr = new XMLHttpRequest();

    xhr.upload.onprogress = function(event) {
        percent = (event.loaded / event.total) * 100
    };

    xhr.responseType = 'json';

    xhr.onloadend = function() {
        clearInterval(interval)
        percent = 100
        uploaded = true
        console.log(xhr.response)
        content_uri = xhr?.response?.content_uri
        dispatch('file-uploaded', {
            id: file?.id,
            url: content_uri,
        })
    };

    xhr.upload.onerror = function() {
      console.log(`Error during the upload: ${xhr.status}`);
    };

    let endpoint = `${homeServer}/_matrix/media/r0/upload`
    xhr.open("POST", endpoint)

    xhr.setRequestHeader('Authorization', `Bearer ${account.matrix_access_token}`)
    xhr.setRequestHeader('Content-Type', file?.info?.mimetype);


    xhr.send(file?.file);
}

async function uploadFileOld() {
    let endpoint = `${homeServer}/_matrix/media/r0/upload`
    console.log(file)

    let resp = await fetch(endpoint, {
    method: 'POST', 
    body: file?.file,
    headers:{
        'Authorization': `Bearer ${account.matrix_access_token}`,
        'Content-Type': file?.info?.mimetype,
    }
    })
    const ret = await resp.json()
    return Promise.resolve(ret)
}


function discard() {
    dispatch('discard', file.id)
}


let percent = 0;
let interval;

onMount(() => {
    interval = setInterval(() => {
        percent += 2
    }, 500)
})

$: url = `${homeServer}/_matrix/media/r0/download/${content_uri?.substring(6)}`

</script>

<div class="file-item-container mb2 mr3 pa2">
    <div class="flex">
        <div class="file-item flex flex-column flex-one">
            <div class="flex-one bold pr1">
                {#if uploaded && URL}
                    <a href={url}>
                        <span class="clmp-1">{file?.info?.filename}</span>
                    </a>
                {:else}
                    {file?.info?.filename}
                {/if}
            </div>
            <div class="wh-s pt1">
                <span class="">{fileSize}</span> - <span class="">{type}</span>
            </div>
        </div>
        <div class="discard gr-center-start pointer" 
            aria-label="Remove file"
            data-microtip-position="top"
            data-microtip-size="fit"
            role="tooltip" 
            on:click={discard}>
            {@html close}
        </div>
    </div>
    <div class="progress-bar mt2">
        <div class="progress" style={`--width: ${percent}%;`}>
        </div>
    </div>
</div>

<style>
.file-item-container {
    border-radius: 8px;
    background-color: var(--background-2);
    max-width: 520px;
    display: grid;
    grid-template-rows: 1fr auto;
}

.hasImage {
    grid-template-columns: 106px 1fr;
}

.discard {
    width: 20px;
    height: 20px;
    fill: var(--text);
}
.discard:hover {
    fill: var(--white);
}

.href {
    color: var(--text-muted);
}
.file-image {
    height: 106px;
    width: 106px;
    border-radius: 8px;
}
.desc {
    font-size: 0.9rem;
}
.wh-s {
    font-size: 0.7rem;
    color: var(--text);
}
.progress-bar {
    border-radius: 500px;
    height: 3px;
    background-color: var(--background-1);
}
.progress {
    border-radius: 500px;
    height: 3px;
    background-color: var(--green);
    width: var(--width);
    transition: 0.2s;
}
</style>
