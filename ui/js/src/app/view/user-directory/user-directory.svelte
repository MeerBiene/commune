<script>
import { useLocation } from 'svelte-navigator'


const location = useLocation()

$: showMe = $location.pathname == `/users`

$: if(showMe) {
    document.title = `User Directory`
}

</script>

<div class="home" class:no-dis={!showMe}>
    <div class="user-profile w-100 scrl">
        User Directory
    </div>
</div>

<style>
.home {
    height: 100%;
    overflow: hidden;
    display: grid;
}
.user-profile {
    background-color: var(--background-3);
    overflow: auto;
}

.no-dis {
     display: none;
}
</style>
