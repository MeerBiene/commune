<script>
import { onMount, createEventDispatcher } from 'svelte'
const dispatch = createEventDispatcher()


let active = false;

function insert() {
    active = !active;
    let opts = {
        container: container,
        dispatch: true,
    }
    window.toggleEmojiPicker(opts)
}

let cords = [
    {x: 0, y: 0},
    {x: -22, y: 0},
    {x: -44, y: 0},
    {x: -66, y: 0},
    {x: -88, y: 0},
    {x: -110, y: 0},
    {x: -132, y: 0},
    {x: -154, y: 0},
    {x: -176, y: 0},
    {x: -198, y: 0},
    {x: -220, y: 0},
    {x: 0, y: -22},
    {x: -22, y: -22},
    {x: -44, y: -22},
    {x: -66, y: -22},
    {x: -88, y: -22},
    {x: -110, y: -22},
    {x: -132, y: -22},
    {x: -154, y: -22},
    {x: -176, y: -22},
    {x: -198, y: -22},
    {x: -220, y: -22},
    {x: 0, y: -44},
    {x: -22, y: -44},
    {x: -44, y: -44},
    {x: -66, y: -44},
    {x: -88, y: -44},
    {x: -110, y: -44},
    {x: -132, y: -44},
    {x: -154, y: -44},
    {x: -176, y: -44},
    {x: -198, y: -44},
    {x: -220, y: -44},
    {x: 0, y: -66},
    {x: -22, y: -66},
    {x: -44, y: -66},
    {x: -66, y: -66},
    {x: -88, y: -66},
    {x: -110, y: -66},
    {x: -132, y: -66},
    {x: -154, y: -66},
    {x: -176, y: -66},
    {x: -198, y: -66},
    {x: -220, y: -66},
    {x: 0, y: -88},
    {x: -22, y: -88},
    {x: -44, y: -88},
    {x: -66, y: -88},
    {x: -88, y: -88},
    {x: -110, y: -88},
]

let selected = cords[0]

function switchEmoji() {
    selected = cords[Math.floor(Math.random() * cords.length)]
}

let container;

onMount(() => {
    window.dispatchEmoji = (emoji) => {
        dispatch('insert', emoji)
    }
    window.killEmoji = (kill) => {
        if(kill == true) {
            active = false
        }
    }
})

</script>

<div class="unset flex mh3 ">
    <div class="emoji-b-c stick">
        <div class="emoji-button" 
            class:active={active}
            bind:this={container}
            on:mouseover={switchEmoji}
            on:click={insert}
            style="--x:{`${selected.x}px`};--y: {`${selected.y}px`}">
        </div>
    </div>
</div>


<style>
.emoji-b-c {
    padding-top: 13px;
    padding-bottom: 12px;
    display: grid;
}

.emoji-button {
    height: 22px;
    width: 22px;
    background-image: url('/static/img/emoji-sprite.png');
    background-size: 242px 110px;
    background-position-x: var(--x);
    background-position-y: var(--y);
    transform: scale(1);
    filter: grayscale(100%);
    cursor: pointer;
}

.active {
    filter: grayscale(0%);
    transform: scale(1.14);
}

.emoji-button:hover {
    filter: none;
    transform: scale(1.14);
}
</style>
