<script>
import { tick, onMount, createEventDispatcher } from 'svelte'
import { fade, fly } from 'svelte/transition'

const dispatch = createEventDispatcher()


export let offset = [0,0];
export let trigger = "click";

export let shadow = `box-shadow: 0 0px 10px rgba(0,0,0,.08)`;
export let borderRadius = `8px`

export let initActive;

$: if(initActive) {
    activate()
}

let active = false;

function activate() {
    active = !active
    if(!active) {
        ready = false
    }
}

let refp;
let reference;
let conp;
let container;

export let placement = "top"


let positions = {}
let activePosition = {}

onMount(() => {
})

let ready = false

$: isTop = placement == "top" || 
    placement == "top-start" || 
    placement == "top-end"

$: isBottom = placement == "bottom" || 
    placement == "bottom-start" || 
    placement == "bottom-end"

$: isLeft = placement == "left" || 
    placement == "left-start" || 
    placement == "left-end"

$: isLeftStart = placement == "left-start" 
$: isLeftEnd = placement == "left-end" 

$: if(active) {
    let p = reference?.getBoundingClientRect()
    let c = container?.getBoundingClientRect()

    positions['bottom'] = {
        top: p.top + p?.height + offset[1],
        left: p.left - (c?.width / 2) - (p?.width / 2) + offset[0],
    }

    positions['bottom-start'] = {
        top: p.top + p?.height + offset[1],
        left: p.left + offset[0],
    }


    positions['bottom-end'] = {
        top: p.top + p?.height + offset[1],
        left: p.left - (c?.width) + (p?.width) + offset[0],
    }


    positions['top'] = {
        top: p.top - (c?.height) - p?.height - offset[1],
        left: p.left - (c?.width / 2) - (p?.width / 2) + offset[0],
    }

    positions['top-start'] = {
        top: p.top - (c?.height) - p?.height - offset[1],
        left: p.left + offset[0],
    }
    positions['top-end'] = {
        top: p.top - (c?.height) - p?.height - offset[1],
        left: p.left - (c?.width) + (p?.width) + offset[0],
    }

    positions['left'] = {
        top: p.top - (c?.height / 2) + (p?.height /2) + offset[1],
        left: p.left - c?.width + offset[0],
    }

    positions['left-start'] = {
        top: p.top + offset[1],
        left: p.left - c?.width + offset[0],
    }

    positions['left-end'] = {
        top: p.top - c?.height + p?.height + offset[1],
        left: p.left - c?.width + offset[0],
    }

    activePosition = positions[placement]

    if(isTop) {
        let y = p.top - (c?.height) - offset[1]
        if(y < 0) {
            activePosition.top = p.top + p?.height + offset[1]
        }
    }

    if(isBottom) {
        let top = p.top + (c?.height) + p?.height + offset[1]
        if(top > document.body.clientHeight) {
            activePosition.top = p.top - (c?.height) - p?.height - offset[1]
        }
    }

    if(isLeftStart) {
        let top = activePosition.top + c?.height
        if(top > document.body.clientHeight) {
            let offsetBy = top - document.body.clientHeight + 20
            activePosition.top = activePosition.top - offsetBy
        }
    }

    if(isLeftEnd) {
        if(activePosition.top < 0) {
            activePosition.top = activePosition.top + (-activePosition.top) + 20
        }
    }


    let x = activePosition.left + c?.width + p?.width
    if(x > document.body.clientWidth) {
        let o = x - document.body.clientWidth
        activePosition.left = activePosition.left  - o
    }


    ready = true

    dispatch('ready', true)

    setTimeout(() => {
        document.addEventListener('click', killMe)
        document.addEventListener('keydown', escape)
    }, 10)
} else {
    document.removeEventListener('click', killMe)
    document.removeEventListener('keydown', escape)
    dispatch('killed', true)
}

let repeat = false;

let killMe = (e) => {
    if(e.target !== container && 
        !container?.contains(e.target)) {
        kill()
    }
}
let escape = (e) => {
    if(e.code == 'Escape' || e.key == 'Escape') {
        kill()
    }
}

export function kill(e) {
    active = false
    ready = false
}

async function focusSearchInput() {
    await tick()
}







</script>

<div class="gr-center" bind:this={reference} on:click={activate}>
    <slot name="reference">
    </slot>
</div>


{#if ready}
    <div class="layer" 
        class:inactive={!active}>

        <div class="popup-container" 
        style={`--top:${activePosition?.top}px;--left:${activePosition?.left}px;--shadow:${shadow};--border-radius:${borderRadius}`}
        in:fly="{{ y: 30, duration: 100 }}"
        bind:this={container}>

            <slot name="content"></slot>

        </div>
    </div>
{/if}


<style>

.inactive {
    visibility: hidden;
    pointer-events: none!important;
}
.layer {
    position: fixed;
    pointer-events: none!important;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: none!important;
    z-index: 1002;
}

.popup-container {
    position: absolute;
    left: var(--left);
    pointer-events: auto;
    overflow: hidden;
    outline: 0;
    transition-property: transform, visibility, opacity;
    z-index: 10000;
    box-shadow: var(--shadow);
    top: var(--top);
    border-radius: var(--border-radius);
}
.shadow {
}
</style>

