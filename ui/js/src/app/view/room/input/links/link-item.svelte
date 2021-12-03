<script>

import { store } from '../../../../store/store.js'
import { onMount, createEventDispatcher } from 'svelte'
import { close } from '../../../../utils/icons.js'

const dispatch = createEventDispatcher()

export let href;
export let room;


let fetched = false;

let metadata = {};

onMount(() => {
    dispatch('fetching-metadata', true)
    fetchMetadata().then(res => {
        console.log(res)
        if(res) {
            fetched = true
            if(res?.title?.length > 0) {
                metadata.title = res?.title
            }
            if(res?.description?.length > 0) {
                metadata.description = res?.description
            }
            if(res?.image?.length > 0) {
                metadata.image = res?.image
            }
            if(res?.author?.length > 0) {
                metadata.author = res?.author
            }
            if(res?.is_youtube && res?.youtube_id?.length > 0) {
                metadata.youtube_id = res?.youtube_id
            }
            if(res?.is_wikipedia) {
                metadata.is_wikipedia = true
            }
            dispatch('metadata-fetched', {
                href: href,
                metadata, metadata,
            })
            window.syncEventsPosition(room?.room_id)
        }
    })
})


$: account = $store.accounts.filter(account => account.user_id == $store.active_account)[0]

async function fetchMetadata() {
    let endpoint = `/link/metadata`

    let data = {
        href: href,
    }

    let resp = await fetch(endpoint, {
    method: 'POST', // or 'PUT'
    body: JSON.stringify(data),
    headers:{
        'Authorization': account.access_token,
    }
    })
    const ret = await resp.json()
    return Promise.resolve(ret)
}


function discard() {
    dispatch('discard', href)
}

$: title = metadata?.title ? metadata?.title : href

$: isYoutube = metadata?.youtube_id?.length > 0

$: imgSrc = `https://img.youtube.com/vi/${metadata?.youtube_id}/mqdefault.jpg`

</script>

<div class="link-item-container mb3 mr3"
    class:hasImage={isYoutube}>
    {#if isYoutube}
        <div class="link-image bg-img" 
            style="background-image: url({imgSrc});">
        </div>
    {/if}
    <div class="link-item pa2 flex flex-column flex-one">
        <div class="flex">
            <div class="">
            </div>
            <div class="flex-one bold pr1">
                <a href={href} target="_blank">
                    <span class="clmp-2">{title} </span>
                </a>
            </div>
            <div class="discard gr-center pointer" 
                aria-label="Remove Link"
                data-microtip-position="top"
                data-microtip-size="fit"
                role="tooltip" 
                on:click={discard}>
                {@html close}
            </div>
        </div>
        {#if fetched && metadata?.description}
            <div class="desc pt2 clmp-2">
                {@html metadata.description}
            </div>
        {/if}
        {#if fetched && metadata?.title}
            <div class="pt2 href clmp-1">
                {href}
            </div>
        {/if}
        {#if !fetched}
            <div class="pt2 href">
                Fetching Metadata...
            </div>
        {/if}
    </div>
</div>

<style>
.link-item-container {
    border-radius: 8px;
    background-color: var(--background-2);
    max-width: 520px;
    display: grid;
}

.hasImage {
    grid-template-columns: 106px 414px;;
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
.link-image {
    height: 106px;
    width: 106px;
    border-radius: 8px;
}
.desc {
    font-size: 0.9rem;
    line-height: 1.2rem;
}
</style>
