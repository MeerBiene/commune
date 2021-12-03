<script>
import { fade, fly } from 'svelte/transition'
import { download, expand } from '../../utils/icons.js'
import { formatBytes } from '../../utils/utils.js'
//import hljs from 'highlight.js'


let hljs;
let hljsLoaded;
$: if(isText && !isPlainText) {
    import('highlight.js/lib/common')
    .then(res => {
        hljs = res.default
        hljsLoaded = true
    });
}

export let title;
export let url;
export let size;
export let type;

$: fileSize = formatBytes(size)

$: isText = type.includes('text/')

$: fetchText = size <= 100000
//$: fetchText = isHTML ? true : size <= 100000

let ready = false

let textContent;

$: splitByLines = textContent?.split(/\r\n|\r|\n/)

$: lines = splitByLines?.length || 0

$: shortened = lines > 1 ? splitByLines?.slice(0, 6)?.join('\n') :
    textContent?.substring(0, 50)

$: highlighted = shortened ? hljs ? hljs.highlightAuto(shortened)?.value : `` : ``

$: fullHighlighted = textContent ? hljs ? hljs.highlightAuto(textContent)?.value : `` : ``

$: formatted = shortened ? highlighted : ``


$: if(isText && url && fetchText) {
    fetchFile().then(res => {
        textContent = res
        ready = true
    })
} else {
    ready = true
}

async function fetchFile() {
    let resp = await fetch(url, {
        method: 'GET',
    })
    const ret = await resp.text()
    return Promise.resolve(ret)
}


$: isPlainText = type == 'text/plain'
$: isHTML = type == 'text/html'

let fullView = false;

function viewFull() {
    fullView = true
}

function killFullView() {
    fullView = false
}

function downloadFile() {
    fetch(url).then(res => {
        return res.blob().then(b =>{
            var a = document.createElement("a");
            a.href = URL.createObjectURL(b);
            a.setAttribute("download", title);
            a.click();
        });
    });
}

</script>


{#if isText && fetchText}


<div class="file-item-container text flex">
    <div class="text-content scrl-s pa2">
        {#if ready}
            {#if isPlainText}
                {@html shortened}
            {:else}
                {@html formatted}
            {/if}
        {:else}
            <div class="placeholder flex flex-column">
                <div class="p-t mb2"></div>
                <div class="p-t mb2"></div>
                <div class="p-t mb2"></div>
                <div class="p-t mb2"></div>
                <div class="p-t mb2"></div>
                <div class="p-t mb2"></div>
            </div>
        {/if}
    </div>
    {#if ready}
        <div class="text-footer pa2">
            <div class="expand-icon gr-center pointer" 
                on:click={viewFull}>
                {@html expand}
            </div>
            <div class="gr-center">
            </div>
            <div class="gr-center">
                {title}
            </div>
            <div class="mute pl2 gr-center">
                <span class="">{fileSize}</span>
            </div>
            <div class="download-icon-s gr-center ml2"
            on:click={downloadFile}>
                {@html download}
            </div>
        </div>
    {/if}



{#if fullView}
    <div class="mask gr-default" 
        on:click={killFullView}
        transition:fade="{{duration: 100}}">
        <div class="modal gr-center" 
            in:fly="{{ y: -200, duration: 100 }}">

            <div class="text-content-full scrl-s pa2">
                {#if isPlainText}
                    {@html textContent}
                {:else}
                    {@html fullHighlighted}
                {/if}
            </div>


            <div class="text-footer pa2">
                <div class="gr-center">
                </div>
                <div class="gr-center">
                </div>
                <div class="gr-center">
                    {title}
                </div>
                <div class="mute pl2 gr-center">
                    <span class="">{fileSize}</span>
                </div>
                <div class="download-icon-s gr-center ml2"
                on:click={downloadFile}>
                    {@html download}
                </div>
            </div>

        </div>
    </div>

{/if}



</div>

{:else}

<div class="file-item-container flex pa2">
    <div class="file-item flex flex-column flex-one">
        <div class="">
            <a href={url}>
                <span class="clmp-1">{title}</span>
            </a>
        </div>
        <div class="wh-s pt1">
            <span class="">{fileSize}</span> - <span class="">{type}</span>
        </div>
    </div>
    <div class="download-icon" 
        on:click={downloadFile}>
        {@html download}
    </div>
</div>

{/if}

<style>
.file-item-container {
    min-width: 400px;
    max-width: 520px;
    border-radius: 4px;
    background-color: var(--background-2);
    border: 1px solid var(--background-5);
    overflow: hidden;
}
.text {
    height: 170px;
    min-width: 400px;
    display: grid;
    grid-template-rows: 1fr auto;
    overflow: hidden;
}

.text-content {
    overflow-x: auto;
    overflow-y: hidden;
    white-space: pre;
    font-family: monospace;
}

.text-content-full {
    white-space: break-spaces;
    font-family: monospace;
    overflow: hidden scroll;
}

.text-footer {
    display: grid;
    grid-template-columns: auto 1fr auto auto auto;
    border-top: 1px solid var(--background-5);
    font-size: 0.9rem;
}

.download-icon {
    width: 24px;
    height: 24px;
    cursor: pointer;
    fill: var(--text);
}
.download-icon:hover {
    fill: var(--white);
}
.download-icon-s {
    width: 20px;
    height: 20px;
    cursor: pointer;
    fill: var(--text);
}
.download-icon-s:hover {
    fill: var(--white);
}
.expand-icon {
    width: 16px;
    height: 16px;
    cursor: pointer;
    fill: var(--text);
}
.expand-icon:hover {
    fill: var(--white);
}
.wh-s {
    font-size: 0.7rem;
    color: var(--text);
}
.p-t {
    background-color: var(--background-3);
    height: 10px;
    width: 100%;
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
    height: 80vh;
    width: 950px;
    border-radius: 4px;
    background-color: var(--background-2);
    border: 1px solid var(--background-5);
    overflow: hidden;
    display: grid;
    grid-template-rows: 1fr auto;
}
</style>
