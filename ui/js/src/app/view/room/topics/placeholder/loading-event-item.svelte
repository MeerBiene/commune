<script>

$: nameWidth = Math.floor(Math.random() * (100 - 50) + 50);

$: words = buildWords()

function buildWords() {

    let items = []

    let count = Math.floor(Math.random() * (90 - 10) + 10);

    for(var i = 0; i < count; i++) {
        let width = Math.floor(Math.random() * (90 - 30) + 30);
        items.push(width)
    }

    return items
}

</script>

<div class="tp-i mb3 flex pa3">
    <div class="ava mr3">
        <div class="profile-avatar">
        </div>
    </div>
    <div class="fl-co flex-one">
        <div class="">
            <div class="name mb2" style={`--width:${nameWidth}px;`}>
            </div>
        </div>
        <div class="mt1 flex flex-wrap">
            {#each words as width, i}
                <div class="word mr2 mb2" style={`--width:${width}px;`}>
                </div>
            {/each}
        </div>
    </div>
</div>

<style>
.profile-avatar {
    width: 40px;
    height: 40px;
    background-color: var(--avatar);
    border-radius: 50%;
    opacity: 0.5;
}
.ava {
    min-width: calc(40px);
    margin-top: 4px;
}

.tp-i {
    background-color: var(--background-12);
    border-radius: 7px;
}
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
</style>
