<script>

export let name;
export let title;
export let label;

$: words = buildWords()

$: nameWidth = Math.floor(Math.random() * (100 - 50) + 50);
$: titleWidth = Math.floor(Math.random() * (90 - 40) + 40);

function buildWords() {

    let items = []

    let count = Math.floor(Math.random() * (8 - 3) + 3);

    for(var i = 0; i < count; i++) {
        let width = Math.floor(Math.random() * (90 - 30) + 30);
        items.push(width)
    }

    return items
}

</script>

{#if name}
    <div class="name mb2" style={`--width:${nameWidth}px;`}>
    </div>
{:else if title}
    <div class="title" style={`--width:${titleWidth}%;`}>
    </div>
{:else if label}
    <div class="label">
    </div>
{:else}
    <div class="flex">
        {#each words as width, i}
            <div class="word mr2 mb2" style={`--width:${width}px;`}>
            </div>
        {/each}
    </div>
{/if}

<style>
.word {
    border-radius: 500px;
    background-color: var(--background-9);
    height: 16px;
    width: var(--width);
}
.name {
    border-radius: 500px;
    background-color: var(--background-13);
    height: 16px;
    width: var(--width);
}
.title {
    border-radius: 500px;
    background-color: var(--background-9);
    height: 16px;
    width: var(--width);
}
.label {
    border-radius: 500px;
    background-color: var(--background-9);
    height: 16px;
    width: 100px;
    opacity: 0.5;
}
</style>
