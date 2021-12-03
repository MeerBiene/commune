<script>
import { store } from '../../../../store/store.js'
import Popup from '../../../../components/popup/popup.svelte'
import ThreadItem from './thread-item.svelte'
import { threadBig, closeBig } from '../../../../utils/icons.js'

export let room;
export let threads;

let container;

function kill() {
    container.kill()
}

function createThread() {
    kill()
    store.newThread(room, null, null)
}

</script>

<Popup 
bind:this={container}
trigger={"click"}
shadow={`0 10px 20px rgba(0,0,0,.15)`}
placement={"bottom-end"}
offset={[0, 6]}>
    <div class="t-c pointer gr-center mr1 flex" slot="reference">
        <div class="gr-center">
            {@html threadBig}
        </div>
        {#if threads.length > 0}
        <div class="gr-center ml1">
            {threads.length}
        </div>
        {/if}
    </div>
    <div class="thread-container flex flex-column" slot="content">
        <div class="header pv2 ph3">
            <div class="t-c gr-center">
                {@html threadBig}
            </div>
            <div class="lh-1 gr-center w-100">
                Threads
            </div>
            <div class="gr-center">
                <button on:click={createThread}>
                    Create
                </button>
            </div>
            <div class="t-c gr-center pointer" on:click={kill}>
                {@html closeBig}
            </div>
        </div>
        <div class="content fl-co pa3">
            {#each threads as thread, i (thread.room_id)}
                <ThreadItem 
                last={i == threads.length - 1}
                thread={thread} />
            {/each}

            {#if threads.length == 0}
                <div class="gr-center">
                    No threads.
                </div>
            {/if}

        </div>
    </div>
</Popup>


<style>
.t-c {
    fill: var(--icon);
}
.t-c:hover {
    fill: var(--white);
}
.thread-container {
    border-radius: 10px;
    border: 1px solid var(--background-1);
    background-color: var(--background-3);
    width: 440px;
    max-height: 448px;
    color: var(--text);
    line-height: 1.4;
    display: grid;
    overflow: hidden;
    grid-template-rows: auto 1fr;
}
.header {
    background-color: var(--background-1);
    display: grid;
    grid-template-columns: auto 1fr auto auto;
    grid-column-gap: 1rem;
}
.content {
    background-color: var(--background-2);
}
</style>
