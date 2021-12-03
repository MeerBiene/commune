<script>
import { onMount } from 'svelte'

let source;
$: src = source?.includes('mxc://') ?
    `${homeServer}/_matrix/media/r0/download/${source?.substring(6)}` : source

$: multiple = source?.length > 0

let active = false;

let isGIF = false;

onMount(() => {
    window.imageViewer = (src, gif) => {
        isGIF = gif
        source = src
        active = true
    }
})

function kill() {
    active = false
    isGIF = false
}

</script>


{#if active && source}
    {#if isGIF}
        <div class="image-viewer gr-default" on:click|self={kill}>
            <div class="image-container gr-center">
                <video 
                src={src}
                loop autoplay muted></video>

            </div>
        </div>
    {:else}
        <div class="image-viewer gr-default" on:click|self={kill}>
            <div class="image-container gr-center">
                <img src={src} loading=lazy/>
            </div>
        </div>
    {/if}
{/if}

<style>
.image-viewer {
    position: fixed;
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
    background: var(--mask);
}
.image-container{
    max-width: 80vw;
}
</style>
