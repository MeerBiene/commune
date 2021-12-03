<script>
import { useLocation } from 'svelte-navigator'


const location = useLocation()

$: showExplore = $location.pathname == '/explore'

</script>

<div class="home" class:no-dis={!showExplore}>
    <div class="chat w-100 scrl">
        TODO
    </div>
</div>

<style>
.home {
    height: 100%;
    overflow: hidden;
    display: grid;
}
.chat {
    overflow: auto;
}

.no-dis {
     display: none;
}
</style>
