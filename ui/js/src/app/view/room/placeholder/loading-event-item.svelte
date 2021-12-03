<script>
import EventItemLine from './event-item-line.svelte'

export let topics;

$: lines = Math.floor(Math.random() * (6 - 1) + 1);

$: isImage = random() == 1
$: imageHeight = Math.floor(Math.random() * (400 - 100) + 100);

function random() {
    let items = [0, 0, 0, 0, 0, 0, 1, 1];
    let x = Math.floor(Math.random() * items.length);
    return items[x];
}

</script>

{#if topics}
    <div class="tp-i mb3 flex pa3 ma3">
        <div class="ava mr3">
            <div class="profile-avatar">
            </div>
        </div>
        <div class="fl-co flex-one">
            <div class="">
                <EventItemLine title={true}/>
            </div>
            <div class="mt3">
                <EventItemLine label={true}/>
            </div>
        </div>
    </div>
{:else}
    <div class="it-c mb3 flex">
        <div class="ava mr3">
            <div class="profile-avatar">
            </div>
        </div>
        <div class="fl-co">

            <div class="mt1">
                <EventItemLine name={true}/>
            </div>

            {#if isImage}
                <div class="image" style={`--height:${imageHeight}px;`}>
                </div>
            {:else}
                {#each Array(lines) as _, i}
                    <EventItemLine/>
                {/each}

            {/if}

        </div>
    </div>
{/if}

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

.image {
    background-color: var(--background-11);
    opacity: 0.7;
    height: var(--height);
    width: 400px;
}

.tp-i {
    background-color: var(--background-12);
    border-radius: 7px;
}

</style>
