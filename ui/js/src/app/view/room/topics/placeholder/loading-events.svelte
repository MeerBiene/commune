<script>

import LoadingEventItem from './loading-event-item.svelte'

$: titleWidth = Math.floor(Math.random() * (100 - 50) + 50);

$: titleWords = buildWords()

function buildWords() {

    let items = []

    let count = Math.floor(Math.random() * (20 - 7) + 7);

    for(var i = 0; i < count; i++) {
        let width = Math.floor(Math.random() * (90 - 30) + 30);
        items.push(width)
    }

    return items
}


let count = 4;

</script>

<div class="t-c mt4">
    <div class="title mb3" style={`--width:${titleWidth}%;`}>
    </div>

    <div class="fl-co">
        {#each Array(count) as _, i}
            <LoadingEventItem />
        {/each}
    </div>
</div>

<style>
.t-c {
    width: 100%;
    justify-self: center;
    display: grid;
    grid-template-rows: auto 1fr;
}
.title {
    border-radius: 500px;
    background-color: var(--background-13);
    height: 20px;
    width: var(--width);
}

</style>
