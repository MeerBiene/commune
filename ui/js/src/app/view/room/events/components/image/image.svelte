<script>
export let image;
export let width;
export let height;
export let isGIF;
export let poster;
export let splitView;

function viewImage() {
    if(isGIF) {
        window.imageViewer(image, true)
    } else {
        window.imageViewer(image)
    }
}

$: src = image.includes('mxc://') ?
    `${homeServer}/_matrix/media/r0/download/${image?.substring(6)}` : image

$: pix = splitView ? 120 : 400

$: base = width < pix ? width : pix

$: w = width > base ? base : width
$: h = (base/width) * height

</script>

<div class="image-item pointer"
    style={`--width:${w}px`}
on:click={viewImage}>
    {#if isGIF}
        <video poster={poster}
        src={src}
        width={w}
        height={h}
        loop autoplay muted></video>
    {:else}
        <img src={src} width={w} height={h} loading=lazy/>
    {/if}
</div>

<style>
.image-item {
    width: var(--width);
}
</style>
